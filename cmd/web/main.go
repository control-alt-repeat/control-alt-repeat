package main

import (
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"
)

func main() {
	e, err := web.Init()

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
