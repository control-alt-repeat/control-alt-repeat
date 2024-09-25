package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Lambda handler function
func ReadS3Object(ctx context.Context, bucket string, key string, region string) (string, error) {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region), // Replace with your region
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Call S3 to get the object
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
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
