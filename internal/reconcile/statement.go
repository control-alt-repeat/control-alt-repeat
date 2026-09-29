package reconcile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// StatementLine is one line from a bank statement, i.e. one thing Xero wants reconciled.
// Amount is negative for money out.
type StatementLine struct {
	ID          string
	Date        time.Time
	Amount      decimal.Decimal
	Payee       string
	Description string
	Reference   string
}

func (l StatementLine) IsSpend() bool {
	return l.Amount.IsNegative()
}

// Text is everything descriptive about the line, used for rule matching.
func (l StatementLine) Text() string {
	return strings.Join([]string{l.Payee, l.Description, l.Reference}, " ")
}

// StatementProfile maps a bank's CSV export onto StatementLine. Either AmountColumn
// (signed) or DebitColumn/CreditColumn (both positive) must be set.
type StatementProfile struct {
	Name              string   `json:"name"`
	SkipRows          int      `json:"skip_rows,omitempty"`
	DateColumn        string   `json:"date_column"`
	DateFormats       []string `json:"date_formats,omitempty"`
	AmountColumn      string   `json:"amount_column,omitempty"`
	NegateAmount      bool     `json:"negate_amount,omitempty"`
	DebitColumn       string   `json:"debit_column,omitempty"`
	CreditColumn      string   `json:"credit_column,omitempty"`
	PayeeColumn       string   `json:"payee_column,omitempty"`
	DescriptionColumn string   `json:"description_column,omitempty"`
	ReferenceColumn   string   `json:"reference_column,omitempty"`
}

func LoadStatement(path string, p StatementProfile) ([]StatementLine, error) {
	if p.AmountColumn == "" && p.DebitColumn == "" && p.CreditColumn == "" {
		return nil, fmt.Errorf("profile %q: set amount_column or debit_column/credit_column", p.Name)
	}

	t, err := readTable(path, p.SkipRows)
	if err != nil {
		return nil, err
	}
	if err := t.require(p.DateColumn, p.AmountColumn, p.DebitColumn, p.CreditColumn, p.PayeeColumn, p.DescriptionColumn, p.ReferenceColumn); err != nil {
		return nil, err
	}

	seen := map[string]int{}
	lines := make([]StatementLine, 0, len(t.rows))

	for i, row := range t.rows {
		rowNum := i + 2 + p.SkipRows

		date, err := parseDate(t.value(row, p.DateColumn), p.DateFormats)
		if err != nil {
			return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
		}

		var amount decimal.Decimal
		if p.AmountColumn != "" {
			if amount, err = parseMoney(t.value(row, p.AmountColumn)); err != nil {
				return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
			}
			if p.NegateAmount {
				amount = amount.Neg()
			}
		} else {
			debit, err := parseMoney(t.value(row, p.DebitColumn))
			if err != nil {
				return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
			}
			credit, err := parseMoney(t.value(row, p.CreditColumn))
			if err != nil {
				return nil, fmt.Errorf("%s row %d: %w", path, rowNum, err)
			}
			amount = credit.Abs().Sub(debit.Abs())
		}

		line := StatementLine{
			Date:        date,
			Amount:      amount,
			Payee:       t.value(row, p.PayeeColumn),
			Description: t.value(row, p.DescriptionColumn),
			Reference:   t.value(row, p.ReferenceColumn),
		}
		line.ID = lineID(line, seen)
		lines = append(lines, line)
	}

	return lines, nil
}

// lineID is a stable fingerprint of a statement line, so re-running over an
// overlapping export produces the same IDs and nothing is posted twice. Identical
// lines on the same day (two £3.50 coffees) are told apart by occurrence order.
func lineID(l StatementLine, seen map[string]int) string {
	key := strings.Join([]string{
		l.Date.Format("2006-01-02"),
		l.Amount.StringFixed(2),
		strings.ToLower(l.Payee),
		strings.ToLower(l.Description),
		strings.ToLower(l.Reference),
	}, "|")
	seen[key]++
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d", key, seen[key])))
	return hex.EncodeToString(sum[:])[:12]
}
