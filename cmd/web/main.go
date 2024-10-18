package main

import "github.com/Control-Alt-Repeat/control-alt-repeat/internal/web"

func main() {
	e := web.Init()

	e.Logger.Fatal(e.Start(":8080"))
}
