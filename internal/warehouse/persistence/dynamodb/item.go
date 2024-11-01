package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

type Item struct {
	ID               string `dynamodbav:"ID"`
	Shelf            string `dynamodbav:"Shelf"`
	EbayListingID    string `dynamodbav:"EbayListingID"`
	FreeagentOwnerID string `dynamodbav:"FreeagentOwnerID"`
	OwnerDisplayName string `dynamodbav:"OwnerDisplayName"`
	CreatedAt        int64  `dynamodbav:"CreatedAt"`
	UpdatedAt        int64  `dynamodbav:"UpdatedAt"`
}

type SaveItemOptions struct {
	Item Item
}

func SaveItem(ctx context.Context, opts SaveItemOptions) error {
	av, err := attributevalue.MarshalMap(opts.Item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Item:      av,
	}

	instance, err := getClient(ctx)
	if err != nil {
		return err
	}

	_, err = instance.client.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
