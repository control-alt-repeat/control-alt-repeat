package ebay

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type TraditionalAPIRequest struct {
	CallName string
}

type TraditionalAPIRequestPayload interface {
	SetEbayAccessToken()
}

func (r *TraditionalAPIRequest) Post(ctx context.Context, payload interface{}) ([]byte, error) {
	xmlData, err := xml.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	xmlData = append([]byte(xml.Header), xmlData...)

	xmlString := string(xmlData)

	reader := strings.NewReader(xmlString)

	request, err := http.NewRequestWithContext(ctx, "POST", "https://api.ebay.com/ws/api.dll", reader)

	if err != nil {
		return nil, err
	}

	request.Header.Add("X-EBAY-API-COMPATIBILITY-LEVEL", "1193")
	request.Header.Add("X-EBAY-API-SITEID", "3")
	request.Header.Add("X-EBAY-API-CALL-NAME", r.CallName)

	client := &http.Client{}

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func newTraditionalAPIRequest(callName string) (*TraditionalAPIRequest, *RequesterCredentials, error) {
	traditionalAPIRequest := &TraditionalAPIRequest{
		CallName: callName,
	}
	requesterCredentials := &RequesterCredentials{}

	err := requesterCredentials.SetEBayAuthToken()

	return traditionalAPIRequest, requesterCredentials, err
}
