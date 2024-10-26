package main

import (
	"github.com/labstack/echo/v4"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
)

func main() {
	var e = echo.New()

	err := web.Init(e)

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
