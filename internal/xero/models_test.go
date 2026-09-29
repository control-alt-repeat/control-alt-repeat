package xero

import (
	"encoding/json"
	"testing"
)

func TestBankTransactionFromXeroJSON(t *testing.T) {
	body := `{"BankTransactions":[{"BankTransactionID":"abc","Type":"SPEND","Total":45.98,"IsReconciled":false,
		"Date":"/Date(1772496000000+0000)/","DateString":"2026-03-03T00:00:00"}]}`

	var resp struct{ BankTransactions []BankTransaction }
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatal(err)
	}
	bt := resp.BankTransactions[0]
	if got := bt.SignedTotal().String(); got != "-45.98" {
		t.Errorf("signed total = %s", got)
	}
	d, err := bt.DateValue()
	if err != nil || d.Format("2006-01-02") != "2026-03-03" {
		t.Errorf("date = %v, %v", d, err)
	}
}

func TestParseXeroDate(t *testing.T) {
	for _, s := range []string{"/Date(1772496000000+0000)/", "2026-03-03T00:00:00", "2026-03-03"} {
		d, err := parseXeroDate(s)
		if err != nil || d.Format("2006-01-02") != "2026-03-03" {
			t.Errorf("%s -> %v %v", s, d, err)
		}
	}
}
