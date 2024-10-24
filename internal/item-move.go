package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/warehouse"
)

func MoveItem(itemID string, newShelf string) error {
	var warehouseItem warehouse.WarehouseItem

	err := aws.LoadJsonObjectS3(warehouse.WarehouseItemsBucketName, itemID, &warehouseItem)

	if err != nil {
		return err
	}

	fmt.Printf("Loaded item '%s' from warehouse\n", warehouseItem.ControlAltRepeatID)

	warehouseItem.Shelf = newShelf

	for _, ebayListingID := range warehouseItem.EbayListingIDs {
		err = ebay.ReviseSKU(ebayListingID, warehouseItem.ToEbaySKU())
		if err != nil {
			return err
		}

		err = aws.SaveJsonObjectS3(warehouse.WarehouseItemsBucketName, warehouseItem.ControlAltRepeatID, warehouseItem)
		if err != nil {
			return err
		}
	}
	return nil
}
