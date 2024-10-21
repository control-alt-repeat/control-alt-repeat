package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

type ItemLookupResult struct {
	ID          string
	Shelf       string
	Title       string
	ImageURL    string
	Description string
	EbayURL     string
}

func LookupItem(itemID string) (ItemLookupResult, error) {
	var warehouseItem WarehouseItem
	var ebayItem EbayItemInternal
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

	return ItemLookupResult{
		ID:       itemID,
		Shelf:    warehouseItem.Shelf,
		Title:    ebayItem.Title,
		ImageURL: ebayItem.PictureURL,
		EbayURL:  ebayItem.ViewItemURL,
	}, nil
}
