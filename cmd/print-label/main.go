package main

import (
	"flag"
	"log"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
)

var (
	text     = flag.String("text", "JLJ-168", "The text to be printed")
	textArea = flag.String("textArea", "Sound Card Asus Xonar DG PCI 5.1 Audio Card ( not pci-e )", "Destination of QR code")
	qr       = flag.String("qr", "https://www.ebay.co.uk/itm/Sound-Card-Asus-Xonar-DG-PCI-5-1-Audio-Card-not-pci-e-/387463431314", "Destination of QR code")
)

func main() {
	flag.Parse()

	// imageBytes, err := labels.Create62mmItemLabel(*text, *qr)
	imageBytes, err := labels.Create102x152mmItemLabel(*text, *textArea, *qr)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("out.png", imageBytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
