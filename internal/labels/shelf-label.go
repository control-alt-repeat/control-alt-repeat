package labels

import (
	"context"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
)

func CreateShelfLabelFromItem(ctx context.Context, item models.WarehouseItem) (label []byte, name string, err error) {
	return Create102x152mmItemLabel(
		item.ControlAltRepeatID,
		item.Title,
		strings.Join([]string{"https://www.ebay.co.uk/itm", item.EbayListingID}, "/"),
	)
}
