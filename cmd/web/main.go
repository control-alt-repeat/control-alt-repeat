package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/control-alt-repeat/control-alt-repeat/internal/web"
)

func main() {
	e, err := web.Init(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
