package xero

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Money marshals as a plain JSON number with 2dp, which is what Xero expects.
type Money struct {
	decimal.Decimal
}

func NewMoney(d decimal.Decimal) *Money {
	return &Money{d}
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(m.StringFixedBank(2)), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	return m.Decimal.UnmarshalJSON(bytes.Trim(data, `"`))
}

type Account struct {
	AccountID         string `json:",omitempty"`
	Code              string `json:",omitempty"`
	Name              string `json:",omitempty"`
	Type              string `json:",omitempty"`
	TaxType           string `json:",omitempty"`
	Status            string `json:",omitempty"`
	BankAccountNumber string `json:",omitempty"`
	CurrencyCode      string `json:",omitempty"`
}

type TaxRate struct {
	Name          string
	TaxType       string
	Status        string
	EffectiveRate decimal.Decimal
}

type Contact struct {
	ContactID string `json:",omitempty"`
	Name      string `json:",omitempty"`
}

type LineItem struct {
	Description string  `json:",omitempty"`
	Quantity    float64 `json:",omitempty"`
	UnitAmount  *Money  `json:",omitempty"`
	AccountCode string  `json:",omitempty"`
	TaxType     string  `json:",omitempty"`
	LineAmount  *Money  `json:",omitempty"`
}

type ValidationError struct {
	Message string
}

type BankTransaction struct {
	BankTransactionID string     `json:",omitempty"`
	Type              string     `json:",omitempty"` // SPEND or RECEIVE
	Contact           *Contact   `json:",omitempty"`
	LineItems         []LineItem `json:",omitempty"`
	BankAccount       *Account   `json:",omitempty"`
	IsReconciled      bool
	Date              string `json:",omitempty"` // yyyy-mm-dd when writing
	DateString        string `json:",omitempty"` // returned by Xero when reading
	Reference         string `json:",omitempty"`
	Status            string `json:",omitempty"`
	LineAmountTypes   string `json:",omitempty"` // Inclusive, Exclusive or NoTax
	Total             *Money `json:",omitempty"`

	StatusAttributeString string            `json:",omitempty"`
	ValidationErrors      []ValidationError `json:",omitempty"`
}

const (
	BankTransactionSpend   = "SPEND"
	BankTransactionReceive = "RECEIVE"
)

// DateValue returns the transaction date regardless of whether it was read or built locally.
func (b BankTransaction) DateValue() (time.Time, error) {
	if b.DateString != "" {
		return parseXeroDate(b.DateString)
	}
	return parseXeroDate(b.Date)
}

// SignedTotal is negative for money out, matching how bank statements present it.
func (b BankTransaction) SignedTotal() decimal.Decimal {
	if b.Total == nil {
		return decimal.Zero
	}
	if strings.HasPrefix(b.Type, BankTransactionSpend) {
		return b.Total.Neg()
	}
	return b.Total.Decimal
}

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var resp struct{ Accounts []Account }
	return resp.Accounts, c.get(ctx, "Accounts", nil, &resp)
}

func (c *Client) ListTaxRates(ctx context.Context) ([]TaxRate, error) {
	var resp struct{ TaxRates []TaxRate }
	return resp.TaxRates, c.get(ctx, "TaxRates", nil, &resp)
}

// FindBankAccount accepts an account code, AccountID, bank account number or exact name.
func (c *Client) FindBankAccount(ctx context.Context, ref string) (Account, error) {
	accounts, err := c.ListAccounts(ctx)
	if err != nil {
		return Account{}, err
	}
	for _, a := range accounts {
		if a.Type != "BANK" {
			continue
		}
		if a.AccountID == ref || (a.Code != "" && a.Code == ref) || a.BankAccountNumber == ref || strings.EqualFold(a.Name, ref) {
			return a, nil
		}
	}
	return Account{}, fmt.Errorf("no xero bank account matching %q (see `car xero accounts`)", ref)
}

// ListBankTransactions returns every spend/receive money transaction on a bank
// account between two dates (inclusive), following pagination.
func (c *Client) ListBankTransactions(ctx context.Context, accountID string, from, to time.Time) ([]BankTransaction, error) {
	where := fmt.Sprintf(`BankAccount.AccountID==guid("%s") AND Date>=DateTime(%d,%d,%d) AND Date<=DateTime(%d,%d,%d) AND Status!="DELETED"`,
		accountID, from.Year(), from.Month(), from.Day(), to.Year(), to.Month(), to.Day())

	var all []BankTransaction
	for page := 1; ; page++ {
		q := url.Values{}
		q.Set("where", where)
		q.Set("page", strconv.Itoa(page))
		q.Set("pageSize", "1000")

		var resp struct{ BankTransactions []BankTransaction }
		if err := c.get(ctx, "BankTransactions", q, &resp); err != nil {
			return nil, err
		}
		if len(resp.BankTransactions) == 0 {
			return all, nil
		}
		all = append(all, resp.BankTransactions...)
	}
}

// CreateBankTransactions creates transactions one batch at a time. Each returned
// transaction carries its own StatusAttributeString ("OK" / "ERROR") so a single
// bad row never blocks the rest.
func (c *Client) CreateBankTransactions(ctx context.Context, txns []BankTransaction) ([]BankTransaction, error) {
	q := url.Values{}
	q.Set("summarizeErrors", "false")

	var resp struct{ BankTransactions []BankTransaction }
	if err := c.put(ctx, "BankTransactions", q, map[string]any{"BankTransactions": txns}, &resp); err != nil {
		return nil, err
	}
	return resp.BankTransactions, nil
}
