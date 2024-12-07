package finances

import (
	"context"
	"fmt"
	"time"
)

func GetTransaction(ctx context.Context, orderID string) (GetTransactionResponse, error) {
	apiopts := requestOptions{
		Path: "/transaction",
		QueryParameters: map[string]string{
			"filter": fmt.Sprintf("orderId:{%s}", orderID),
		},
	}

	var response GetTransactionResponse
	if err := apiGet(ctx, apiopts, &response); err != nil {
		return GetTransactionResponse{}, err
	}

	return response, nil
}

type GetTransactionResponse struct {
	Transactions []struct {
		TransactionID   string `json:"transactionId"`
		OrderID         string `json:"orderId"`
		TransactionType string `json:"transactionType"`
		TransactionMemo string `json:"transactionMemo"`
		Amount          struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"amount"`
		NetAmount struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		} `json:"netAmount"`
		Description string    `json:"description"`
		OccurredAt  time.Time `json:"occurredAt"`
	} `json:"transactions"`
}
