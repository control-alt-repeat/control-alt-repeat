package dynamodb

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func OverwriteItem(ctx context.Context, item Item) error {
	if item.Shelf == "" {
		item.Shelf = UnsetShelfDefault
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	return put(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Item:      av,
	})
}

func GetItem(ctx context.Context, id string) (*Item, error) {
	return get[Item](ctx, &dynamodb.GetItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
}

func GetUnshelvedItems(ctx context.Context) ([]Item, error) {
	return query[Item](ctx, &dynamodb.QueryInput{
		TableName:              aws.String("control-alt-repeat-warehouse"),
		KeyConditionExpression: aws.String("Shelf = :unshelved"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":unshelved": &types.AttributeValueMemberS{Value: UnsetShelfDefault},
		},
	})
}

func GetItemsUpdatedSince(ctx context.Context, since time.Time) ([]Item, error) {
	return scan[Item](ctx, &dynamodb.ScanInput{
		TableName:        aws.String("control-alt-repeat-warehouse"),
		IndexName:        aws.String("UpdatedAtIndex"),
		FilterExpression: aws.String("UpdatedAt > :since"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":since": &types.AttributeValueMemberN{Value: strconv.Itoa(int(since.Unix()))},
		},
	})
}

func UpdateOwner(ctx context.Context, itemID, newOwnerID, newOwnerName string) error {
	return update(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
		UpdateExpression: aws.String("SET FreeagentOwnerID = :newOwnerID, OwnerDisplayName = :newOwnerName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newOwnerID":   &types.AttributeValueMemberS{Value: newOwnerID},
			":newOwnerName": &types.AttributeValueMemberS{Value: newOwnerName},
		},
	})
}

func UpdateShelf(ctx context.Context, itemID, newShelf string) error {
	return update(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("control-alt-repeat-warehouse"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: itemID},
		},
		UpdateExpression: aws.String("SET Shelf = :newShelf"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newShelf": &types.AttributeValueMemberS{Value: newShelf},
		},
	})
}
