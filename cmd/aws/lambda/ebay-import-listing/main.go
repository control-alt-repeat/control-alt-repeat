package main

import (
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay-import-listing/ebay"
	"github.com/aws/aws-lambda-go/lambda"
)

func hello() error {
	err := ebay.ImportListing()

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}
