package reconcile

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		if r.BillID != "" && len(g.Rows) > 1 {
			return fmt.Errorf("line %s: a bill payment must be a single row", g.LineID)
		}
		if r.AccountCode == "" && r.BillID == "" {
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
	existing, err := c.ListAccountActivity(ctx, opts.BankAccountID, from.AddDate(0, 0, -existingMatchDays), to.AddDate(0, 0, existingMatchDays))
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
			first := g.Rows[0]
			kind := "SPEND"
			if first.BillID != "" {
				kind = "PAY BILL"
			} else if !first.Amount.IsNegative() {
				kind = "RECEIVE"
			}
			opts.Log("DRY RUN %s %-8s %8s %-20s %s", first.Date.Format("2006-01-02"), kind, first.Amount.StringFixed(2), first.Contact, first.OrderID)
		}
		res.DryRun = len(toPost)
		return res, nil
	}

	var bills, txns []Group
	for _, g := range toPost {
		if g.Rows[0].BillID != "" {
			bills = append(bills, g)
		} else {
			txns = append(txns, g)
		}
	}

	for start := 0; start < len(bills); start += applyBatchSize {
		batch := bills[start:min(start+applyBatchSize, len(bills))]
		payments := make([]xero.Payment, len(batch))
		for i, g := range batch {
			payments[i] = g.ToPayment(opts.BankAccountID)
		}
		created, err := c.CreatePayments(ctx, payments)
		if err != nil {
			return res, err
		}
		if len(created) != len(batch) {
			return res, fmt.Errorf("xero returned %d results for %d payments - check Xero before re-running", len(created), len(batch))
		}
		for i, g := range batch {
			if err := settle(g, created[i].PaymentID, created[i].StatusAttributeString, created[i].ValidationErrors, payments[i].Reference, ledger, opts, &res); err != nil {
				return res, err
			}
		}
	}

	for start := 0; start < len(txns); start += applyBatchSize {
		batch := txns[start:min(start+applyBatchSize, len(txns))]
		bankTxns := make([]xero.BankTransaction, len(batch))
		for i, g := range batch {
			bankTxns[i] = g.ToBankTransaction(opts.BankAccountID)
		}

		created, err := c.CreateBankTransactions(ctx, bankTxns)
		if err != nil {
			return res, err
		}
		if len(created) != len(batch) {
			return res, fmt.Errorf("xero returned %d results for %d transactions - check Xero before re-running", len(created), len(batch))
		}

		for i, g := range batch {
			ct := created[i]
			if err := settle(g, ct.BankTransactionID, ct.StatusAttributeString, ct.ValidationErrors, bankTxns[i].Reference, ledger, opts, &res); err != nil {
				return res, err
			}
			if g.Rows[0].Status == StatusApplied && g.Rows[0].OrderID != "" {
				name, body := g.OrderAttachment()
				// The transaction is already posted and recorded, so a failed upload
				// is reported but does not stop the run.
				if err := c.AttachToBankTransaction(ctx, ct.BankTransactionID, name, "text/plain", body); err != nil {
					opts.Log("WARNING could not attach order details to %s: %v", ct.BankTransactionID, err)
				}
			}
		}
	}
	return res, nil
}

// settle records the outcome of one posted group in the ledger and the plan.
func settle(g Group, xeroID, status string, errs []xero.ValidationError, reference string, ledger *Ledger, opts ApplyOptions, res *ApplyResult) error {
	first := g.Rows[0]
	if status == "ERROR" || xeroID == "" {
		msg := "xero rejected it"
		for _, v := range errs {
			msg += ": " + v.Message
		}
		setStatus(g, StatusReview, msg, "")
		opts.Log("FAILED  %s %8s %s", first.Date.Format("2006-01-02"), first.Amount.StringFixed(2), msg)
		res.Failed++
		return nil
	}
	if err := ledger.record(ledgerEntry{
		LineID: g.LineID, XeroID: xeroID, PostedAt: time.Now().UTC(),
		Amount: first.Amount.StringFixed(2), Date: first.Date.Format("2006-01-02"), Reference: reference,
	}); err != nil {
		return fmt.Errorf("posted %s to xero as %s but could not write ledger: %w", g.LineID, xeroID, err)
	}
	setStatus(g, StatusApplied, "", xeroID)
	opts.Log("POSTED  %s %8s %s", first.Date.Format("2006-01-02"), first.Amount.StringFixed(2), xeroID)
	res.Posted++
	return nil
}

// ToPayment pays an existing bill from the bank account, left unreconciled so it
// lines up with the statement line in Xero's Reconcile tab.
func (g Group) ToPayment(bankAccountID string) xero.Payment {
	first := g.Rows[0]
	ref := first.BankReference
	if ref == "" {
		ref = first.BankPayee
	}
	return xero.Payment{
		Invoice:   &xero.Invoice{InvoiceID: first.BillID},
		Account:   &xero.Account{AccountID: bankAccountID},
		Date:      first.Date.Format("2006-01-02"),
		Amount:    xero.NewMoney(first.Amount.Abs()),
		Reference: truncate(ref, 255),
	}
}

// OrderAttachment is a plain-text record of the marketplace order behind a
// transaction, kept in Xero as evidence for VAT and margin scheme checks.
func (g Group) OrderAttachment() (string, []byte) {
	first := g.Rows[0]
	var b strings.Builder
	fmt.Fprintf(&b, "Source:     %s\n", first.Source)
	fmt.Fprintf(&b, "Order:      %s\n", first.OrderID)
	if first.OrderURL != "" {
		fmt.Fprintf(&b, "Order link: %s\n", first.OrderURL)
	}
	fmt.Fprintf(&b, "Paid:       %s on %s\n", first.Amount.Abs().StringFixed(2), first.Date.Format("2006-01-02"))
	fmt.Fprintf(&b, "Bank line:  %s\n\nItems:\n", strings.TrimSpace(first.BankPayee+" "+first.BankDescription))
	for _, r := range g.Rows {
		fmt.Fprintf(&b, "  %8s  %s\n", r.LineAmount.StringFixed(2), r.Description)
	}
	name := strings.NewReplacer("/", "-", "\\", "-", " ", "_").Replace(fmt.Sprintf("%s-order-%s.txt", first.Source, first.OrderID))
	return name, []byte(b.String())
}

func setStatus(g Group, status, reason, xeroID string) {
	for _, r := range g.Rows {
		r.Status, r.Reason, r.XeroID = status, reason, xeroID
	}
}
