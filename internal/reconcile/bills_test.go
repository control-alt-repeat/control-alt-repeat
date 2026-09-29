package reconcile

import (
	"strings"
	"testing"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

func bill(id, number, contact, date, due string) xero.Invoice {
	return xero.Invoice{
		InvoiceID: id, InvoiceNumber: number, Contact: &xero.Contact{Name: contact},
		DateString: date + "T00:00:00", AmountDue: xero.NewMoney(mustDec(due)),
	}
}

func line(id, date, amount, payee string) StatementLine {
	d, _ := time.Parse("2006-01-02", date)
	return StatementLine{ID: id, Date: d, Amount: mustDec(amount), Payee: payee}
}

func TestBillMatching(t *testing.T) {
	rules := &Rules{}
	bills := []xero.Invoice{
		bill("b1", "INV-100", "Royal Mail Group Ltd", "2026-03-01", "54.20"),
		bill("b2", "A-1", "Screwfix", "2026-03-01", "20.00"),
		bill("b3", "A-2", "Toolstation", "2026-03-02", "20.00"),
		bill("b4", "X-9", "Octopus Energy", "2026-01-01", "99.00"),
	}
	lines := []StatementLine{
		line("l1", "2026-03-05", "-54.20", "ROYALMAIL"),             // no name match needed: unique amount
		line("l2", "2026-03-06", "-20.00", "SCREWFIX DIRECT"),       // two £20 bills, supplier name decides
		line("l3", "2026-03-07", "-20.00", "CARD PAYMENT"),          // remaining £20 bill is now unique
		line("l4", "2026-09-01", "-99.00", "OCTOPUS ENERGY"),        // outside the 120 day window
		line("l5", "2026-03-08", "54.20", "REFUND FROM ROYAL MAIL"), // money in never pays a bill
	}

	rows, err := BuildPlan(PlanInput{Lines: lines, Rules: rules, Bills: bills})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]PlanRow{}
	for _, r := range rows {
		got[r.LineID] = r
	}

	want := map[string]string{"l1": "b1", "l2": "b2", "l3": "b3", "l4": "", "l5": ""}
	for id, billID := range want {
		if got[id].BillID != billID {
			t.Errorf("%s: bill %q, want %q (%s)", id, got[id].BillID, billID, got[id].Reason)
		}
	}
	if got["l1"].Status != StatusReady || got["l1"].Contact != "Royal Mail Group Ltd" {
		t.Errorf("l1: %+v", got["l1"])
	}
}

func TestBillAmbiguous(t *testing.T) {
	bills := []xero.Invoice{
		bill("b2", "A-1", "Screwfix", "2026-03-01", "20.00"),
		bill("b3", "A-2", "Toolstation", "2026-03-02", "20.00"),
	}
	rows, err := BuildPlan(PlanInput{Lines: []StatementLine{line("l", "2026-03-06", "-20.00", "CARD PAYMENT")}, Rules: &Rules{}, Bills: bills})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].BillID != "" || rows[0].Status != StatusReview || !strings.Contains(rows[0].Reason, "several unpaid bills") {
		t.Errorf("got %+v", rows[0])
	}
}

func TestToPaymentAndAttachment(t *testing.T) {
	rows, err := BuildPlan(PlanInput{
		Lines: []StatementLine{line("l1", "2026-03-05", "-54.20", "ROYALMAIL")},
		Rules: &Rules{}, Bills: []xero.Invoice{bill("b1", "INV-100", "Royal Mail", "2026-03-01", "54.20")},
	})
	if err != nil {
		t.Fatal(err)
	}
	g := GroupPlan(rows)[0]
	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}
	p := g.ToPayment("bank")
	if p.Invoice.InvoiceID != "b1" || p.Amount.String() != "54.2" || p.Date != "2026-03-05" || p.Account.AccountID != "bank" {
		t.Errorf("payment %+v", p)
	}

	eb := Group{LineID: "x", Rows: []*PlanRow{{Source: "ebay", OrderID: "22-1/2", OrderURL: "https://x", Date: g.Rows[0].Date, Amount: mustDec("-8.99"), LineAmount: mustDec("8.99"), Description: "PSU"}}}
	name, body := eb.OrderAttachment()
	if name != "ebay-order-22-1-2.txt" || !strings.Contains(string(body), "https://x") || !strings.Contains(string(body), "8.99  PSU") {
		t.Errorf("attachment %s:\n%s", name, body)
	}
}
