package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

type ItemLookupResult struct {
	ID          string
	Title       string
	ImageURL    string
	Description string
	EbayURL     string
}

func LookupItem(itemID string) (ItemLookupResult, error) {
	var warehouseItem WarehouseItem
	var ebayItem EbayListingItem
	var result ItemLookupResult

	err := aws.LoadJsonObjectS3(WarehouseItemsBucketName, itemID, &warehouseItem)
	if err != nil {
		return result, err
	}
	fmt.Printf("Loaded item '%s' from warehouse\n", warehouseItem.ControlAltRepeatID)

	err = aws.LoadJsonObjectS3(EbayListingsBucketName, warehouseItem.Ebay.ID, &ebayItem)
	if err != nil {
		return result, err
	}
	fmt.Printf("Loaded item '%s' from ebay listings\n", warehouseItem.Ebay.ID)

	ebayResult, err := ebay.GetItem(warehouseItem.Ebay.ID)
	if err != nil {
		return result, err
	}
	fmt.Printf("Loaded listing '%s' direct from ebay\n", warehouseItem.Ebay.ID)

	return ItemLookupResult{
		ID:          itemID,
		Title:       ebayItem.Title,
		ImageURL:    ebayResult.PictureDetails.PictureURL[0],
		Description: ebayResult.ConditionDisplayName,
		EbayURL:     ebayItem.ViewItemURL,
	}, nil
}
