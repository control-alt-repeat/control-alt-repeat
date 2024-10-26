package ebay

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func ReviseSKU(ctx context.Context, ebayListingID string, sku string) error {
	request, requesterCredentials, err := newTraditionalAPIRequest("ReviseItem")
	if err != nil {
		return err
	}

	payload := models.ReviseItemRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
	}
	payload.Item.ItemID = ebayListingID
	payload.Item.SKU = sku

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return err
	}

	var reviseItemResponse models.ReviseItemResponse
	if err = xml.Unmarshal(resp, &reviseItemResponse); err != nil {
		return err
	}

	if reviseItemResponse.Ack == "Failure" {
		return errors.New(reviseItemResponse.Errors.LongMessage)
	}

	if reviseItemResponse.Ack == "Warning" {
		fmt.Printf("eBay API Warning: %s\n", reviseItemResponse.Errors.LongMessage)
	}

	return nil
}
