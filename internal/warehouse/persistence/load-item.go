package persistence

import (
	"context"

	models "github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/s3"
)

type LoadItemOptions struct {
	ID string
}

func LoadItem(ctx context.Context, opt LoadItemOptions) (models.WarehouseItem, error) {
	result, err := s3.LoadItem(ctx, s3.LoadItemOptions{ID: opt.ID})

	return models.WarehouseItem{
		ControlAltRepeatID: result.ControlAltRepeatID,
		Shelf:              result.Shelf,
		AddedTime:          result.AddedTime,
		EbayListingIDs:     result.EbayListingIDs,
	}, err
}
