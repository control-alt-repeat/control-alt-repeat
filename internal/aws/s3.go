package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SaveBytesToS3(ctx context.Context, bucket, key string, data []byte, contentType string) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	_, err = svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return nil
}

func SaveJsonObjectS3(ctx context.Context, bucket, key string, item interface{}) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	}

	_, err = svc.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func LoadJsonObjectS3(ctx context.Context, bucket string, key string, object interface{}) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}
	svc := s3.NewFromConfig(cfg)
	fmt.Println(bucket)
	fmt.Println(key)
	resp, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to get object from S3: %w", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(object); err != nil {
		return fmt.Errorf("unable to decode JSON: %w", err)
	}
	return nil
}

func KeyExistsInS3(ctx context.Context, bucket string, key string) (bool, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		return false, err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ReadS3Object(ctx context.Context, bucket string, key string, region string) (string, error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		return "", err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object content: %v", err)
	}

	content := string(body)
	fmt.Println("Content of the S3 object:", content)

	return content, nil
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
		log.Fatalln(err)
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}

	for object := range s3Client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil
		}

		err := f(ctx, object.Key)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return nil
}
