package warehouse

import (
	"context"
	"errors"
	"strconv"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func SaveEbayListing(ctx context.Context, ebayListing models.WarehouseEbayListing) error {
	return persistence.SaveEbayListing(ctx, persistence.SaveEbayListingOptions{EbayListing: ebayListing})
}

func LoadEbayListing(ctx context.Context, ebayListingID string) (models.WarehouseEbayListing, error) {
	return persistence.LoadEbayListing(ctx, persistence.LoadEbayListingOptions{ID: ebayListingID})
}

func ValidateListingID(ebayListingID string) error {
	id, err := strconv.Atoi(ebayListingID)
	if err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("ebay listing ID does not look valid - should be a biggish number")
	}

	return nil
}
