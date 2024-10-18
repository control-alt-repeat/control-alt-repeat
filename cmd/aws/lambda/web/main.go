package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	log.Printf("Echo cold start")
	app := web.Init()
	echoLambda = echoadapter.New(app)

}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
