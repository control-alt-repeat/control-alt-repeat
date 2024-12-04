package ebay

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/sell/inventory"
)

func InventoryImportListing(ctx context.Context, listingID string) error {
	return BulkMigrateListings(ctx, inventory.BulkMigrateListingOptions{
		RequestObject: inventory.BulkMigrateRequest{
			Requests: []inventory.BulkMigrateRequestListing{
				{
					ListingID: listingID,
				},
			},
		},
	})
}
