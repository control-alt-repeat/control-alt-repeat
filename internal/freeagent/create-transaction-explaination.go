package freeagent

import "context"

type CreateBankTransactionExplainationRequest struct {
	BankTransactionExplanation BankTransactionExplanation `json:"bank_transaction_explanation,omitempty"`
}

func CreateBankTransactionExplaination(ctx context.Context, requestObject CreateBankTransactionExplainationRequest) error {
	apiopts := requestOptions{
		Path: "bank_transaction_explanations",
	}
	return apiPost(ctx, apiopts, requestObject)
}
