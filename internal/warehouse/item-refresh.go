package warehouse

import (
	"context"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func RefreshItemsFromEbay(ctx context.Context) error {
	return persistence.IterateItems(ctx, RefreshItemFromEbay)
}

func RefreshItemFromEbay(ctx context.Context, itemID string) error {
	warehouseItem, err := persistence.LoadItem(ctx, persistence.LoadItemOptions{ID: itemID})
	if err != nil {
		return err
	}

	ebayItemInternal, err := persistence.LoadEbayListing(ctx, persistence.LoadEbayListingOptions{ID: warehouseItem.EbayListingID})
	if err != nil {
		return err
	}

	ebayItemSource, err := ebay.GetItem(ctx, warehouseItem.EbayListingID, []string{
		"ItemID",
		"Title",
		"PictureDetails",
		"ListingDetails",
		"SKU",
	})
	if err != nil {
		return err
	}

	startTime, err := time.Parse(time.RFC3339, ebayItemSource.ListingDetails.StartTime)
	if err != nil {
		return err
	}

	ebayItemInternal.Title = ebayItemSource.Title
	ebayItemInternal.PictureURL = ebayItemSource.PictureDetails.PictureURL[0]
	ebayItemInternal.ViewItemURL = ebayItemSource.ListingDetails.ViewItemURL
	ebayItemInternal.StartTime = startTime

	return persistence.SaveEbayListing(ctx, persistence.SaveEbayListingOptions{EbayListing: ebayItemInternal})
}
