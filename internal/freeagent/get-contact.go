package freeagent

import (
	"context"
	"fmt"
	"strings"
	"time"
)

func GetContact(ctx context.Context, opts GetContactOptions) (Contact, error) {
	var response GetContactResponse
	apiopts := ApiGetOptions{
		Path: fmt.Sprintf("contacts/%s", opts.ContactID),
	}
	if err := FreeagentApiGet(ctx, apiopts, &response); err != nil {
		return Contact{}, err
	}

	return response.Contact, nil
}

type GetContactOptions struct {
	ContactID string
}

type GetContactResponse struct {
	Contact Contact `json:"contact,omitempty"`
}

type Contact struct {
	URL                        string    `json:"url,omitempty"`
	FirstName                  string    `json:"first_name,omitempty"`
	LastName                   string    `json:"last_name,omitempty"`
	OrganisationName           string    `json:"organisation_name,omitempty"`
	Email                      string    `json:"email,omitempty"`
	BillingEmail               string    `json:"billing_email,omitempty"`
	PhoneNumber                string    `json:"phone_number,omitempty"`
	Mobile                     string    `json:"mobile,omitempty"`
	Address1                   string    `json:"address1,omitempty"`
	Address2                   string    `json:"address2,omitempty"`
	Address3                   string    `json:"address3,omitempty"`
	Town                       string    `json:"town,omitempty"`
	Region                     string    `json:"region,omitempty"`
	Postcode                   string    `json:"postcode,omitempty"`
	Country                    string    `json:"country,omitempty"`
	ContactNameOnInvoices      bool      `json:"contact_name_on_invoices,omitempty"`
	DefaultPaymentTermsInDays  int       `json:"default_payment_terms_in_days,omitempty"`
	Locale                     string    `json:"locale,omitempty"`
	AccountBalance             string    `json:"account_balance,omitempty"`
	UsesContactInvoiceSequence bool      `json:"uses_contact_invoice_sequence,omitempty"`
	ChargeSalesTax             string    `json:"charge_sales_tax,omitempty"`
	SalesTaxRegistrationNumber string    `json:"sales_tax_registration_number,omitempty"`
	ActiveProjectsCount        int       `json:"active_projects_count,omitempty"`
	DirectDebitMandateState    string    `json:"direct_debit_mandate_state,omitempty"`
	Status                     string    `json:"status,omitempty"`
	CreatedAt                  time.Time `json:"created_at,omitempty"`
	UpdatedAt                  time.Time `json:"updated_at,omitempty"`
}

func (c Contact) ID() string {
	parts := strings.Split(c.URL, "/")
	lastSegment := parts[len(parts)-1]
	return lastSegment
}

func (c Contact) DisplayName() string {
	var displayName = ""
	if c.FirstName != "" || c.LastName != "" {
		displayName = strings.TrimSpace(strings.Join([]string{c.FirstName, c.LastName}, " "))
	}

	if c.OrganisationName == "" {
		return displayName
	}

	if displayName == "" {
		return c.OrganisationName
	}

	return strings.Join([]string{c.OrganisationName, displayName}, " - ")
}
