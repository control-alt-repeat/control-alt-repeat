package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func ImportEbayListing(ebayListingID string) error {
	err := validateListingID(ebayListingID)
	if err != nil {
		return err
	}

	fmt.Printf("Listing ID valid: %s\n", ebayListingID)

	ebayListing, err := ebay.GetItem(ebayListingID, []string{
		"ItemID",
		"Title",
		"Description",
		"PictureDetails",
		"ListingDetails",
	})
	if err != nil {
		return err
	}

	if ebayListing.SKU != "" {
		err = validateSKU(ebayListing.SKU)
		if err != nil {
			return err
		}
	}

	warehouseItem := &WarehouseItem{}
	warehouseItem.EbayListingIDs = []string{ebayListing.ItemID}

	warehouseItem.initialiseFromSKU(ebayListing.SKU)

	if warehouseItem.ControlAltRepeatID == "" {
		warehouseItem.ControlAltRepeatID = generateControlAltRepeatID()

		newSKU := warehouseItem.toEbaySKU()

		ebay.ReviseSKU(ebayListing.ItemID, newSKU)
	}

	EbayItemInternal := &EbayItemInternal{
		ID:          ebayListing.ItemID,
		Title:       ebayListing.Title,
		PictureURL:  ebayListing.PictureDetails.PictureURL[0],
		ViewItemURL: ebayListing.ListingDetails.ViewItemURL,
	}

	err = aws.SaveJsonObjectS3(
		EbayListingsBucketName,
		EbayItemInternal.ID,
		EbayItemInternal,
	)
	if err != nil {
		fmt.Printf("Failed to save eBay listing '%s'\n", EbayItemInternal.ID)
		return err
	}

	err = aws.SaveJsonObjectS3(
		WarehouseItemsBucketName,
		warehouseItem.ControlAltRepeatID,
		warehouseItem,
	)
	if err != nil {
		fmt.Printf("Failed to save warehouse item '%s'\n", warehouseItem.ControlAltRepeatID)
		return err
	}

	fmt.Printf("Successfully imported eBay listing %s with ID %s\n", ebayListingID, warehouseItem.ControlAltRepeatID)

	return err
}
