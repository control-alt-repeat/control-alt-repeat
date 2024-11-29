package internal

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func LookupItem(ctx context.Context, itemID string) (models.WarehouseItem, []models.WarehouseEbayListing, error) {
	warehouseItem, err := warehouse.GetItem(ctx, itemID)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	ebayListing, err := warehouse.LoadEbayListing(ctx, warehouseItem.EbayListingID)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	return warehouseItem, []models.WarehouseEbayListing{ebayListing}, err
}
