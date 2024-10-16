package main

import (
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func main() {
	err := ebay.SetNotificationPreferences()

	if err != nil {
		fmt.Println(err)
	}
}
