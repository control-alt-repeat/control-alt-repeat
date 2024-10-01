package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func SaveJsonObjectS3(bucket, key string, item interface{}) error {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create an S3 client
	svc := s3.NewFromConfig(cfg)

	// Marshal the item to JSON
	jsonData, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	// Upload input parameters
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	}

	// Upload the JSON data to S3
	_, err = svc.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func LoadJsonObjectS3(bucket string, key string, object interface{}) error {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	// Create an S3 client
	svc := s3.NewFromConfig(cfg)

	// Get the object from S3
	resp, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to get object from S3: %w", err)
	}
	defer resp.Body.Close()

	// Unmarshal the JSON data into the struct
	if err := json.NewDecoder(resp.Body).Decode(object); err != nil {
		return fmt.Errorf("unable to decode JSON: %w", err)
	}

	return nil
}

func KeyExistsInS3(bucket string, key string) (bool, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		return false, err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	_, err = s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
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

// Lambda handler function
func ReadS3Object(ctx context.Context, bucket string, key string, region string) (string, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		return "", err
	}
	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	// Call S3 to get the object
	if err != nil {
		return "", fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer result.Body.Close()

	// Read the object's content
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read object content: %v", err)
	}

	// Convert the content to a string
	content := string(body)
	fmt.Println("Content of the S3 object:", content)

	return content, nil
}
