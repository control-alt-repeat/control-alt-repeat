package ebay

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type TraditionalAPIRequest struct {
	HttpRequest http.Request
	Payload     TraditionalAPIRequestPayload
}

type TraditionalAPIRequestPayload interface {
	SetEbayAccessToken()
}

func (s *RequesterCredentials) SetEBayAuthToken() error {
	ebayAccessToken, err := getAccessToken()

	if err != nil {
		return err
	}

	s.EBayAuthToken = ebayAccessToken

	fmt.Println("Access token added to 'EBayAuthToken'")

	return nil
}

type RequesterCredentials struct {
	EBayAuthToken string `xml:"eBayAuthToken"`
}

func newTraditionalAPIRequest(callName string, payload interface{}, requesterCredentials RequesterCredentials) (*TraditionalAPIRequest, error) {
	err := requesterCredentials.SetEBayAuthToken()
	if err != nil {
		return nil, err
	}

	traditionalAPIRequest := &TraditionalAPIRequest{}

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
