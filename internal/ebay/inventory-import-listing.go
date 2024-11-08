package ebay

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay/sell/inventory"
)

func InventoryImportListing(ctx context.Context, listingID string) error {
	return inventory.BulkMigrateListing(ctx, inventory.BulkMigrateListingOptions{
		RequestObject: inventory.BulkMigrateRequest{
			Requests: []inventory.BulkMigrateRequestListing{
				{
					ListingID: listingID,
				},
			},
		},
	})
}
