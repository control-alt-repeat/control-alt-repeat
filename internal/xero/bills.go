package xero

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Invoice struct {
	InvoiceID     string   `json:",omitempty"`
	InvoiceNumber string   `json:",omitempty"`
	Type          string   `json:",omitempty"` // ACCPAY = bill
	Status        string   `json:",omitempty"`
	Contact       *Contact `json:",omitempty"`
	Reference     string   `json:",omitempty"`
	DateString    string   `json:",omitempty"`
	DueDateString string   `json:",omitempty"`
	Total         *Money   `json:",omitempty"`
	AmountDue     *Money   `json:",omitempty"`
	CurrencyCode  string   `json:",omitempty"`
}

func (i Invoice) DateValue() (time.Time, error) {
	return parseXeroDate(i.DateString)
}

type Payment struct {
	PaymentID   string   `json:",omitempty"`
	Invoice     *Invoice `json:",omitempty"`
	Account     *Account `json:",omitempty"`
	Date        string   `json:",omitempty"` // yyyy-mm-dd when writing
	Amount      *Money   `json:",omitempty"`
	Reference   string   `json:",omitempty"`
	PaymentType string   `json:",omitempty"`
	Status      string   `json:",omitempty"`

	StatusAttributeString string            `json:",omitempty"`
	ValidationErrors      []ValidationError `json:",omitempty"`
}

// parseXeroDate handles both "2026-03-03T00:00:00" and the older "/Date(1772496000000+0000)/".
func parseXeroDate(s string) (time.Time, error) {
	if strings.HasPrefix(s, "/Date(") {
		ms := strings.TrimPrefix(s, "/Date(")
		if i := strings.IndexAny(ms, "+-)"); i > 0 {
			ms = ms[:i]
		}
		n, err := strconv.ParseInt(ms, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot parse xero date %q", s)
		}
		t := time.UnixMilli(n).UTC()
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	if len(s) >= 10 {
		s = s[:10]
	}
	return time.Parse("2006-01-02", s)
}

// ListUnpaidBills returns every approved supplier bill with money still owing.
func (c *Client) ListUnpaidBills(ctx context.Context) ([]Invoice, error) {
	var all []Invoice
	for page := 1; ; page++ {
		q := url.Values{}
		q.Set("where", `Type=="ACCPAY" AND Status=="AUTHORISED"`)
		q.Set("summaryOnly", "true")
		q.Set("page", strconv.Itoa(page))
		q.Set("pageSize", "1000")

		var resp struct{ Invoices []Invoice }
		if err := c.get(ctx, "Invoices", q, &resp); err != nil {
			return nil, err
		}
		if len(resp.Invoices) == 0 {
			return all, nil
		}
		for _, inv := range resp.Invoices {
			if inv.AmountDue != nil && inv.AmountDue.IsPositive() {
				all = append(all, inv)
			}
		}
	}
}

// ListPayments returns bill and invoice payments recorded against a bank account.
func (c *Client) ListPayments(ctx context.Context, accountID string, from, to time.Time) ([]Payment, error) {
	where := fmt.Sprintf(`Account.AccountID==guid("%s") AND Date>=DateTime(%d,%d,%d) AND Date<=DateTime(%d,%d,%d) AND Status!="DELETED"`,
		accountID, from.Year(), from.Month(), from.Day(), to.Year(), to.Month(), to.Day())

	var all []Payment
	for page := 1; ; page++ {
		q := url.Values{}
		q.Set("where", where)
		q.Set("page", strconv.Itoa(page))
		q.Set("pageSize", "1000")

		var resp struct {
			Payments []struct {
				Payment
				DateString string `json:",omitempty"`
			}
		}
		if err := c.get(ctx, "Payments", q, &resp); err != nil {
			return nil, err
		}
		if len(resp.Payments) == 0 {
			return all, nil
		}
		for _, p := range resp.Payments {
			if p.DateString != "" {
				p.Date = p.DateString
			}
			all = append(all, p.Payment)
		}
	}
}

// ListAccountActivity returns everything already recorded against a bank account -
// spend/receive money and bill/invoice payments - as bank transactions, so callers
// can check a statement line is not already accounted for either way.
func (c *Client) ListAccountActivity(ctx context.Context, accountID string, from, to time.Time) ([]BankTransaction, error) {
	txns, err := c.ListBankTransactions(ctx, accountID, from, to)
	if err != nil {
		return nil, err
	}
	payments, err := c.ListPayments(ctx, accountID, from, to)
	if err != nil {
		return nil, err
	}
	for _, p := range payments {
		t := BankTransaction{BankTransactionID: p.PaymentID, Date: p.Date, Total: p.Amount, Type: BankTransactionReceive, Reference: p.Reference}
		// Money leaves the bank when paying a bill or refunding a customer.
		if p.PaymentType == "ACCPAYPAYMENT" || p.PaymentType == "ARCREDITPAYMENT" || p.PaymentType == "AROVERPAYMENTPAYMENT" || p.PaymentType == "ARPREPAYMENTPAYMENT" {
			t.Type = BankTransactionSpend
		}
		txns = append(txns, t)
	}
	return txns, nil
}

// CreatePayments records payments against bills. Like bank transactions they are
// left unreconciled so Xero offers them as the match for the statement line.
func (c *Client) CreatePayments(ctx context.Context, payments []Payment) ([]Payment, error) {
	q := url.Values{}
	q.Set("summarizeErrors", "false")

	var resp struct{ Payments []Payment }
	if err := c.put(ctx, "Payments", q, map[string]any{"Payments": payments}, &resp); err != nil {
		return nil, err
	}
	return resp.Payments, nil
}

// AttachToBankTransaction uploads a file to a spend/receive money transaction.
func (c *Client) AttachToBankTransaction(ctx context.Context, bankTransactionID, fileName, contentType string, content []byte) error {
	path := fmt.Sprintf("BankTransactions/%s/Attachments/%s", bankTransactionID, url.PathEscape(fileName))
	return c.doRaw(ctx, http.MethodPut, path, contentType, content)
}

func (c *Client) doRaw(ctx context.Context, method, path, contentType string, body []byte) error {
	if err := c.throttle(ctx); err != nil {
		return err
	}
	token, err := c.tokens.AccessToken(ctx)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, apiBaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", contentType)
	if c.tenantID != "" {
		req.Header.Set("xero-tenant-id", c.tenantID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("xero %s %s failed: %d %s", method, path, resp.StatusCode, b)
	}
	return nil
}
