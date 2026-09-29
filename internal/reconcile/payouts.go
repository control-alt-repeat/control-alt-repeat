package reconcile

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/control-alt-repeat/control-alt-repeat/internal/stripe"
)

// PayoutCoding says how a payment processor's payouts are booked. It lives in the
// rules file, e.g. "stripe": {"sales_account": "200", ...}.
type PayoutCoding struct {
	Contact      string `json:"contact"`
	PayeePattern string `json:"payee_pattern"`
	SalesAccount string `json:"sales_account"`
	SalesTax     string `json:"sales_tax,omitempty"`
	FeesAccount  string `json:"fees_account"`
	FeesTax      string `json:"fees_tax,omitempty"`
}

// Payout is money a processor paid into the bank, broken down by what made it up.
// Line amounts are signed and add up exactly to Amount.
type Payout struct {
	Source string
	ID     string
	Amount decimal.Decimal
	Date   time.Time
	Lines  []PayoutLine
	// Problem is set when the breakdown cannot be trusted, which sends the line to review.
	Problem string
}

type PayoutLine struct {
	Kind        string // sales, refunds, fees or other
	Description string
	Amount      decimal.Decimal
}

const payoutMatchDays = 4

// StripePayout summarises a Stripe payout's balance transactions into gross sales,
// refunds and fees, checking they reconcile to the amount paid out.
func StripePayout(p stripe.Payout, txns []stripe.BalanceTransaction) Payout {
	sales, refunds, fees, other := decimal.Zero, decimal.Zero, decimal.Zero, decimal.Zero
	var otherTypes []string

	for _, t := range txns {
		fees = fees.Sub(t.Fee)
		switch t.Type {
		case "payout":
			continue // the payout itself
		case "charge", "payment":
			sales = sales.Add(t.Amount)
		case "refund", "payment_refund":
			refunds = refunds.Add(t.Amount)
		case "stripe_fee", "stripe_fx_fee", "tax_fee":
			fees = fees.Add(t.Amount)
		default:
			other = other.Add(t.Amount)
			otherTypes = append(otherTypes, t.Type)
		}
	}

	out := Payout{Source: "stripe", ID: p.ID, Amount: p.Amount, Date: p.ArrivalDate}
	for _, l := range []PayoutLine{
		{"sales", "Stripe sales " + p.ID, sales},
		{"refunds", "Stripe refunds " + p.ID, refunds},
		{"fees", "Stripe fees " + p.ID, fees},
		{"other", "Stripe other (" + strings.Join(otherTypes, ", ") + ") " + p.ID, other},
	} {
		if !l.Amount.IsZero() {
			out.Lines = append(out.Lines, l)
		}
	}

	total := sales.Add(refunds).Add(fees).Add(other)
	switch {
	case strings.ToLower(p.Currency) != "gbp":
		out.Problem = "payout currency is " + p.Currency
	case !total.Equal(p.Amount):
		out.Problem = fmt.Sprintf("stripe breakdown adds up to %s, payout is %s", total.StringFixed(2), p.Amount.StringFixed(2))
	case !other.IsZero():
		out.Problem = "payout includes " + strings.Join(otherTypes, ", ") + " - check how to code it"
	}
	return out
}

// matchPayouts pairs money-in lines with processor payouts of the same amount
// arriving within a few days. Two payouts that could fit send the line to review.
func matchPayouts(lines []StatementLine, payouts []Payout, coding *PayoutCoding) (map[string]Payout, map[string]string, error) {
	matched, reasons := map[string]Payout{}, map[string]string{}
	if coding == nil || len(payouts) == 0 {
		return matched, reasons, nil
	}
	payee, err := regexp.Compile(coding.PayeePattern)
	if err != nil {
		return nil, nil, fmt.Errorf("payout payee_pattern: %w", err)
	}

	used := map[string]bool{}
	for _, l := range lines {
		if l.IsSpend() || !payee.MatchString(l.Text()) {
			continue
		}
		var found []Payout
		for _, p := range payouts {
			if !used[p.ID] && p.Amount.Equal(l.Amount) && absDays(p.Date.Sub(l.Date)) <= payoutMatchDays {
				found = append(found, p)
			}
		}
		switch len(found) {
		case 1:
			matched[l.ID] = found[0]
			used[found[0].ID] = true
		case 0:
			reasons[l.ID] = "no Stripe payout of this amount arrived around this date"
		default:
			reasons[l.ID] = "several Stripe payouts could match"
		}
	}
	return matched, reasons, nil
}

func payoutRows(l StatementLine, p Payout, coding *PayoutCoding) []PlanRow {
	var out []PlanRow
	for _, pl := range p.Lines {
		row := baseRow(l, StatusReady, "")
		row.Source = p.Source + " payout"
		row.OrderID = p.ID
		row.Contact = coding.Contact
		row.Description = pl.Description
		row.LineAmount = pl.Amount
		switch pl.Kind {
		case "fees":
			row.AccountCode, row.TaxType = coding.FeesAccount, coding.FeesTax
		default:
			row.AccountCode, row.TaxType = coding.SalesAccount, coding.SalesTax
		}
		if p.Problem != "" {
			row.Status, row.Reason = StatusReview, p.Problem
		}
		out = append(out, row)
	}
	return out
}
