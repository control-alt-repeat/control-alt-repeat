package ebay

import (
	"encoding/xml"
	"net/http"
	"strings"
)

type TraditionalAPIRequest struct {
	AccessToken string
	HttpRequest http.Request
}

func newTraditionalAPIRequest(callName string, payload interface{}) (*TraditionalAPIRequest, error) {
	traditionalAPIRequest := &TraditionalAPIRequest{}

	ebayAccessToken, err := getAccessToken()

	if err != nil {
		return nil, err
	}

	traditionalAPIRequest.AccessToken = ebayAccessToken

	xmlData, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	xmlData = append([]byte(xml.Header), xmlData...)

	xmlString := string(xmlData)

	reader := strings.NewReader(xmlString)

	req, err := http.NewRequest("POST", "https://api.ebay.com/ws/api.dll", reader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("X-EBAY-API-COMPATIBILITY-LEVEL", "1193")
	req.Header.Add("X-EBAY-API-SITEID", "3")
	req.Header.Add("X-EBAY-API-CALL-NAME", callName)

	traditionalAPIRequest.HttpRequest = *req

	return traditionalAPIRequest, nil
}
