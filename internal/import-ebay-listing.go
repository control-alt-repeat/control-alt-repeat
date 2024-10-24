package internal

import (
	"fmt"
	"time"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func ImportEbayListing(ebayListing *models.EbayItem) error {
	if ebayListing.SKU != "" {
		err := warehouse.ValidateSKU(ebayListing.SKU)
		if err != nil {
			return err
		}
	}

	warehouseItem := &warehouse.WarehouseItem{}
	warehouseItem.EbayListingIDs = []string{ebayListing.ItemID}

	warehouseItem.InitialiseFromSKU(ebayListing.SKU)

	if warehouseItem.ControlAltRepeatID == "" {
		warehouseItem.ControlAltRepeatID = warehouse.GenerateControlAltRepeatID()

		newSKU := warehouseItem.ToEbaySKU()

		ebay.ReviseSKU(ebayListing.ItemID, newSKU)
	}

	startTime, err := time.Parse(time.RFC3339, ebayListing.ListingDetails.StartTime)
	if err != nil {
		return err
	}

	ebayItemInternal := &warehouse.EbayItemInternal{
		ID:          ebayListing.ItemID,
		Title:       ebayListing.Title,
		PictureURL:  ebayListing.PictureDetails.PictureURL[0],
		ViewItemURL: ebayListing.ListingDetails.ViewItemURL,
		StartTime:   startTime,
	}

	warehouseItem.AddedTime = ebayItemInternal.StartTime

	err = aws.SaveJsonObjectS3(
		warehouse.EbayListingsBucketName,
		ebayItemInternal.ID,
		ebayItemInternal,
	)
	if err != nil {
		fmt.Printf("Failed to save eBay listing '%s'\n", ebayItemInternal.ID)
		return err
	}

	err = aws.SaveJsonObjectS3(
		warehouse.WarehouseItemsBucketName,
		warehouseItem.ControlAltRepeatID,
		warehouseItem,
	)
	if err != nil {
		fmt.Printf("Failed to save warehouse item '%s'\n", warehouseItem.ControlAltRepeatID)
		return err
	}

	fmt.Printf("Successfully imported eBay listing %s with ID %s\n", ebayListing.ItemID, warehouseItem.ControlAltRepeatID)

	return nil
}

func ImportEbayListingByID(ebayListingID string) error {
	err := warehouse.ValidateListingID(ebayListingID)
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
		"SKU",
	})
	if err != nil {
		return err
	}

	return ImportEbayListing(ebayListing)
}
