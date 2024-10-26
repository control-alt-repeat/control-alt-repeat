package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/labstack/echo/v4"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
)

var log zerolog.Logger

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.With().
		Timestamp().
		Str("service", "web").
		Logger().
		Level(zerolog.DebugLevel)

	var e = echo.New()

	err := web.Init(e)

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
