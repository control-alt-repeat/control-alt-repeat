package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal"
)

func main() {
	err := internal.MoveItem(os.Args[1], os.Args[2])

	if err != nil {
		fmt.Println(err)
	}
}

func checkStringIsUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	isURL := err == nil

	return isURL
}

func getListingIDFromURL(url string) (string, error) {
	// Regular expression to find the item number after "/itm/"
	re := regexp.MustCompile(`/itm/(\d+)`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("no item number found in the URL")
}
