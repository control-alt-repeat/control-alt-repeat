package persistence

import (
	"context"

	models "github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/s3"
)

type SaveItemOptions struct {
	Item models.WarehouseItem
}

func SaveItem(ctx context.Context, opt SaveItemOptions) error {
	return s3.SaveItem(ctx, s3.SaveItemOptions{Item: s3.Item{
		ControlAltRepeatID: opt.Item.ControlAltRepeatID,
		Shelf:              opt.Item.Shelf,
		AddedTime:          opt.Item.AddedTime,
		EbayListingIDs:     []string{opt.Item.EbayListingID},
	}})
}

type LoadItemOptions struct {
	ID string
}

func LoadItem(ctx context.Context, opt LoadItemOptions) (models.WarehouseItem, error) {
	result, err := s3.LoadItem(ctx, s3.LoadItemOptions{ID: opt.ID})

	return models.WarehouseItem{
		ControlAltRepeatID: result.ControlAltRepeatID,
		Shelf:              result.Shelf,
		AddedTime:          result.AddedTime,
		EbayListingID:      result.EbayListingIDs[0],
	}, err
}
