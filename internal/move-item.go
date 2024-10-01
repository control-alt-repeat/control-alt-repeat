package internal

import (
	"fmt"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

func MoveItem(itemID string, newShelf string) error {
	var warehouseItem WarehouseItem

	err := aws.LoadJsonObjectS3(WarehouseItemsBucketName, itemID, &warehouseItem)

	if err != nil {
		return err
	}

	fmt.Printf("Loaded from S3 warehouse item with ID '%s'\n", warehouseItem.ControlAltRepeatID)

	warehouseItem.Shelf = newShelf

	err = ebay.ReviseSKU(warehouseItem.Ebay.ID, warehouseItem.toEbaySKU())
	if err != nil {
		return err
	}

	aws.SaveJsonObjectS3(WarehouseItemsBucketName, warehouseItem.ControlAltRepeatID, warehouseItem)

	return nil
}
