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
	fmt.Println("Request Path: ", req.Path)
	fmt.Println("Request HTTP Method: ", req.HTTPMethod)
	fmt.Println("Request Body: ", req.Body)

	if message, ok := req.QueryStringParameters["message"]; ok {
		fmt.Println("Message from URL query: ", message)
	} else {
		fmt.Println("No message found in URL query.")
	}

	ebay.HandleNotification()

	// Return the message or default response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Message received: " + req.Body,
	}, nil
}

func main() {
	lambda.Start(handler)
}
