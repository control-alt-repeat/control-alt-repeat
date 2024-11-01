package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/control-alt-repeat/control-alt-repeat/internal/web"

	"github.com/labstack/echo/v4"
)

var log zerolog.Logger

var (
	echoLambda *echoadapter.EchoLambdaV2
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.With().
		Timestamp().
		Str("service", "web").
		Logger().
		Level(zerolog.DebugLevel)

	e := echo.New()
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
