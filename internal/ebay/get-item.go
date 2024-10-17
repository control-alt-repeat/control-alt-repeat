package ebay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func GetItem(ebayListingID string) (*models.EbayItem, error) {
	fmt.Printf("Getting eBay listing with ID: %s\n", ebayListingID)

	request, requesterCredentials, err := newTraditionalAPIRequest("GetItem")
	if err != nil {
		fmt.Println(err)
		return &models.EbayItem{}, err
	}

	payload := models.GetItemRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		ItemID:               ebayListingID,
	}

	resp, err := request.Post(payload)
	if err != nil {
		fmt.Println(err)
		return &models.EbayItem{}, err
	}

	var getItemResponse models.GetItemResponse
	err = xml.Unmarshal(resp, &getItemResponse)
	if err != nil {
		fmt.Println(err)
		return &models.EbayItem{}, err
	}

	if getItemResponse.Ack != "Success" {
		err = errors.New(getItemResponse.Errors.LongMessage)
	}

	return &getItemResponse.Item, err
}
