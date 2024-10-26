package main

import (
	"context"
	"log"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/labstack/echo/v4"
)

var (
	echoLambda *echoadapter.EchoLambdaV2
)

func init() {
	e := echo.New()
	err := web.Init(e)
	if err != nil {
		log.Fatalf("Failed to initialize Echo app: %v", err)
	}
	echoLambda = echoadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
