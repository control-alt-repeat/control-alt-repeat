package ebay

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func GetItem(ctx context.Context, ebayListingID string, outputSelector []string) (*models.EbayItem, error) {
	request, requesterCredentials, err := newTraditionalAPIRequest("GetItem")
	if err != nil {
		fmt.Println(err)
		return &models.EbayItem{}, err
	}

	payload := models.GetItemRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		ItemID:               ebayListingID,
		OutputSelector:       outputSelector,
	}

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return &models.EbayItem{}, err
	}

	var getItemResponse models.GetItemResponse
	if err = xml.Unmarshal(resp, &getItemResponse); err != nil {
		return &models.EbayItem{}, err
	}

	if getItemResponse.Ack != "Success" {
		err = errors.New(getItemResponse.Errors.LongMessage)
	}

	return &getItemResponse.Item, err
}
