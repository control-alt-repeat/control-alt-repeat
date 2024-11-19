package accounting

import (
	"context"
	"fmt"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/reports"
	"github.com/control-alt-repeat/control-alt-repeat/internal/freeagent"
)

const (
	EmptyCellToken = "--"

	EbayPostageLabel = "Postage label"
	EbayOrder        = "Order"
	EbayOtherFee     = "Other fee"
	EbayPayout       = "Payout"
	EbayRefund       = "Refund"
	EbayHold         = "Hold"
	EbayCharge       = "Charge"

	FreeAgentSales          = "https://api.freeagent.com/v2/categories/001"
	FreeAgentCostOfSales    = "https://api.freeagent.com/v2/categories/101"
	FreeAgentCommissionPaid = "https://api.freeagent.com/v2/categories/102"
)

func MapEbayTransactionsToFreeAgent(ctx context.Context, t reports.Transaction) ([]freeagent.BankTransactionExplanation, error) {
	if t.PayoutCurrency != "GBP" {
		return nil, fmt.Errorf("only GBP currency supported - %s no catered for", t.PayoutCurrency)
	}
	if t.Type == EbayPayout {
		return nil, nil // Payouts are "linked transactions" in FreeAgent TODO
	}

	var explanation freeagent.BankTransactionExplanation

	explanation.GrossValue = t.GrossTransactionAmount.StringFixedBank(2)
	explanation.DatedOn = freeagent.FreeAgentDate(t.TransactionCreationDate)
	explanation.IsDeletable = true
	explanation.Description = generateDescription(t)
	category, err := mapEbayTransactionType(t.Type)
	if err != nil {
		return nil, err
	}
	explanation.Category = category

	explanations := []freeagent.BankTransactionExplanation{}

	if t.Type == EbayOrder || t.Type == EbayRefund {
		finalValueFixed := explanation
		finalValueFixed.Category = FreeAgentCommissionPaid
		finalValueFixed.GrossValue = t.FinalValueFeeFixed.Decimal.StringFixedBank(2)
		finalValueFixed.Description = fmt.Sprintf("%s | Fixed Fee", explanation.Description)

		finalValueVariable := explanation
		finalValueVariable.Category = FreeAgentCommissionPaid
		finalValueVariable.GrossValue = t.FinalValueFeeVariable.Decimal.StringFixedBank(2)
		finalValueVariable.Description = fmt.Sprintf("%s | Variable Fee", explanation.Description)

		regulatoryOperatingFee := explanation
		regulatoryOperatingFee.Category = FreeAgentCommissionPaid
		regulatoryOperatingFee.GrossValue = t.RegulatoryOperatingFee.Decimal.StringFixedBank(2)
		regulatoryOperatingFee.Description = fmt.Sprintf("%s | Regulatory Operating Fee", explanation.Description)

		explanations = append(explanations, finalValueFixed)
		explanations = append(explanations, finalValueVariable)
		explanations = append(explanations, regulatoryOperatingFee)

		explanation.Description = explanation.Description + " | " + t.Type
	}

	explanations = append(explanations, explanation)

	return explanations, nil
}

func mapEbayTransactionType(eBayTransactionType string) (string, error) {
	switch strings.TrimSpace(eBayTransactionType) {
	case EbayPostageLabel, EbayCharge, EbayOtherFee:
		return FreeAgentCostOfSales, nil
	case EbayOrder, EbayHold, EbayRefund:
		return FreeAgentSales, nil
	}

	return "", fmt.Errorf("eBay transaction type '%s' not recognised", eBayTransactionType)
}

func generateDescription(t reports.Transaction) string {
	tokens := []string{}

	tokens = appendIfHasValue(tokens, t.OrderNumber)
	tokens = appendIfHasValue(tokens, t.ItemID)
	tokens = appendIfHasValue(tokens, t.CustomLabel)
	tokens = appendIfHasValue(tokens, t.ReferenceID)
	tokens = appendIfHasValue(tokens, t.Description)

	return strings.TrimSpace(strings.Join(tokens, " | "))
}

func appendIfHasValue(tokens []string, token string) []string {
	if strings.TrimSpace(token) == EmptyCellToken {
		return tokens
	}

	return append(tokens, token)
}
