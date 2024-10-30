package aws

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

func IterateS3Objects(ctx context.Context, bucket string, region string, f func(context.Context, string) error) error {
	s3Client, err := minio.New("s3.amazonaws.com", &minio.Options{
		Creds: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvAWS{},             // Check environment variables
			&credentials.FileAWSCredentials{}, // Check ~/.aws/credentials file
			&credentials.IAM{Client: nil},     // Check IAM roles (if running on AWS)
		}),
		Region: region,
		Secure: true,
	})
	if err != nil {
		return err
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}

	for object := range s3Client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			return err
		}

		err := f(ctx, object.Key)
		if err != nil {
			return err
		}
	}

	return nil
}
