package main

import (
	"fmt"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func main() {
	err := ebay.GetNotificationUsage(os.Args[1])

	if err != nil {
		fmt.Println(err)
	}
}
