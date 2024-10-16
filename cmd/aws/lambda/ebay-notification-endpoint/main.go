package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ebay.HandleNotification(req.Body)

	fmt.Println("Responding that everything is OK.")

	// Return the message or default response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
