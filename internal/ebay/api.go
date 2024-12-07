package ebay

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/developer/keymanagement"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/sell/finances"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/sell/inventory"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/signing"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/traditionalapi"
)

func GetItem(ctx context.Context, ebayListingID string, outputSelector []string) (*traditionalapi.EbayItem, error) {
	return traditionalapi.GetItem(ctx, ebayListingID, outputSelector)
}

func GetItemTransactions(ctx context.Context, ebayListingID string) (*traditionalapi.GetItemTransactionsResponse, error) {
	return traditionalapi.GetItemTransactions(ctx, ebayListingID)
}

func SetSKU(ctx context.Context, ebayListingID string, sku string) error {
	return traditionalapi.ReviseSKU(ctx, ebayListingID, sku)
}

func GetNotificationsUsage(ctx context.Context, itemID string) error {
	return traditionalapi.GetNotificationUsage(ctx, itemID)
}

func SetNotificationPreferences(ctx context.Context) error {
	return traditionalapi.SetNotificationPreferences(ctx, applicationDeliveryPreferences, userDeliveryPreferenceArray)
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

func GetTransaction(ctx context.Context, orderID string) (finances.GetTransactionResponse, error) {
	return finances.GetTransaction(ctx, orderID)
}

func CreateSigningKey(ctx context.Context) (keymanagement.CreateSigningKeyResponse, error) {
	return keymanagement.CreateSigningKey(ctx)
}

func VerifySignature(ctx context.Context) error {
	return signing.VerifyGet(ctx)
}
