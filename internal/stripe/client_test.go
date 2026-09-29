package stripe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPaginationAndParsing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer rk_test" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Query().Get("starting_after") {
		case "":
			_, _ = w.Write([]byte(`{"has_more":true,"data":[{"id":"txn_1","type":"charge","amount":10000,"fee":170,"net":9830}]}`))
		case "txn_1":
			_, _ = w.Write([]byte(`{"has_more":false,"data":[{"id":"txn_2","type":"refund","amount":-2000,"fee":0,"net":-2000}]}`))
		}
	}))
	defer srv.Close()

	c := &Client{key: "rk_test", baseURL: srv.URL + "/"}
	txns, err := c.PayoutTransactions(context.Background(), "po_1")
	if err != nil {
		t.Fatal(err)
	}
	if len(txns) != 2 || txns[0].Amount.String() != "100" || txns[0].Fee.String() != "1.7" || txns[1].Amount.String() != "-20" {
		t.Errorf("got %+v", txns)
	}
}
