package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/control-alt-repeat/control-alt-repeat/internal"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
)

var log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.With().
		Timestamp().
		Str("service", "ebay-import-listing").
		Logger().
		Level(zerolog.DebugLevel)
}

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
		body, err := io.ReadAll(result.Body)
		if err != nil {
			return fmt.Errorf("failed to read object body: %v", err)
		}

		var notification ebay.ItemNotificationEnvelope
		err = xml.Unmarshal(body, &notification)
		if err != nil {
			return fmt.Errorf("failed to Unmarshal body: %v", err)
		}

		_, err = internal.ImportEbayListing(ctx, &notification.Body.GetItemResponse.Item)
		if err != nil {
			return fmt.Errorf("failed to import listing: %v", err)
		}
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
