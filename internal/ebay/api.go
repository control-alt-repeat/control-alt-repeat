package ebay

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/sell/inventory"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/traditionalapi"
)

func GetItem(ctx context.Context, ebayListingID string, outputSelector []string) (*traditionalapi.EbayItem, error) {
	return traditionalapi.GetItem(ctx, ebayListingID, outputSelector)
}

func SetSKU(ctx context.Context, ebayListingID string, sku string) error {
	return traditionalapi.ReviseSKU(ctx, ebayListingID, sku)
}

func GetNotificationsUsage(ctx context.Context, itemID string) error {
	return traditionalapi.GetNotificationUsage(ctx, itemID)
}

func SetNotificationPreferences(ctx context.Context) error {
	return traditionalapi.SetNotificationPreferences(ctx)
}

func BulkMigrateListings(ctx context.Context, opts inventory.BulkMigrateListingOptions) error {
	return inventory.BulkMigrateListings(ctx, opts)
}

func CreateLocation(ctx context.Context, opts inventory.CreateLocationOptions) error {
	return inventory.CreateLocation(ctx, opts)
}

func GetLocations(ctx context.Context, opts inventory.GetLocationsOptions) (inventory.GetLocationsResponse, error) {
	return inventory.GetLocations(ctx, opts)
}

func UpdateLocation(ctx context.Context, opts inventory.UpdateLocationOptions) error {
	return inventory.UpdateLocation(ctx, opts)
}
