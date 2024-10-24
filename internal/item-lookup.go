package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

type ItemLookupResult struct {
	ID             string
	Shelf          string
	EbayReferences []ItemLookupEbayReference
}

type ItemLookupEbayReference struct {
	Title       string
	Description string
	ImageURL    string
	ListingURL  string
}

func LookupItem(itemID string) (ItemLookupResult, error) {
	var warehouseItem warehouse.WarehouseItem

	fmt.Printf("Loading item '%s' from warehouse\n", itemID)
	err := aws.LoadJsonObjectS3(warehouse.WarehouseItemsBucketName, itemID, &warehouseItem)
	if err != nil {
		return ItemLookupResult{}, err
	}
	fmt.Printf("Loaded item '%s' from warehouse\n", itemID)

	var ebayReferences []ItemLookupEbayReference

	for _, ebayListingID := range warehouseItem.EbayListingIDs {
		var ebayItem warehouse.EbayItemInternal
		fmt.Printf("Loading item '%s' from ebay listings\n", ebayListingID)
		err = aws.LoadJsonObjectS3(warehouse.EbayListingsBucketName, ebayListingID, &ebayItem)
		if err != nil {
			return ItemLookupResult{}, err
		}

		ebayReferences = append(ebayReferences, ItemLookupEbayReference{
			Title:      ebayItem.Title,
			ImageURL:   ebayItem.PictureURL,
			ListingURL: ebayItem.ViewItemURL,
		})
	}

	return ItemLookupResult{
		ID:             itemID,
		Shelf:          warehouseItem.Shelf,
		EbayReferences: ebayReferences,
	}, nil
}
