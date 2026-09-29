package reconcile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRefundsCodedLikeTheirCharge(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, []byte(`{"rules":[{"name":"ebay","text":"(?i)ebay","direction":"out","contact":"eBay","account_code":"310","tax_type":"NONE"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	rules, err := LoadRules(path)
	if err != nil {
		t.Fatal(err)
	}
	withRef := func(l StatementLine, ref string) StatementLine { l.Reference = ref; return l }
	lines := []StatementLine{
		withRef(line("c1", "2026-06-05", "-13.49", "eBay"), "eBay O*03-1"), // cancelled, fully refunded
		withRef(line("r1", "2026-06-16", "13.49", "eBay"), "eBay O*03-1"),
		withRef(line("c2", "2026-06-05", "-50.00", "eBay"), "eBay O*04-2"), // partial refund
		withRef(line("r2", "2026-06-20", "10.00", "eBay"), "eBay O*04-2"),
		withRef(line("r3", "2026-06-20", "5.00", "eBay"), "eBay O*99-9"),  // no matching charge
		withRef(line("c4", "2026-06-25", "-7.00", "eBay"), "eBay O*05-5"), // refund before charge: not a refund
		withRef(line("r4", "2026-06-20", "7.00", "eBay"), "eBay O*05-5"),
	}
	rows, err := BuildPlan(PlanInput{Lines: lines, Rules: rules})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]PlanRow{}
	for _, r := range rows {
		got[r.LineID] = r
	}

	for _, id := range []string{"c1", "r1", "r2"} {
		if got[id].Status != StatusReady || got[id].AccountCode != "310" || got[id].TaxType != "NONE" || got[id].Contact != "eBay" {
			t.Errorf("%s: %+v", id, got[id])
		}
	}
	// A partially refunded charge still needs its order to itemise it.
	if got["c2"].AccountCode != "310" || got["c2"].Source != "rule: ebay" {
		t.Errorf("c2 should keep its own coding: %+v", got["c2"])
	}
	for _, id := range []string{"r3", "r4"} {
		if got[id].Status != StatusReview {
			t.Errorf("%s should need review, got %s", id, got[id].Status)
		}
	}
	if g := GroupPlan(rows); len(g) != len(lines) {
		t.Fatalf("groups %d", len(g))
	}
	for _, g := range GroupPlan(rows) {
		if g.Rows[0].Status == StatusReady {
			if err := g.Validate(); err != nil {
				t.Error(err)
			}
		}
	}
}
