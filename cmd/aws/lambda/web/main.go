package main

import (
	"context"
	"log"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	app, err := web.Init()
	if err != nil {
		log.Fatalf("Failed to initialize Echo app: %v", err)
	}

	// Initialize the Lambda adapter
	echoLambda = echoadapter.New(app)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Use echoLambda to proxy the request
	response, err := echoLambda.ProxyWithContext(ctx, req)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Ensure the Content-Type is set correctly
	if response.Headers == nil {
		response.Headers = map[string]string{}
	}
	response.Headers["Content-Type"] = "text/html; charset=utf-8" // Set to HTML

	return response, nil
}

func main() {
	lambda.Start(handler)
}
