package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context, s3Event events.S3Event) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	// Create a new S3 service client
	svc := s3.New(sess)

	for _, record := range s3Event.Records {
		bucket := record.S3.Bucket.Name
		key, err := url.QueryUnescape(record.S3.Object.Key)
		if err != nil {
			return fmt.Errorf("failed to unescape key %s: %v", key, err)
		}

		fmt.Println("Processing event: ", key)

		// Get the object from S3
		result, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return fmt.Errorf("failed to get object %s from bucket %s: %v", key, bucket, err)
		}
		defer result.Body.Close() // Ensure the body is closed after reading

		// Read the object body
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return fmt.Errorf("failed to read object body: %v", err)
		}

		var notification models.ItemNotificationEnvelope
		err = xml.Unmarshal(body, &notification)
		if err != nil {
			return fmt.Errorf("failed to Unmarshal body: %v", err)
		}

		return internal.ImportEbayListing(notification.Body.GetItemResponse.Item.ItemID)
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
