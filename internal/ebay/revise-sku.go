package ebay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func ReviseSKU(ebayListingID string, sku string) error {
	fmt.Printf("Updating SKU to '%s' on eBay listing '%s'\n", sku, ebayListingID)

	request, requesterCredentials, err := newTraditionalAPIRequest("ReviseItem")
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := models.ReviseItemRequest{
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

	var reviseItemResponse models.ReviseItemResponse
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
