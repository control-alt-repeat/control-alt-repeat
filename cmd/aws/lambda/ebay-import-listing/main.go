package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {
	ebayToken, err := aws.ReadS3Object(
		ctx,
		os.Getenv("EBAY_TOKEN_S3_BUCKET"),
		os.Getenv("EBAY_TOKEN_S3_KEY"),
		"eu-west-2",
	)

	if err != nil {
		fmt.Println(err)

		return err
	}

	err = ebay.ImportListing(ebayToken)

	if err != nil {
		fmt.Println(err)

		return err
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
