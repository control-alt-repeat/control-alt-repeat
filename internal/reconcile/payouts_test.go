package reconcile

import (
	"strings"
	"testing"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/stripe"
)

var stripeCoding = &PayoutCoding{Contact: "Stripe", PayeePattern: "(?i)stripe", SalesAccount: "200", SalesTax: "OUTPUT2", FeesAccount: "504", FeesTax: "NONE"}

func payout(id, amount string, day int, txns ...stripe.BalanceTransaction) Payout {
	return StripePayout(stripe.Payout{ID: id, Amount: mustDec(amount), Currency: "gbp", ArrivalDate: time.Date(2026, 6, day, 0, 0, 0, 0, time.UTC)}, txns)
}

func bt(typ, amount, fee string) stripe.BalanceTransaction {
	return stripe.BalanceTransaction{Type: typ, Amount: mustDec(amount), Fee: mustDec(fee)}
}

func TestStripePayoutSplit(t *testing.T) {
	// £100 + £50 sales, £20 refund, fees £2.20 + £1.10, payout = 150 - 20 - 3.30 = 126.70
	p := payout("po_1", "126.70", 11, bt("charge", "100", "2.20"), bt("charge", "50", "1.10"), bt("refund", "-20", "0"), bt("payout", "-126.70", "0"))
	if p.Problem != "" {
		t.Fatal(p.Problem)
	}
	want := map[string]string{"sales": "150", "refunds": "-20", "fees": "-3.3"}
	for _, l := range p.Lines {
		if want[l.Kind] != l.Amount.String() {
			t.Errorf("%s = %s", l.Kind, l.Amount)
		}
	}

	rows, err := BuildPlan(PlanInput{
		Lines:   []StatementLine{line("l1", "2026-06-11", "126.70", "Stripe Payments UK Ltd")},
		Rules:   &Rules{Stripe: stripeCoding},
		Payouts: []Payout{p},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("want 3 rows, got %d", len(rows))
	}
	g := GroupPlan(rows)[0]
	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, r := range rows {
		if r.Status != StatusReady {
			t.Errorf("%s: %s %s", r.Description, r.Status, r.Reason)
		}
		if strings.HasPrefix(r.Description, "Stripe fees") != (r.AccountCode == "504") {
			t.Errorf("%s coded to %s", r.Description, r.AccountCode)
		}
	}
	bt := g.ToBankTransaction("bank")
	if bt.Type != "RECEIVE" || len(bt.LineItems) != 3 {
		t.Errorf("%+v", bt)
	}
}

func TestStripePayoutProblems(t *testing.T) {
	if p := payout("po_2", "100", 11, bt("charge", "100", "1")); !strings.Contains(p.Problem, "adds up to 99.00") {
		t.Errorf("mismatch not caught: %q", p.Problem)
	}
	if p := payout("po_3", "50", 11, bt("charge", "100", "0"), bt("dispute", "-50", "0")); !strings.Contains(p.Problem, "dispute") {
		t.Errorf("dispute not flagged: %q", p.Problem)
	}

	rows, err := BuildPlan(PlanInput{
		Lines:   []StatementLine{line("l1", "2026-06-11", "99.00", "Stripe"), line("l2", "2026-06-20", "40.00", "Stripe")},
		Rules:   &Rules{Stripe: stripeCoding},
		Payouts: []Payout{payout("a", "99", 10, bt("charge", "99", "0")), payout("b", "99", 12, bt("charge", "99", "0"))},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Status != StatusReview || !strings.Contains(rows[0].Reason, "several") {
		t.Errorf("ambiguous payout: %+v", rows[0])
	}
	if rows[1].Status != StatusReview || !strings.Contains(rows[1].Reason, "no Stripe payout") {
		t.Errorf("missing payout: %+v", rows[1])
	}
}
