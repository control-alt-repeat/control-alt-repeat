package reconcile

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/shopspring/decimal"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

const applyBatchSize = 50

// Ledger is an append-only record of every bank line this tool has posted, so a
// line is never posted twice even if the plan is regenerated.
type Ledger struct {
	path    string
	Applied map[string]string // line_id -> xero BankTransactionID
}

type ledgerEntry struct {
	LineID    string    `json:"line_id"`
	XeroID    string    `json:"xero_id"`
	PostedAt  time.Time `json:"posted_at"`
	Amount    string    `json:"amount"`
	Date      string    `json:"date"`
	Reference string    `json:"reference,omitempty"`
}

func OpenLedger(path string) (*Ledger, error) {
	l := &Ledger{path: path, Applied: map[string]string{}}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		var e ledgerEntry
		if err := json.Unmarshal(s.Bytes(), &e); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		l.Applied[e.LineID] = e.XeroID
	}
	return l, s.Err()
}

func (l *Ledger) record(e ledgerEntry) error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	l.Applied[e.LineID] = e.XeroID
	return nil
}

// Group is all plan rows for one bank line.
type Group struct {
	LineID string
	Rows   []*PlanRow
}

func GroupPlan(rows []PlanRow) []Group {
	var groups []Group
	index := map[string]int{}
	for i := range rows {
		id := rows[i].LineID
		gi, ok := index[id]
		if !ok {
			gi = len(groups)
			index[id] = gi
			groups = append(groups, Group{LineID: id})
		}
		groups[gi].Rows = append(groups[gi].Rows, &rows[i])
	}
	return groups
}

// Validate checks a ready group is internally consistent before anything is sent.
func (g Group) Validate() error {
	first := g.Rows[0]
	total := decimal.Zero
	for _, r := range g.Rows {
		if r.Status != first.Status {
			return fmt.Errorf("line %s: rows have mixed statuses (%s / %s)", g.LineID, first.Status, r.Status)
		}
		if !r.Amount.Equal(first.Amount) || !r.Date.Equal(first.Date) {
			return fmt.Errorf("line %s: rows disagree on bank amount or date", g.LineID)
		}
		if r.AccountCode == "" {
			return fmt.Errorf("line %s: missing account_code", g.LineID)
		}
		if r.Contact == "" {
			return fmt.Errorf("line %s: missing contact", g.LineID)
		}
		if r.LineAmount.IsNegative() {
			return fmt.Errorf("line %s: line_amount must be positive", g.LineID)
		}
		total = total.Add(r.LineAmount)
	}
	if !total.Equal(first.Amount.Abs()) {
		return fmt.Errorf("line %s: line amounts add up to %s but the bank line is %s", g.LineID, total.StringFixed(2), first.Amount.Abs().StringFixed(2))
	}
	if first.Amount.IsZero() {
		return fmt.Errorf("line %s: zero amount", g.LineID)
	}
	return nil
}

// ToBankTransaction builds an unreconciled spend/receive money transaction. It is
// left unreconciled on purpose: Xero then offers it as the match for the bank feed
// line and reconciling is a single "OK" click, keeping the statement authoritative.
func (g Group) ToBankTransaction(bankAccountID string) xero.BankTransaction {
	first := g.Rows[0]
	t := xero.BankTransaction{
		Type:            xero.BankTransactionReceive,
		Contact:         &xero.Contact{Name: first.Contact},
		BankAccount:     &xero.Account{AccountID: bankAccountID},
		Date:            first.Date.Format("2006-01-02"),
		Reference:       first.BankReference,
		LineAmountTypes: "Inclusive",
	}
	if first.Amount.IsNegative() {
		t.Type = xero.BankTransactionSpend
	}
	if first.OrderID != "" {
		t.Reference = first.OrderID
	}
	t.Reference = truncate(t.Reference, 255)
	for _, r := range g.Rows {
		t.LineItems = append(t.LineItems, xero.LineItem{
			Description: r.Description,
			Quantity:    1,
			UnitAmount:  xero.NewMoney(r.LineAmount),
			AccountCode: r.AccountCode,
			TaxType:     r.TaxType,
		})
	}
	return t
}

type ApplyOptions struct {
	BankAccountID string
	Commit        bool // false = dry run, nothing is sent
	Log           func(format string, args ...any)
}

type ApplyResult struct {
	Posted, Failed, Duplicates, DryRun int
}

// Apply posts every ready group in the plan and updates the rows in place.
func Apply(ctx context.Context, c *xero.Client, ledger *Ledger, rows []PlanRow, opts ApplyOptions) (ApplyResult, error) {
	var res ApplyResult
	var ready []Group

	for _, g := range GroupPlan(rows) {
		if g.Rows[0].Status != StatusReady {
			continue
		}
		if err := g.Validate(); err != nil {
			return res, err
		}
		if id, ok := ledger.Applied[g.LineID]; ok {
			setStatus(g, StatusApplied, "already posted by this tool", id)
			res.Duplicates++
			continue
		}
		ready = append(ready, g)
	}
	if len(ready) == 0 {
		return res, nil
	}

	// Belt and braces: re-check Xero right before posting in case someone keyed the
	// transaction by hand since the plan was built.
	from, to := ready[0].Rows[0].Date, ready[0].Rows[0].Date
	for _, g := range ready {
		d := g.Rows[0].Date
		if d.Before(from) {
			from = d
		}
		if d.After(to) {
			to = d
		}
	}
	existing, err := c.ListBankTransactions(ctx, opts.BankAccountID, from.AddDate(0, 0, -existingMatchDays), to.AddDate(0, 0, existingMatchDays))
	if err != nil {
		return res, err
	}
	claimed := map[string]bool{}
	var toPost []Group
	for _, g := range ready {
		line := StatementLine{Date: g.Rows[0].Date, Amount: g.Rows[0].Amount}
		if id := findExisting(line, existing, claimed); id != "" {
			setStatus(g, StatusExists, "found in Xero at apply time", id)
			res.Duplicates++
			continue
		}
		toPost = append(toPost, g)
	}

	if !opts.Commit {
		for _, g := range toPost {
			t := g.ToBankTransaction(opts.BankAccountID)
			opts.Log("DRY RUN %s %s %8s %-20s %s", t.Date, t.Type, g.Rows[0].Amount.StringFixed(2), t.Contact.Name, t.Reference)
		}
		res.DryRun = len(toPost)
		return res, nil
	}

	for start := 0; start < len(toPost); start += applyBatchSize {
		batch := toPost[start:min(start+applyBatchSize, len(toPost))]
		txns := make([]xero.BankTransaction, len(batch))
		for i, g := range batch {
			txns[i] = g.ToBankTransaction(opts.BankAccountID)
		}

		created, err := c.CreateBankTransactions(ctx, txns)
		if err != nil {
			return res, err
		}
		if len(created) != len(batch) {
			return res, fmt.Errorf("xero returned %d results for %d transactions - check Xero before re-running", len(created), len(batch))
		}

		for i, g := range batch {
			ct := created[i]
			if ct.StatusAttributeString == "ERROR" || ct.BankTransactionID == "" {
				msg := "xero rejected it"
				for _, v := range ct.ValidationErrors {
					msg += ": " + v.Message
				}
				setStatus(g, StatusReview, msg, "")
				opts.Log("FAILED  %s %8s %s", g.Rows[0].Date.Format("2006-01-02"), g.Rows[0].Amount.StringFixed(2), msg)
				res.Failed++
				continue
			}
			if err := ledger.record(ledgerEntry{
				LineID: g.LineID, XeroID: ct.BankTransactionID, PostedAt: time.Now().UTC(),
				Amount: g.Rows[0].Amount.StringFixed(2), Date: g.Rows[0].Date.Format("2006-01-02"), Reference: txns[i].Reference,
			}); err != nil {
				return res, fmt.Errorf("posted %s to xero as %s but could not write ledger: %w", g.LineID, ct.BankTransactionID, err)
			}
			setStatus(g, StatusApplied, "", ct.BankTransactionID)
			opts.Log("POSTED  %s %8s %s", g.Rows[0].Date.Format("2006-01-02"), g.Rows[0].Amount.StringFixed(2), ct.BankTransactionID)
			res.Posted++
		}
	}
	return res, nil
}

func setStatus(g Group, status, reason, xeroID string) {
	for _, r := range g.Rows {
		r.Status, r.Reason, r.XeroID = status, reason, xeroID
	}
}
