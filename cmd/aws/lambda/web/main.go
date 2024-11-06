package main

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/control-alt-repeat/control-alt-repeat/internal/logger"
	"github.com/control-alt-repeat/control-alt-repeat/internal/web"

	"github.com/labstack/echo/v4"
)

var (
	echoLambda *echoadapter.EchoLambdaV2
)

func init() {
	log := logger.Get(os.Stdout)

	log.With().
		Timestamp().
		Str("service", "web").
		Logger().
		Level(zerolog.InfoLevel)

	var e = echo.New()

	logger := lecho.From(log)
	e.Logger = logger

	e.Use(lecho.Middleware(lecho.Config{
		Logger: logger,
	}))

	err := web.Init(e)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Echo app")
	}
	echoLambda = echoadapter.NewV2(e)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
