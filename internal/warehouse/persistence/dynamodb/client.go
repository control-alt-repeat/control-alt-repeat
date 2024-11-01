package dynamodb

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
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
