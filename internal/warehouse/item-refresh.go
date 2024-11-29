package warehouse

import (
	"context"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/traditionalapi"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func RefreshItemFromEbay(ctx context.Context, itemID string) error {
	warehouseItem, err := persistence.GetItem(ctx, itemID)
	if err != nil {
		return err
	}

	ebayItemInternal, err := persistence.LoadEbayListing(ctx, persistence.LoadEbayListingOptions{ID: warehouseItem.EbayListingID})
	if err != nil {
		return err
	}

	ebayItemSource, err := traditionalapi.GetItem(ctx, warehouseItem.EbayListingID, []string{
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
