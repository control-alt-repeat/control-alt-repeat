package dynamodb

import (
	"context"
	"fmt"
	"strings"
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

type UpdateItemOptions struct {
	ItemID               string
	UpdateItemAttributes []UpdateItemAttributes
}

type UpdateItemAttributes struct {
	Name  string
	Value string
}

func UpdateItem(ctx context.Context, opts UpdateItemOptions) error {
	instance, err := getClient(ctx)
	if err != nil {
		return err
	}

	var updateExpressions []string
	expressionAttributeNames := make(map[string]string)
	expressionAttributeValues := make(map[string]types.AttributeValue)

	for _, opt := range opts.UpdateItemAttributes {
		placeholderName := fmt.Sprintf("#%s", opt.Name)
		expressionAttributeNames[placeholderName] = opt.Name

		placeholderValue := fmt.Sprintf(":val_%s", opt.Name)
		expressionAttributeValues[placeholderValue] = &types.AttributeValueMemberS{Value: opt.Value}

		updateExpressions = append(updateExpressions, fmt.Sprintf("%s = %s", placeholderName, placeholderValue))
	}

	updateExpression := strings.Join(updateExpressions, ", ")

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: opts.ItemID},
		},
		UpdateExpression:          aws.String("SET " + updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err = instance.client.UpdateItem(ctx, input)

	return err
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

type ScanItemsOptions struct {
	IndexName              string
	FilterExpression       string
	NumberExpressionValues map[string]string
}

func ScanItems(ctx context.Context, opts ScanItemsOptions) ([]Item, error) {
	expressionAttributeValues := map[string]types.AttributeValue{}
	for k, v := range opts.NumberExpressionValues {
		expressionAttributeValues[k] = &types.AttributeValueMemberN{Value: v}
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String("control-alt-repeat-warehouse"),
		FilterExpression:          aws.String(opts.FilterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	if opts.IndexName != "" {
		input.IndexName = aws.String(opts.IndexName)
	}

	instance, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := instance.client.Scan(ctx, input)
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
