package inventory

import (
	"context"
)

type BulkMigrateListingOptions struct {
	RequestObject BulkMigrateRequest
}

type BulkMigrateRequest struct {
	Requests []BulkMigrateRequestListing `json:"requests"`
}

type BulkMigrateRequestListing struct {
	ListingID string `json:"listingId"`
}

func BulkMigrateListings(ctx context.Context, opts BulkMigrateListingOptions) error {
	apiopts := requestOptions{
		Path: "/bulk_migrate_listing",
	}
	return apiPost(ctx, apiopts, opts.RequestObject)
}
