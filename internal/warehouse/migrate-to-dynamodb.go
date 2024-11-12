package warehouse

import (
	"context"
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/traditionalapi"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func MigrateToDynamoDb(ctx context.Context) error {
	return persistence.IterateItems(ctx, func(ctx context.Context, key string) error {
		fmt.Printf("Migrating %s\n", key)
		item, err := persistence.LoadItem(ctx, persistence.LoadItemOptions{ID: key})
		if err != nil {
			return err
		}

		ebayListingInternal, err := persistence.LoadEbayListing(ctx, persistence.LoadEbayListingOptions{
			ID: item.EbayListingID,
		})
		if err != nil {
			return err
		}

		if item.AddedTime.IsZero() {
			ebayItemSource, err := traditionalapi.GetItem(ctx, item.EbayListingID, []string{
				"ListingDetails",
			})
			if err != nil {
				return err
			}
			item.AddedTime, err = time.Parse(time.RFC3339, ebayItemSource.ListingDetails.StartTime)
			if err != nil {
				return err
			}
		}

		item.Title = ebayListingInternal.Title
		item.PictureURL = ebayListingInternal.PictureURL

		return persistence.SaveItem(ctx, persistence.SaveItemOptions{
			Item: item,
		})
	})
}
