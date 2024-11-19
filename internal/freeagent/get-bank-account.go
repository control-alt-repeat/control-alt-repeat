package freeagent

import (
	"context"
	"fmt"
)

func GetBankAccount(ctx context.Context, opts GetBankAccountOptions) (BankAccount, error) {
	var response GetBankAccountResponse
	apiopts := ApiGetOptions{
		Path: fmt.Sprintf("bank_accounts/%s", opts.BankAccountID),
	}
	if err := FreeagentApiGet(ctx, apiopts, &response); err != nil {
		return BankAccount{}, err
	}

	bankAccountID, err := apiopts.buildUrl()
	if err != nil {
		return BankAccount{}, err
	}

	response.BankAccount.ID = bankAccountID

	return response.BankAccount, nil
}

type GetBankAccountOptions struct {
	BankAccountID string
}

type GetBankAccountResponse struct {
	BankAccount BankAccount `json:"bank_account,omitempty"`
}
