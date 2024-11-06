package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/control-alt-repeat/control-alt-repeat/internal/web"
)

var (
	echoLambda *echoadapter.EchoLambdaV2
)

func init() {
	e, err := web.Init(os.Stdout)
	if err != nil {
		panic(err)
	}
	echoLambda = echoadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
