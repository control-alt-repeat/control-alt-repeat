package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/control-alt-repeat/control-alt-repeat/internal/ebay"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func MoveItem(ctx context.Context, itemID string, newShelf string) error {
	var warehouseItem models.WarehouseItem

	warehouseItem, exists, err := warehouse.LoadItem(ctx, itemID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("item does not exist for ID '%s'", itemID)
	}

	warehouseItem.Shelf = newShelf

	err = ebay.ReviseSKU(ctx, warehouseItem.EbayListingID, warehouseItem.ToEbaySKU())
	if err != nil {
		slog.Warn("unable to update SKU on eBay.")
		slog.Warn(err.Error())
	}

	return warehouse.SaveItem(ctx, warehouseItem)
}
