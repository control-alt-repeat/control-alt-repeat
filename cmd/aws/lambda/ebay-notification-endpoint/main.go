package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

var log zerolog.Logger

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.With().
		Timestamp().
		Str("service", "ebay-notification-endpoint").
		Logger().
		Level(zerolog.DebugLevel)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := ebay.HandleNotification(ctx, req.Body); err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			// this is our eBay notification endpoint. If they see errors one this they may cut us off.
			StatusCode: http.StatusOK,
		}, nil
	}

	log.Info().Msg("Responding that everything is OK.")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	lambda.Start(handler)
}
