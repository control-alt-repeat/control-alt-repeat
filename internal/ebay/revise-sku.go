package ebay

import (
	"encoding/xml"
	"errors"
	"fmt"
)

func ReviseSKU(ebayListingID string, sku string) error {
	fmt.Printf("Updating SKU to '%s' on eBay listing '%s'\n", sku, ebayListingID)

	request, requesterCredentials, err := newTraditionalAPIRequest("ReviseItem")
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := ReviseItemRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
	}
	payload.Item.ItemID = ebayListingID
	payload.Item.SKU = sku

	resp, err := request.Post(payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var reviseItemResponse ReviseItemResponse
	err = xml.Unmarshal(resp, &reviseItemResponse)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if reviseItemResponse.Ack == "Failure" {
		err = errors.New(reviseItemResponse.Errors.LongMessage)
	}

	if reviseItemResponse.Ack == "Warning" {
		fmt.Printf("eBay API Warning: %s\n", reviseItemResponse.Errors.LongMessage)
	}

	return err
}

// Models generated using https://blog.kowalczyk.info/tools/xmltogo/
type ReviseItemRequest struct {
	XMLName              xml.Name             `xml:"ReviseItemRequest"`
	Text                 string               `xml:",chardata"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	Item                 struct {
		Text   string `xml:",chardata"`
		ItemID string `xml:"ItemID"`
		SKU    string `xml:"SKU"`
	} `xml:"Item"`
}

type ReviseItemResponse struct {
	XMLName   xml.Name `xml:"ReviseItemResponse"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Timestamp string   `xml:"Timestamp"`
	Ack       string   `xml:"Ack"`
	Errors    struct {
		Text            string `xml:",chardata"`
		ShortMessage    string `xml:"ShortMessage"`
		LongMessage     string `xml:"LongMessage"`
		ErrorCode       string `xml:"ErrorCode"`
		SeverityCode    string `xml:"SeverityCode"`
		ErrorParameters struct {
			Text    string `xml:",chardata"`
			ParamID string `xml:"ParamID,attr"`
			Value   string `xml:"Value"`
		} `xml:"ErrorParameters"`
		ErrorClassification string `xml:"ErrorClassification"`
	} `xml:"Errors"`
	Version               string `xml:"Version"`
	Build                 string `xml:"Build"`
	HardExpirationWarning string `xml:"HardExpirationWarning"`
	ItemID                string `xml:"ItemID"`
	StartTime             string `xml:"StartTime"`
	EndTime               string `xml:"EndTime"`
	Fees                  struct {
		Text string `xml:",chardata"`
		Fee  []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name"`
			Fee  struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"Fee"`
		} `xml:"Fee"`
	} `xml:"Fees"`
	DiscountReason string `xml:"DiscountReason"`
}