package reconcile

import (
	"strings"
)

// A refund can land a while after the purchase it reverses.
const refundWindowDays = 120

// codeRefunds handles money coming back for something bought earlier. A refund
// carrying the same bank reference as an earlier charge (eBay reuses "eBay O*<order>")
// is coded exactly like that charge, so the two net off in the same account. The
// charge itself is released from review too when the refund covers it in full -
// e.g. a cancelled order that is missing from the marketplace export.
func codeRefunds(lines []StatementLine, rows map[string][]PlanRow, rules *Rules) {
	for _, refund := range lines {
		rr := rows[refund.ID]
		if refund.IsSpend() || len(rr) != 1 || rr[0].Status != StatusReview {
			continue
		}
		ref := normaliseRef(refund.Reference)
		if ref == "" {
			continue
		}

		for _, charge := range lines {
			if !charge.IsSpend() || normaliseRef(charge.Reference) != ref || charge.Date.After(refund.Date) ||
				refund.Date.Sub(charge.Date).Hours() > refundWindowDays*24 || charge.Amount.Abs().LessThan(refund.Amount) {
				continue
			}
			rule := rules.MatchLine(charge)
			if rule == nil || rule.Skip {
				continue
			}

			row := &rows[refund.ID][0]
			row.Status, row.Reason = StatusReady, ""
			row.Source = "refund of " + charge.Date.Format("2006-01-02") + " charge"
			row.AccountCode, row.TaxType = rule.AccountCode, rule.TaxType
			row.Contact = rule.Contact
			if row.Contact == "" {
				row.Contact = charge.Payee
			}
			row.Description = "Refund: " + strings.TrimSpace(refund.Reference)

			// A full refund of a charge still waiting for review settles that charge too.
			cr := rows[charge.ID]
			if len(cr) == 1 && cr[0].Status == StatusReview && charge.Amount.Abs().Equal(refund.Amount) {
				c := &rows[charge.ID][0]
				c.Status, c.Reason = StatusReady, ""
				c.Source = "refunded " + refund.Date.Format("2006-01-02")
				c.AccountCode, c.TaxType, c.Contact = rule.AccountCode, rule.TaxType, row.Contact
				c.Description = "Refunded purchase: " + strings.TrimSpace(charge.Reference)
			}
			break
		}
	}
}

func normaliseRef(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}
