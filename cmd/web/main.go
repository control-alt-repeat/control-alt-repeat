package main

import (
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
	"github.com/labstack/echo/v4"
)

func main() {
	var e = echo.New()

	err := web.Init(e)

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
