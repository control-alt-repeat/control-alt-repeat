package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/traditionalapi"
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

	var pictureURL string
	if len(ebayListing.PictureDetails.PictureURL) > 0 {
		pictureURL = ebayListing.PictureDetails.PictureURL[0]
	} else {
		pictureURL = ebayListing.ProductListingDetails.StockPhotoURL
	}

	warehouseItem := &models.WarehouseItem{
		Title:         ebayListing.Title,
		PictureURL:    pictureURL,
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

		if err := ebay.SetSKU(ctx, ebayListing.ItemID, newSKU); err != nil {
			return "", err
		}
	}

	ebayItemInternal := &models.WarehouseEbayListing{
		ID:          ebayListing.ItemID,
		Title:       ebayListing.Title,
		PictureURL:  pictureURL,
		ViewItemURL: ebayListing.ListingDetails.ViewItemURL,
		StartTime:   startTime,
	}

	warehouseItem.AddedTime = ebayItemInternal.StartTime

	err = warehouse.SaveEbayListing(ctx, *ebayItemInternal)
	if err != nil {
		fmt.Printf("Failed to save eBay listing '%s'\n", ebayItemInternal.ID)
		return "", err
	}

	err = warehouse.OverwriteItem(ctx, *warehouseItem)
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

	ebayListing, err := ebay.GetItem(ctx, ebayListingID, []string{
		"ItemID",
		"Title",
		"Description",
		"PictureDetails",
		"ListingDetails",
		"ProductListingDetails",
		"SKU",
	})
	if err != nil {
		return "", err
	}

	return ImportEbayListing(ctx, ebayListing)
}
