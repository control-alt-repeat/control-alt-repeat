package reconcile

import (
	"sort"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

// A bill can be paid a little before its date (card pre-auth) or months after.
const (
	billDaysBefore = 7
	billDaysAfter  = 120
)

// matchBills pairs spend lines with unpaid supplier bills already entered in Xero.
// Paying the bill is always preferable to creating a new spend transaction, which
// would count the cost twice. Like order matching, it never guesses between equals.
func matchBills(lines []StatementLine, bills []xero.Invoice) (map[string]xero.Invoice, map[string]string) {
	matched := map[string]xero.Invoice{}
	reasons := map[string]string{}
	used := map[string]bool{}

	candidates := func(l StatementLine) []xero.Invoice {
		if !l.IsSpend() {
			return nil
		}
		want := l.Amount.Abs()
		text := strings.ToLower(l.Text())

		var inWindow, named, referenced []xero.Invoice
		for _, b := range bills {
			if used[b.InvoiceID] || b.AmountDue == nil || !b.AmountDue.Equal(want) {
				continue
			}
			d, err := b.DateValue()
			if err != nil || l.Date.Before(d.AddDate(0, 0, -billDaysBefore)) || l.Date.After(d.AddDate(0, 0, billDaysAfter)) {
				continue
			}
			inWindow = append(inWindow, b)
			if mentionsReference(text, b) {
				referenced = append(referenced, b)
			}
			if b.Contact != nil && mentionsContact(text, b.Contact.Name) {
				named = append(named, b)
			}
		}
		// Strongest evidence first: a bill number on the statement, then the supplier's name.
		switch {
		case len(referenced) > 0:
			return referenced
		case len(named) > 0:
			return named
		}
		return inWindow
	}

	for progress := true; progress; {
		progress = false
		for _, l := range lines {
			if _, done := matched[l.ID]; done {
				continue
			}
			if c := candidates(l); len(c) == 1 {
				matched[l.ID] = c[0]
				used[c[0].InvoiceID] = true
				progress = true
			}
		}
	}

	for _, l := range lines {
		if _, done := matched[l.ID]; done {
			continue
		}
		if c := candidates(l); len(c) > 1 {
			names := make([]string, 0, len(c))
			for _, b := range c {
				names = append(names, billLabel(b))
			}
			sort.Strings(names)
			reasons[l.ID] = "matches several unpaid bills: " + strings.Join(names, ", ")
		}
	}
	return matched, reasons
}

func mentionsReference(text string, b xero.Invoice) bool {
	for _, ref := range []string{b.InvoiceNumber, b.Reference} {
		ref = strings.ToLower(strings.TrimSpace(ref))
		if len(ref) >= 4 && strings.Contains(text, ref) {
			return true
		}
	}
	return false
}

// mentionsContact checks for the first meaningful word of the supplier name, since
// banks truncate and mangle payee names ("ROYAL MAIL GROUP LTD" -> "ROYALMAIL").
func mentionsContact(text, name string) bool {
	for _, word := range strings.Fields(strings.ToLower(name)) {
		word = strings.Trim(word, ".,()&-")
		if len(word) < 3 || word == "the" || word == "ltd" || word == "limited" {
			continue
		}
		return strings.Contains(text, word)
	}
	return false
}

func billLabel(b xero.Invoice) string {
	label := b.InvoiceNumber
	if label == "" {
		label = b.Reference
	}
	if b.Contact != nil {
		label = strings.TrimSpace(b.Contact.Name + " " + label)
	}
	return label
}

func billRow(l StatementLine, b xero.Invoice) PlanRow {
	row := baseRow(l, StatusReady, "")
	row.Source = "bill"
	row.OrderID = b.InvoiceNumber
	row.BillID = b.InvoiceID
	if b.Contact != nil {
		row.Contact = b.Contact.Name
	}
	row.Description = "Payment of bill " + billLabel(b)
	return row
}
