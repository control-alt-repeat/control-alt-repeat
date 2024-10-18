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
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
