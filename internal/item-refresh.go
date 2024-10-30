package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/aws"
	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func RefreshItemsFromEbay(ctx context.Context) error {
	return aws.IterateS3Objects(ctx, "control-alt-repeat-warehouse", "eu-west-2", RefreshItemFromEbay)
}

func RefreshItemFromEbay(ctx context.Context, itemID string) error {
	var ebayItemInternal warehouse.EbayItemInternal

	warehouseItem, err := warehouse.GetWarehouseItem(ctx, itemID)
	if err != nil {
		return err
	}

	for _, ebayListingID := range warehouseItem.EbayListingIDs {
		fmt.Println("loading ebay item ", ebayListingID)
		err = aws.LoadJsonObjectS3(ctx, warehouse.EbayListingsBucketName, ebayListingID, &ebayItemInternal)
		if err != nil {
			return err
		}

		ebayItemSource, err := ebay.GetItem(ctx, ebayListingID, []string{
			"ItemID",
			"Title",
			"PictureDetails",
			"ListingDetails",
			"SKU",
		})
		if err != nil {
			return err
		}

		startTime, err := time.Parse(ebayItemSource.ListingDetails.StartTime, time.RFC3339)
		if err != nil {
			return err
		}

		ebayItemInternal.Title = ebayItemSource.Title
		ebayItemInternal.PictureURL = ebayItemSource.PictureDetails.PictureURL[0]
		ebayItemInternal.ViewItemURL = ebayItemSource.ListingDetails.ViewItemURL
		ebayItemInternal.StartTime = startTime

		warehouseItem.EbayListingIDs = []string{ebayListingID}

		err = aws.SaveJsonObjectS3(ctx, warehouse.EbayListingsBucketName, ebayListingID, &ebayItemInternal)
		if err != nil {
			return err
		}
	}
	return nil
}
