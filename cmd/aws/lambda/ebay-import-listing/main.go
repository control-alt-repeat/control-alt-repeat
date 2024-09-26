package main

import (
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler() {
	err := ebay.ImportListing()

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
