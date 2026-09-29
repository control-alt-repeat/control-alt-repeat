package reconcile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

// Plan statuses. A human edits the plan CSV to move rows from review to ready (or skip).
const (
	StatusReady   = "ready"   // fully coded, `apply` will post it
	StatusReview  = "review"  // needs a human decision, see reason
	StatusSkip    = "skip"    // deliberately not posted by this tool
	StatusExists  = "exists"  // a matching transaction is already in Xero
	StatusApplied = "applied" // posted to Xero, see xero_id
)

// PlanRow is one Xero line item. A bank line split across several purchase items
// has several rows sharing the same line_id.
type PlanRow struct {
	LineID          string
	Date            time.Time
	Amount          decimal.Decimal // the whole bank line, signed
	BankPayee       string
	BankDescription string
	BankReference   string
	Status          string
	Reason          string
	Source          string
	OrderID         string
	Contact         string
	AccountCode     string
	TaxType         string
	Description     string
	LineAmount      decimal.Decimal // this row's share, always positive
	XeroID          string
	BillID          string // set when the line pays an existing Xero bill
	OrderURL        string
}

var planHeader = []string{
	"line_id", "date", "amount", "bank_payee", "bank_description", "bank_reference",
	"status", "reason", "source", "order_id", "contact", "account_code", "tax_type",
	"description", "line_amount", "xero_id", "bill_id", "order_url",
}

// Columns added after the first release, so older plan files still load.
var optionalPlanColumns = map[string]bool{"bill_id": true, "order_url": true}

// PurchaseSource is one loaded marketplace export plus how to match it to the bank.
type PurchaseSource struct {
	Profile   PurchaseProfile
	Purchases []Purchase
}

type PlanInput struct {
	Lines    []StatementLine
	Rules    *Rules
	Sources  []PurchaseSource
	Existing []xero.BankTransaction // optional: already in Xero for this bank account
	Bills    []xero.Invoice         // optional: unpaid supplier bills
	Applied  map[string]string      // optional: line_id -> xero id from the ledger
}

const existingMatchDays = 3

func BuildPlan(in PlanInput) ([]PlanRow, error) {
	lines := append([]StatementLine(nil), in.Lines...)
	sort.SliceStable(lines, func(i, j int) bool { return lines[i].Date.Before(lines[j].Date) })

	rows := map[string][]PlanRow{}
	pending := []StatementLine{}

	claimedExisting := map[string]bool{}
	for _, l := range lines {
		if id, ok := in.Applied[l.ID]; ok {
			rows[l.ID] = []PlanRow{baseRow(l, StatusApplied, "already posted by this tool")}
			rows[l.ID][0].XeroID = id
			continue
		}
		if id := findExisting(l, in.Existing, claimedExisting); id != "" {
			row := baseRow(l, StatusExists, "a transaction with this amount and date is already in Xero")
			row.XeroID = id
			rows[l.ID] = []PlanRow{row}
			continue
		}
		pending = append(pending, l)
	}

	billMatches, billReasons := matchBills(pending, in.Bills)
	unbilled := pending[:0:0]
	for _, l := range pending {
		if b, ok := billMatches[l.ID]; ok {
			rows[l.ID] = []PlanRow{billRow(l, b)}
			continue
		}
		unbilled = append(unbilled, l)
	}

	m, err := newPurchaseMatcher(in.Sources)
	if err != nil {
		return nil, err
	}
	matches := m.match(unbilled)
	for id, reason := range billReasons {
		if _, ok := matches[id]; !ok {
			m.reasons[id] = reason
		}
	}

	for _, l := range unbilled {
		if r, ok := matches[l.ID]; ok {
			rows[l.ID] = purchaseRows(l, r, in.Rules)
			continue
		}
		rows[l.ID] = []PlanRow{ruleRow(l, in.Rules, m.reasons[l.ID])}
	}

	codeRefunds(lines, rows, in.Rules)

	out := make([]PlanRow, 0, len(lines))
	for _, l := range lines {
		out = append(out, rows[l.ID]...)
	}
	return out, nil
}

func baseRow(l StatementLine, status, reason string) PlanRow {
	return PlanRow{
		LineID:          l.ID,
		Date:            l.Date,
		Amount:          l.Amount,
		BankPayee:       l.Payee,
		BankDescription: l.Description,
		BankReference:   l.Reference,
		Status:          status,
		Reason:          reason,
		LineAmount:      l.Amount.Abs(),
	}
}

func findExisting(l StatementLine, existing []xero.BankTransaction, claimed map[string]bool) string {
	for _, e := range existing {
		if claimed[e.BankTransactionID] || !e.SignedTotal().Equal(l.Amount) {
			continue
		}
		d, err := e.DateValue()
		if err != nil || absDays(d.Sub(l.Date)) > existingMatchDays {
			continue
		}
		claimed[e.BankTransactionID] = true
		return e.BankTransactionID
	}
	return ""
}

func ruleRow(l StatementLine, rules *Rules, purchaseReason string) PlanRow {
	row := baseRow(l, StatusReview, "no rule matches - add one to the rules file or code this row by hand")
	row.Contact = l.Payee
	row.Description = strings.TrimSpace(strings.Join([]string{l.Payee, l.Description}, " "))

	rule := rules.MatchLine(l)
	if rule != nil {
		row.Source = "rule: " + rule.Name
		row.AccountCode = rule.AccountCode
		row.TaxType = rule.TaxType
		if rule.Contact != "" {
			row.Contact = rule.Contact
		}
		if rule.Narrative != "" {
			row.Description = rule.Narrative
		}
		row.Status = StatusReady
		row.Reason = ""
		if rule.Review {
			row.Status = StatusReview
			row.Reason = "rule asks for manual review"
		}
		if rule.Skip {
			row.Status = StatusSkip
			row.Reason = "handled elsewhere (rule: " + rule.Name + ")"
			return row
		}
	}

	// A marketplace charge that could not be tied to an order is worth a look even
	// when a rule would code it - the order export may be incomplete.
	if purchaseReason != "" {
		row.Status = StatusReview
		row.Reason = purchaseReason
	}
	return row
}

func purchaseRows(l StatementLine, r purchaseMatch, rules *Rules) []PlanRow {
	lineRule := rules.MatchLine(l)

	contact := r.source.contactName()
	if lineRule != nil && lineRule.Contact != "" {
		contact = lineRule.Contact
	}

	var out []PlanRow
	for _, item := range r.items {
		row := baseRow(l, StatusReady, "")
		row.Source = r.source.Source
		row.OrderID = r.purchase.OrderID
		row.OrderURL = r.purchase.URL
		row.Contact = contact
		row.Description = truncate(fmt.Sprintf("%s %s: %s", row.Contact, r.purchase.OrderID, item.Description), 4000)
		if item.Quantity > 1 {
			row.Description = truncate(fmt.Sprintf("%s %s: %dx %s", row.Contact, r.purchase.OrderID, item.Quantity, item.Description), 4000)
		}
		row.LineAmount = item.Total

		switch rule := rules.MatchItem(l, item); {
		case rule != nil:
			row.AccountCode, row.TaxType = rule.AccountCode, rule.TaxType
			if rule.Review {
				row.Status, row.Reason = StatusReview, "item rule "+rule.Name+" asks for manual review"
			}
		case lineRule != nil:
			row.AccountCode, row.TaxType = lineRule.AccountCode, lineRule.TaxType
		default:
			row.AccountCode, row.TaxType = rules.DefaultPurchaseAccountCode, rules.DefaultPurchaseTaxType
		}
		if row.AccountCode == "" {
			row.Status, row.Reason = StatusReview, "no account code for this item"
		}
		out = append(out, row)
	}
	return out
}

// purchaseMatcher ties card charges to marketplace orders by exact amount within a
// date window. It only accepts unambiguous matches: when two orders could explain
// the same charge, the line goes to review rather than guessing.
type purchaseMatcher struct {
	sources  []PurchaseSource
	payeeRes []*regexp.Regexp
	reasons  map[string]string
}

type purchaseMatch struct {
	source   PurchaseProfile
	purchase Purchase
	items    []PurchaseItem
	key      string // what the match consumes, so nothing is used twice
}

func newPurchaseMatcher(sources []PurchaseSource) (*purchaseMatcher, error) {
	m := &purchaseMatcher{sources: sources, reasons: map[string]string{}}
	for _, s := range sources {
		re, err := s.Profile.payeeMatcher()
		if err != nil {
			return nil, fmt.Errorf("profile %q bank_payee_pattern: %w", s.Profile.Name, err)
		}
		m.payeeRes = append(m.payeeRes, re)
	}
	return m, nil
}

func (m *purchaseMatcher) candidates(l StatementLine, used map[string]bool) (relevant bool, found []purchaseMatch) {
	if !l.IsSpend() {
		return false, nil
	}
	want := l.Amount.Abs()

	for si, s := range m.sources {
		if !m.payeeRes[si].MatchString(l.Text()) {
			continue
		}
		relevant = true
		before, after := s.Profile.DaysBefore, s.Profile.DaysAfter
		if before == 0 && after == 0 {
			before, after = 2, 14
		}

		// An order number printed on the statement line beats any amount/date inference.
		if pm, ok := matchByOrderID(l, s, want, used); ok {
			found = append(found, pm)
			continue
		}

		var orderLevel, itemLevel []purchaseMatch
		for _, p := range s.Purchases {
			if l.Date.Before(p.Date.AddDate(0, 0, -before)) || l.Date.After(p.Date.AddDate(0, 0, after)) {
				continue
			}
			orderKey := s.Profile.Source + "|" + p.OrderID
			if used[orderKey] {
				continue
			}
			if p.Total().Equal(want) && !used[orderKey+"#partial"] {
				orderLevel = append(orderLevel, purchaseMatch{source: s.Profile, purchase: p, items: p.Items, key: orderKey})
				continue
			}
			// Multi-item orders are often charged per shipment, i.e. per item.
			if len(p.Items) > 1 {
				for ii, item := range p.Items {
					itemKey := fmt.Sprintf("%s#%d", orderKey, ii)
					if !used[itemKey] && item.Total.Equal(want) {
						itemLevel = append(itemLevel, purchaseMatch{source: s.Profile, purchase: p, items: []PurchaseItem{item}, key: itemKey})
					}
				}
			}
		}
		if len(orderLevel) > 0 {
			found = append(found, orderLevel...)
		} else {
			found = append(found, itemLevel...)
		}
	}
	return relevant, found
}

func matchByOrderID(l StatementLine, s PurchaseSource, want decimal.Decimal, used map[string]bool) (purchaseMatch, bool) {
	text := l.Text()
	for _, p := range s.Purchases {
		if len(p.OrderID) < 6 || !strings.Contains(text, p.OrderID) {
			continue
		}
		orderKey := s.Profile.Source + "|" + p.OrderID
		if !used[orderKey] && !used[orderKey+"#partial"] && p.Total().Equal(want) {
			return purchaseMatch{source: s.Profile, purchase: p, items: p.Items, key: orderKey}, true
		}
	}
	return purchaseMatch{}, false
}

func (m *purchaseMatcher) match(lines []StatementLine) map[string]purchaseMatch {
	used := map[string]bool{}
	result := map[string]purchaseMatch{}

	// Resolve unique matches first; each one removes an order from contention,
	// which can make other lines unique. Repeat until nothing changes.
	for progress := true; progress; {
		progress = false
		for _, l := range lines {
			if _, done := result[l.ID]; done {
				continue
			}
			if _, found := m.candidates(l, used); len(found) == 1 {
				result[l.ID] = found[0]
				markUsed(used, found[0])
				progress = true
			}
		}
	}

	for _, l := range lines {
		if _, done := result[l.ID]; done {
			continue
		}
		relevant, found := m.candidates(l, used)
		switch {
		case len(found) > 1:
			ids := make([]string, 0, len(found))
			for _, f := range found {
				ids = append(ids, f.purchase.OrderID)
			}
			m.reasons[l.ID] = "ambiguous: could be orders " + strings.Join(ids, ", ")
		case relevant:
			m.reasons[l.ID] = "looks like a marketplace charge but no order in the export matches amount and date"
		}
	}
	return result
}

func markUsed(used map[string]bool, pm purchaseMatch) {
	used[pm.key] = true
	// Consuming a whole order also consumes its items, and vice versa once every item is taken.
	orderKey := pm.source.Source + "|" + pm.purchase.OrderID
	if pm.key == orderKey {
		for i := range pm.purchase.Items {
			used[fmt.Sprintf("%s#%d", orderKey, i)] = true
		}
		return
	}
	used[orderKey+"#partial"] = true
	for i := range pm.purchase.Items {
		if !used[fmt.Sprintf("%s#%d", orderKey, i)] {
			return
		}
	}
	used[orderKey] = true
}

func absDays(d time.Duration) int {
	days := int(d.Hours() / 24)
	if days < 0 {
		return -days
	}
	return days
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func WritePlan(path string, rows []PlanRow) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(planHeader); err != nil {
		return err
	}
	for _, r := range rows {
		rec := []string{
			r.LineID, r.Date.Format("2006-01-02"), r.Amount.StringFixed(2), r.BankPayee, r.BankDescription, r.BankReference,
			r.Status, r.Reason, r.Source, r.OrderID, r.Contact, r.AccountCode, r.TaxType,
			r.Description, r.LineAmount.StringFixed(2), r.XeroID, r.BillID, r.OrderURL,
		}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func ReadPlan(path string) ([]PlanRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(stripBOM(f))
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	col := map[string]int{}
	for i, h := range header {
		col[normaliseHeader(h)] = i
	}
	for _, h := range planHeader {
		if _, ok := col[h]; !ok && !optionalPlanColumns[h] {
			return nil, fmt.Errorf("%s: missing column %q", path, h)
		}
	}

	var rows []PlanRow
	for n := 2; ; n++ {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			return rows, nil
		}
		if err != nil {
			return nil, err
		}
		get := func(name string) string {
			i, ok := col[name]
			if !ok || i >= len(rec) {
				return ""
			}
			return strings.TrimSpace(rec[i])
		}

		date, err := time.Parse("2006-01-02", get("date"))
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", path, n, err)
		}
		amount, err := parseMoney(get("amount"))
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", path, n, err)
		}
		lineAmount, err := parseMoney(get("line_amount"))
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", path, n, err)
		}
		rows = append(rows, PlanRow{
			LineID: get("line_id"), Date: date, Amount: amount,
			BankPayee: get("bank_payee"), BankDescription: get("bank_description"), BankReference: get("bank_reference"),
			Status: strings.ToLower(get("status")), Reason: get("reason"), Source: get("source"), OrderID: get("order_id"),
			Contact: get("contact"), AccountCode: get("account_code"), TaxType: get("tax_type"),
			Description: get("description"), LineAmount: lineAmount, XeroID: get("xero_id"),
			BillID: get("bill_id"), OrderURL: get("order_url"),
		})
	}
}
