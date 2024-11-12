package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/traditionalapi"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func ImportEbayListing(ctx context.Context, ebayListing *traditionalapi.EbayItem) (string, error) {
	if ebayListing.SKU != "" {
		err := models.ValidateSKU(ebayListing.SKU)
		if err != nil {
			return "", err
		}
	}

	startTime, err := time.Parse(time.RFC3339, ebayListing.ListingDetails.StartTime)
	if err != nil {
		return "", err
	}

	warehouseItem := &models.WarehouseItem{
		Title:         ebayListing.Title,
		PictureURL:    ebayListing.PictureDetails.PictureURL[0],
		AddedTime:     startTime,
		EbayListingID: ebayListing.ItemID,
	}

	warehouseItem.InitialiseFromSKU(ebayListing.SKU)

	if warehouseItem.ControlAltRepeatID == "" {
		newID, err := warehouse.GenerateControlAltRepeatID()
		if err != nil {
			return "", err
		}

		warehouseItem.ControlAltRepeatID = newID

		newSKU := warehouseItem.ToEbaySKU()

		if err := traditionalapi.ReviseSKU(ctx, ebayListing.ItemID, newSKU); err != nil {
			return "", err
		}
	}

	ebayItemInternal := &models.WarehouseEbayListing{
		ID:          ebayListing.ItemID,
		Title:       ebayListing.Title,
		PictureURL:  ebayListing.PictureDetails.PictureURL[0],
		ViewItemURL: ebayListing.ListingDetails.ViewItemURL,
		StartTime:   startTime,
	}

	warehouseItem.AddedTime = ebayItemInternal.StartTime

	err = warehouse.SaveEbayListing(ctx, *ebayItemInternal)
	if err != nil {
		fmt.Printf("Failed to save eBay listing '%s'\n", ebayItemInternal.ID)
		return "", err
	}

	err = warehouse.SaveItem(ctx, *warehouseItem)
	if err != nil {
		fmt.Printf("Failed to save warehouse item '%s'\n", warehouseItem.ControlAltRepeatID)
		return "", err
	}

	fmt.Printf("Successfully imported eBay listing %s with ID %s\n", ebayListing.ItemID, warehouseItem.ControlAltRepeatID)

	return warehouseItem.ControlAltRepeatID, nil
}

func ImportEbayListingByID(ctx context.Context, ebayListingID string) (string, error) {
	err := warehouse.ValidateListingID(ebayListingID)
	if err != nil {
		return "", err
	}

	fmt.Printf("Listing ID valid: %s\n", ebayListingID)

	ebayListing, err := traditionalapi.GetItem(ctx, ebayListingID, []string{
		"ItemID",
		"Title",
		"Description",
		"PictureDetails",
		"ListingDetails",
		"SKU",
	})
	if err != nil {
		return "", err
	}

	return ImportEbayListing(ctx, ebayListing)
}
