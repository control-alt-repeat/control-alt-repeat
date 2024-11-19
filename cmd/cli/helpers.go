package main

import (
	"fmt"
	"os"
)

func handleError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
