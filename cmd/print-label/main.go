package main

import (
	"flag"
	"log"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
)

var (
	text = flag.String("text", "ABC-123", "The text to be printed")
	qr   = flag.String("qr", "https://controlaltrepeat.net", "Destination of QR code")
)

func main() {
	flag.Parse()

	// imageBytes, err := labels.Create62mmItemLabel(*text, *qr)
	imageBytes, err := labels.Create102x152mmItemLabel(*text, *qr)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("out.png", imageBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
