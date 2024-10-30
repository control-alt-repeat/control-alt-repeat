package internal

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func LookupItem(ctx context.Context, itemID string) (models.WarehouseItem, []warehouse.EbayItemInternal, error) {
	warehouseItem, err := warehouse.GetWarehouseItem(ctx, itemID)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	ebayInternalItems, err := warehouse.GetEbayInternalItems(ctx, warehouseItem.EbayListingIDs)
	if err != nil {
		return models.WarehouseItem{}, nil, err
	}

	return warehouseItem, ebayInternalItems, err
}
