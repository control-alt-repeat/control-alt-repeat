package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := ebay.HandleNotification(req.Body); err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK, // this is our eBay notification endpoint. If they see errors one this they may cut us off.
		}, nil
	}

	fmt.Println("Responding that everything is OK.")

	// Return the message or default response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
