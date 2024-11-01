package aws

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func SaveBytesToS3(ctx context.Context, bucket, key string, data []byte, contentType string) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return err
	}

	svc := s3.NewFromConfig(cfg)

	_, err = svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}
	return nil
}
