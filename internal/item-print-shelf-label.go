package internal

import (
	"fmt"
	"strings"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/labels"
)

func ItemPrintShelfLabel(itemID string) error {
	var warehouseItem WarehouseItem
	var ebayItem EbayItemInternal

	label, err := labels.Create102x152mmItemLabel(
		warehouseItem.ControlAltRepeatID,
		ebayItem.Title,
		strings.Join([]string{"https://www.ebay.co.uk/itm", ebayItem.ID}, "/"),
	)
	if err != nil {
		return err
	}

	err = aws.SaveBytesToS3("control-alt-repeat-label-print-buffer", fmt.Sprintf("102x152-%s.png", warehouseItem.ControlAltRepeatID), label, "image/png")
	if err != nil {
		return err
	}

	return labels.NotifyLabelPrintServer()
}
