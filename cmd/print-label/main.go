package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
)

var (
	text = flag.String("text", "ABC-123", "The text to be printed")
	qr   = flag.String("qr", "https://controlaltrepeat.net", "Destination of QR code")
)

func main() {
	flag.Parse()

	imageBytes, err := labels.Create62mmItemLabel(*text, *qr)
	if err != nil {
		fmt.Errorf("Failed to create 62mm label", err)
	}

	err = os.WriteFile("out.png", imageBytes, 0644)
	if err != nil {
		fmt.Errorf("Failed to write 62mm label to file", err)
	}
}
