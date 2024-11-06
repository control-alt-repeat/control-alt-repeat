package main

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"

	"github.com/control-alt-repeat/control-alt-repeat/internal/logger"
	"github.com/control-alt-repeat/control-alt-repeat/internal/web"
)

func main() {
	log := logger.Get(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

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
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
