package internal

import (
	"context"
	"fmt"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func LookupItem(ctx context.Context, itemID string) (models.WarehouseItem, []models.WarehouseEbayListing, error) {
	warehouseItem, exists, err := warehouse.LoadItem(ctx, itemID)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	if !exists {
		return models.WarehouseItem{}, nil, fmt.Errorf("item does not exist for ID '%s'", itemID)
	}

	fmt.Println("SHOULD NOT SEE THIS")

	ebayListing, err := warehouse.LoadEbayListing(ctx, warehouseItem.EbayListingID)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	return warehouseItem, []models.WarehouseEbayListing{ebayListing}, err
}
