package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/labels"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func ItemPrintShelfLabel(ctx context.Context, itemID string) error {
	warehouseItem, err := warehouse.LoadItem(ctx, itemID)
	if err != nil {
		return err
	}

	ebayItemInternal, err := warehouse.LoadEbayListing(ctx, warehouseItem.EbayListingID)
	if err != nil {
		return err
	}

	label, err := labels.Create102x152mmItemLabel(
		warehouseItem.ControlAltRepeatID,
		ebayItemInternal.Title,
		strings.Join([]string{"https://www.ebay.co.uk/itm", ebayItemInternal.ID}, "/"),
	)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("102x152-%s.png", warehouseItem.ControlAltRepeatID)

	return labels.UploadFileFromBytes(label, key)
}
