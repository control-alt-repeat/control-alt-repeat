package reconcile

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestLoadEbayXLSX(t *testing.T) {
	profiles, err := LoadProfiles("../../configs/reconcile/profiles.json")
	if err != nil {
		t.Fatal(err)
	}
	pp, err := profiles.Purchase("ebay")
	if err != nil {
		t.Fatal(err)
	}
	purchases, err := LoadPurchases("testdata/ebay.xlsx", pp)
	if err != nil {
		t.Fatal(err)
	}

	byID := map[string]Purchase{}
	for _, p := range purchases {
		byID[p.OrderID] = p
	}
	if _, ok := byID["44-44444-44444"]; ok {
		t.Error("£0 order should be dropped")
	}
	if len(purchases) != 4 {
		t.Fatalf("want 4 orders, got %d", len(purchases))
	}

	cases := map[string][]string{
		"11-11111-11111": {"8.99"},
		"22-22222-22222": {"74.81"},         // discount applied
		"33-33333-33333": {"9.73", "29.91"}, // items already add up
		"55-55555-55555": {"7.5", "7.5"},    // discount spread in proportion
	}
	for id, want := range cases {
		p := byID[id]
		if len(p.Items) != len(want) {
			t.Fatalf("%s: want %d items, got %d", id, len(want), len(p.Items))
		}
		for i, w := range want {
			if got := p.Items[i].Total.String(); got != w {
				t.Errorf("%s item %d: got %s want %s", id, i, got, w)
			}
		}
	}
	if d := byID["11-11111-11111"].Date.Format("2006-01-02"); d != "2026-09-28" {
		t.Errorf("date parsed as %s", d)
	}
	if pp.contactName() != "eBay" {
		t.Errorf("contact = %s", pp.contactName())
	}
}

func TestScaleItemsRoundingPennies(t *testing.T) {
	items := []PurchaseItem{{Total: mustDec("10")}, {Total: mustDec("10")}, {Total: mustDec("10")}}
	out := scaleItems(items, mustDec("10"))
	sum := out[0].Total.Add(out[1].Total).Add(out[2].Total)
	if !sum.Equal(mustDec("10")) {
		t.Errorf("scaled items add up to %s", sum)
	}
}

func TestMatchByOrderIDBeatsAmbiguity(t *testing.T) {
	profiles, err := LoadProfiles("../../configs/reconcile/profiles.json")
	if err != nil {
		t.Fatal(err)
	}
	pp, err := profiles.Purchase("ebay")
	if err != nil {
		t.Fatal(err)
	}
	purchases, err := LoadPurchases("testdata/ebay.xlsx", pp)
	if err != nil {
		t.Fatal(err)
	}
	// Add a second £8.99 order on the same day so amount+date alone is ambiguous.
	dup := purchases[0]
	for _, p := range purchases {
		if p.OrderID == "11-11111-11111" {
			dup = p
		}
	}
	dup.OrderID = "99-99999-99999"
	purchases = append(purchases, dup)

	rules, err := LoadRules("../../configs/reconcile/rules.example.json")
	if err != nil {
		t.Fatal(err)
	}
	line := StatementLine{ID: "x", Date: dup.Date, Amount: mustDec("-8.99"), Payee: "EBAY O*11-11111-11111"}
	rows, err := BuildPlan(PlanInput{Lines: []StatementLine{line}, Rules: rules, Sources: []PurchaseSource{{Profile: pp, Purchases: purchases}}})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Status != StatusReady || rows[0].OrderID != "11-11111-11111" {
		t.Errorf("got %s %s %q", rows[0].Status, rows[0].OrderID, rows[0].Reason)
	}

	line.Payee = "EBAY"
	rows, err = BuildPlan(PlanInput{Lines: []StatementLine{line}, Rules: rules, Sources: []PurchaseSource{{Profile: pp, Purchases: purchases}}})
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Status != StatusReview {
		t.Errorf("without an order number the match must be ambiguous, got %s", rows[0].Status)
	}
}

func mustDec(s string) decimal.Decimal {
	return decimal.RequireFromString(s)
}
