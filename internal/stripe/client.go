// Package stripe reads payouts and the balance transactions inside them, so a
// Stripe payout in the bank can be booked as gross sales, refunds and fees.
package stripe

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

const apiBaseURL = "https://api.stripe.com/v1/"

type Client struct {
	key     string
	baseURL string
}

// NewClientFromEnv uses STRIPE_API_KEY - a restricted key with read access to
// Balance and Payouts is all this needs.
func NewClientFromEnv() (*Client, error) {
	key := os.Getenv("STRIPE_API_KEY")
	if key == "" {
		return nil, errors.New("STRIPE_API_KEY must be set")
	}
	return &Client{key: key, baseURL: apiBaseURL}, nil
}

type Payout struct {
	ID          string
	Amount      decimal.Decimal
	Currency    string
	ArrivalDate time.Time
	Status      string
}

type BalanceTransaction struct {
	ID          string
	Type        string
	Amount      decimal.Decimal // gross, negative for money out
	Fee         decimal.Decimal // positive
	Net         decimal.Decimal
	Description string
}

type list[T any] struct {
	Data    []T  `json:"data"`
	HasMore bool `json:"has_more"`
}

type apiPayout struct {
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	ArrivalDate int64  `json:"arrival_date"`
	Status      string `json:"status"`
}

type apiBalanceTransaction struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Amount      int64  `json:"amount"`
	Fee         int64  `json:"fee"`
	Net         int64  `json:"net"`
	Description string `json:"description"`
}

func pence(v int64) decimal.Decimal {
	return decimal.New(v, -2)
}

// ListPayouts returns paid payouts arriving between two dates (inclusive).
func (c *Client) ListPayouts(ctx context.Context, from, to time.Time) ([]Payout, error) {
	q := url.Values{}
	q.Set("arrival_date[gte]", strconv.FormatInt(from.Unix(), 10))
	q.Set("arrival_date[lte]", strconv.FormatInt(to.AddDate(0, 0, 1).Unix()-1, 10))

	var out []Payout
	err := paginate(ctx, c, "payouts", q, func(p apiPayout) string {
		if p.Status == "paid" {
			out = append(out, Payout{
				ID: p.ID, Amount: pence(p.Amount), Currency: p.Currency, Status: p.Status,
				ArrivalDate: time.Unix(p.ArrivalDate, 0).UTC().Truncate(24 * time.Hour),
			})
		}
		return p.ID
	})
	return out, err
}

// PayoutTransactions returns what makes up a payout: charges, refunds, fees, etc.
func (c *Client) PayoutTransactions(ctx context.Context, payoutID string) ([]BalanceTransaction, error) {
	q := url.Values{}
	q.Set("payout", payoutID)

	var out []BalanceTransaction
	err := paginate(ctx, c, "balance_transactions", q, func(b apiBalanceTransaction) string {
		out = append(out, BalanceTransaction{
			ID: b.ID, Type: b.Type, Amount: pence(b.Amount), Fee: pence(b.Fee), Net: pence(b.Net), Description: b.Description,
		})
		return b.ID
	})
	return out, err
}

func paginate[T any](ctx context.Context, c *Client, path string, q url.Values, each func(T) string) error {
	q.Set("limit", "100")
	for {
		var page list[T]
		if err := c.get(ctx, path, q, &page); err != nil {
			return err
		}
		last := ""
		for _, item := range page.Data {
			last = each(item)
		}
		if !page.HasMore || last == "" {
			return nil
		}
		q.Set("starting_after", last)
	}
}

func (c *Client) get(ctx context.Context, path string, q url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path+"?"+q.Encode(), http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stripe GET %s failed: %d %s", path, resp.StatusCode, body)
	}
	return json.Unmarshal(body, out)
}
