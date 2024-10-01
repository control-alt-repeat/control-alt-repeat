package main

import (
	"context"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) {
	panic("not implemented!")
	err := internal.ImportEbayListing("")

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
