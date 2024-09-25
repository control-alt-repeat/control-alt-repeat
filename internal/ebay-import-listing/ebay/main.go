package ebay

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ImportListing() error {
	url := "https://api.ebay.com/ws/api.dll"
	method := "POST"

	payload := GetItemRequest{
		Xmlns: "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: RequesterCredentials{
			EBayAuthToken: "abc123",
		},
		ItemID: "387372844761",
	}

	xmlData, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling XML:", err)
		return err
	}

	xmlData = append([]byte(xml.Header), xmlData...)

	xmlString := string(xmlData)
	fmt.Println("Marshalled XML:")
	fmt.Println(xmlString)

	reader := strings.NewReader(xmlString)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, reader)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("X-EBAY-API-COMPATIBILITY-LEVEL", "1193")
	req.Header.Add("X-EBAY-API-SITEID", "3")
	req.Header.Add("X-EBAY-API-CALL-NAME", "GetItem")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(body))

	return nil
}
