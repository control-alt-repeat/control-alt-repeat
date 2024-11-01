package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
)

type Item struct {
	ID               string `dynamodbav:"ID"`
	Shelf            string `dynamodbav:"Shelf"`
	Title            string `dynamodbav:"Title"`
	PictureURL       string `dynamodbav:"PictureURL"`
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

type QueryItemsOptions struct {
	IndexName              string
	KeyConditionExpression string
	StringExpressionValues map[string]string
}

func QueryItems(ctx context.Context, opts QueryItemsOptions) ([]Item, error) {
	expressionAttributeValues := map[string]types.AttributeValue{}
	for k, v := range opts.StringExpressionValues {
		expressionAttributeValues[k] = &types.AttributeValueMemberS{Value: v}
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String("control-alt-repeat-warehouse"),
		KeyConditionExpression:    aws.String(opts.KeyConditionExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	if opts.IndexName != "" {
		input.IndexName = aws.String(opts.IndexName)
	}

	instance, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := instance.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []Item
	for _, item := range result.Items {
		var i Item
		err = attributevalue.UnmarshalMap(item, &i)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (i Item) Map() models.WarehouseItem {
	return models.WarehouseItem{
		ControlAltRepeatID: i.ID,
		Title:              i.Title,
		Shelf:              i.Shelf,
		AddedTime:          time.Unix(i.CreatedAt, 0),
		EbayListingID:      i.EbayListingID,
		FreeagentOwnerID:   i.FreeagentOwnerID,
		OwnerDisplayName:   i.OwnerDisplayName,
	}
}
