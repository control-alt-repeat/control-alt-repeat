package reconcile

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/control-alt-repeat/control-alt-repeat/internal/xero"
)

func loadFixtures(t *testing.T) PlanInput {
	t.Helper()

	profiles, err := LoadProfiles("../../configs/reconcile/profiles.json")
	if err != nil {
		t.Fatal(err)
	}
	rules, err := LoadRules("../../configs/reconcile/rules.example.json")
	if err != nil {
		t.Fatal(err)
	}
	sp, err := profiles.Statement("xero")
	if err != nil {
		t.Fatal(err)
	}
	lines, err := LoadStatement("testdata/statement.csv", sp)
	if err != nil {
		t.Fatal(err)
	}
	pp, err := profiles.Purchase("amazon")
	if err != nil {
		t.Fatal(err)
	}
	purchases, err := LoadPurchases("testdata/amazon.csv", pp)
	if err != nil {
		t.Fatal(err)
	}

	return PlanInput{Lines: lines, Rules: rules, Sources: []PurchaseSource{{Profile: pp, Purchases: purchases}}}
}

func byLine(rows []PlanRow) map[string][]PlanRow {
	out := map[string][]PlanRow{}
	for _, r := range rows {
		key := r.Date.Format("0102") + " " + r.Amount.StringFixed(2)
		out[key] = append(out[key], r)
	}
	return out
}

func TestLoadStatement(t *testing.T) {
	in := loadFixtures(t)
	if len(in.Lines) != 9 {
		t.Fatalf("want 9 lines, got %d", len(in.Lines))
	}
	if got := in.Lines[6].Amount.String(); got != "1234.56" {
		t.Errorf("thousands separator: got %s", got)
	}
	if in.Lines[3].ID == in.Lines[4].ID {
		t.Error("identical same-day lines must get different IDs")
	}

	again := loadFixtures(t)
	for i := range in.Lines {
		if in.Lines[i].ID != again.Lines[i].ID {
			t.Fatal("line IDs must be stable across loads")
		}
	}
}

func TestLoadPurchasesExcludesCancelled(t *testing.T) {
	in := loadFixtures(t)
	for _, p := range in.Sources[0].Purchases {
		if p.OrderID == "202-5555555-5555555" {
			t.Fatal("cancelled order should be excluded")
		}
	}
	if len(in.Sources[0].Purchases) != 4 {
		t.Fatalf("want 4 orders, got %d", len(in.Sources[0].Purchases))
	}
}

func TestBuildPlan(t *testing.T) {
	rows, err := BuildPlan(loadFixtures(t))
	if err != nil {
		t.Fatal(err)
	}
	got := byLine(rows)

	cases := []struct {
		line, status, order, account, reason string
	}{
		{"0303 -45.98", StatusReady, "202-1111111-1111111", "429", ""},
		{"0304 -12.50", StatusReady, "202-2222222-2222222", "310", ""}, // split shipment, item rule
		{"0305 -7.99", StatusReady, "202-2222222-2222222", "429", ""},
		{"0307 -3.50", StatusReady, "", "425", ""},
		{"0308 1234.56", StatusReview, "", "200", "manual review"},
		{"0309 -15.00", StatusReview, "", "", "no rule"},
		{"0310 -99.99", StatusReview, "", "429", "no order"},
	}
	for _, c := range cases {
		r := got[c.line]
		if len(r) != 1 {
			t.Errorf("%s: want 1 row, got %d", c.line, len(r))
			continue
		}
		if r[0].Status != c.status || r[0].OrderID != c.order || r[0].AccountCode != c.account {
			t.Errorf("%s: got status=%s order=%s account=%s", c.line, r[0].Status, r[0].OrderID, r[0].AccountCode)
		}
		if c.reason != "" && !strings.Contains(r[0].Reason, c.reason) {
			t.Errorf("%s: reason %q does not mention %q", c.line, r[0].Reason, c.reason)
		}
	}

	// Two £20 charges and two £20 orders: never guess which is which.
	twenty := got["0306 -20.00"]
	if len(twenty) != 2 {
		t.Fatalf("want 2 rows for the £20 lines, got %d", len(twenty))
	}
	for _, r := range twenty {
		if r.Status != StatusReview || !strings.Contains(r.Reason, "ambiguous") {
			t.Errorf("£20 line should be ambiguous, got %s %q", r.Status, r.Reason)
		}
	}
}

func TestBuildPlanSkipsAppliedAndExisting(t *testing.T) {
	in := loadFixtures(t)
	in.Applied = map[string]string{in.Lines[0].ID: "xero-1"}
	in.Existing = []xero.BankTransaction{{
		BankTransactionID: "xero-2",
		Type:              xero.BankTransactionSpend,
		DateString:        "2026-03-08T00:00:00",
		Total:             xero.NewMoney(decimal.RequireFromString("3.50")),
	}}

	rows, err := BuildPlan(in)
	if err != nil {
		t.Fatal(err)
	}
	got := byLine(rows)
	if r := got["0303 -45.98"][0]; r.Status != StatusApplied || r.XeroID != "xero-1" {
		t.Errorf("ledger line: got %s %s", r.Status, r.XeroID)
	}
	if r := got["0307 -3.50"][0]; r.Status != StatusExists || r.XeroID != "xero-2" {
		t.Errorf("existing line: got %s %s", r.Status, r.XeroID)
	}
}

func TestPlanRoundTripAndValidate(t *testing.T) {
	rows, err := BuildPlan(loadFixtures(t))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "plan.csv")
	if err := WritePlan(path, rows); err != nil {
		t.Fatal(err)
	}
	back, err := ReadPlan(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(back) != len(rows) {
		t.Fatalf("round trip lost rows: %d vs %d", len(back), len(rows))
	}

	for _, g := range GroupPlan(back) {
		if g.Rows[0].Status == StatusReady {
			if err := g.Validate(); err != nil {
				t.Errorf("ready group invalid: %v", err)
			}
		}
	}

	// A human editing an amount so it no longer adds up must be caught.
	g := GroupPlan(back)[0]
	g.Rows[0].LineAmount = g.Rows[0].LineAmount.Add(decimal.NewFromInt(1))
	if err := g.Validate(); err == nil {
		t.Error("expected validation error for mismatched total")
	}
}

func TestToBankTransactionJSON(t *testing.T) {
	rows, err := BuildPlan(loadFixtures(t))
	if err != nil {
		t.Fatal(err)
	}
	g := GroupPlan(rows)[0]
	data, err := json.Marshal(g.ToBankTransaction("bank-id"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	for _, want := range []string{`"Type":"SPEND"`, `"UnitAmount":45.98`, `"Reference":"202-1111111-1111111"`, `"IsReconciled":false`, `"Date":"2026-03-03"`, `"LineAmountTypes":"Inclusive"`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %s in %s", want, s)
		}
	}
}

func TestApplyDryRunSkipsLedgerLines(t *testing.T) {
	rows, err := BuildPlan(loadFixtures(t))
	if err != nil {
		t.Fatal(err)
	}
	ledger := &Ledger{path: filepath.Join(t.TempDir(), "l.jsonl"), Applied: map[string]string{}}
	for _, r := range rows {
		if r.Status == StatusReady {
			ledger.Applied[r.LineID] = "done"
		}
	}
	// Every ready line is in the ledger, so Apply must not need Xero at all (nil client).
	res, err := Apply(context.Background(), nil, ledger, rows, ApplyOptions{Log: t.Logf})
	if err != nil {
		t.Fatal(err)
	}
	if res.Duplicates != 4 || res.DryRun != 0 {
		t.Errorf("got %+v", res)
	}
}

func TestParseMoney(t *testing.T) {
	cases := map[string]string{"£1,234.56": "1234.56", "(12.00)": "-12", "12.00 DR": "-12", "5.00CR": "5", "+3": "3", "": "0", "-0.99": "-0.99"}
	for in, want := range cases {
		got, err := parseMoney(in)
		if err != nil || got.String() != want {
			t.Errorf("parseMoney(%q) = %s, %v; want %s", in, got, err, want)
		}
	}
}
