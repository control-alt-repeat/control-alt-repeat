package freeagent

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Attachment struct {
	URL         string `json:"url,omitempty"`
	ContentSrc  string `json:"content_src,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	FileSize    int    `json:"file_size,omitempty"`
}

type BankAccount struct {
	ID               string
	OpeningBalance   decimal.Decimal `json:"opening_balance"`
	Type             string          `json:"type"`
	Name             string          `json:"name"`
	IsPersonal       bool            `json:"is_personal"`
	Status           string          `json:"status"`
	Currency         string          `json:"currency"`
	CurrentBalance   decimal.Decimal `json:"current_balance"`
	UpdatedAt        time.Time       `json:"updated_at"`
	CreatedAt        time.Time       `json:"created_at"`
	BankGuessEnabled bool            `json:"bank_guess_enabled"`
}

type FreeAgentDate time.Time

type BankTransactionExplanation struct {
	URL             string        `json:"url,omitempty"`
	BankTransaction string        `json:"bank_transaction,omitempty"`
	BankAccount     string        `json:"bank_account,omitempty"`
	Category        string        `json:"category,omitempty"`
	DatedOn         FreeAgentDate `json:"dated_on,omitempty"`
	Description     string        `json:"description,omitempty"`
	GrossValue      string        `json:"gross_value,omitempty"`
	Project         string        `json:"project,omitempty"`
	RebillType      string        `json:"rebill_type,omitempty"`
	RebillFactor    string        `json:"rebill_factor,omitempty"`
	SalesTaxStatus  string        `json:"sales_tax_status,omitempty"`
	SalesTaxRate    string        `json:"sales_tax_rate,omitempty"`
	SalesTaxValue   string        `json:"sales_tax_value,omitempty"`
	IsDeletable     bool          `json:"is_deletable,omitempty"`
}

func (cd FreeAgentDate) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, time.Time(cd).Format("2006-01-02"))
	return []byte(formatted), nil
}

func (cd *FreeAgentDate) UnmarshalJSON(data []byte) error {
	parsed, err := time.Parse(`"2006-01-02"`, string(data))
	if err != nil {
		return err
	}
	*cd = FreeAgentDate(parsed)
	return nil
}

func (cd FreeAgentDate) ToTime() time.Time {
	return time.Time(cd)
}
