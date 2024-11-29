package dynamodb

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamodbClient struct {
	client *dynamodb.Client
	once   sync.Once
}

var instance *dynamodbClient

func getClient(ctx context.Context) (*dynamodbClient, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &dynamodbClient{}
	instance.once.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
		if err != nil {
			return
		}
		instance.client = dynamodb.NewFromConfig(cfg)
	})
	return instance, nil
}

func put(ctx context.Context, input *dynamodb.PutItemInput) error {
	instance, err := getClient(ctx)
	if err != nil {
		return err
	}

	_, err = instance.client.PutItem(ctx, input)

	return err
}

func get[T any](ctx context.Context, input *dynamodb.GetItemInput) (*T, error) {
	instance, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	response, err := instance.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if response.Item == nil {
		return nil, nil
	}

	var result T
	err = attributevalue.UnmarshalMap(response.Item, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func query[T any](ctx context.Context, input *dynamodb.QueryInput) ([]T, error) {
	instance, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := instance.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []T
	for _, item := range result.Items {
		var i T
		err = attributevalue.UnmarshalMap(item, &i)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func scan[T any](ctx context.Context, input *dynamodb.ScanInput) ([]T, error) {
	instance, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	result, err := instance.client.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []T
	for _, item := range result.Items {
		var i T
		err = attributevalue.UnmarshalMap(item, &i)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func update(ctx context.Context, input *dynamodb.UpdateItemInput) error {
	instance, err := getClient(ctx)
	if err != nil {
		return err
	}

	_, err = instance.client.UpdateItem(ctx, input)

	return err
}
