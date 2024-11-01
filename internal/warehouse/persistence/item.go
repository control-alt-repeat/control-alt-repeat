package persistence

import (
	"context"
	"time"

	models "github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/dynamodb"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/s3"
)

type SaveItemOptions struct {
	Item models.WarehouseItem
}

func SaveItem(ctx context.Context, opt SaveItemOptions) error {
	err := s3.SaveItem(ctx, s3.SaveItemOptions{Item: s3.Item{
		ControlAltRepeatID: opt.Item.ControlAltRepeatID,
		Shelf:              opt.Item.Shelf,
		AddedTime:          opt.Item.AddedTime,
		EbayListingIDs:     []string{opt.Item.EbayListingID},
		EbayListingID:      opt.Item.EbayListingID,
		FreeagentOwnerID:   opt.Item.FreeagentOwnerID,
		OwnerDisplayName:   opt.Item.OwnerDisplayName,
	}})
	if err != nil {
		return err
	}

	shelf := opt.Item.Shelf
	if shelf == "" {
		shelf = "NOT_SET_YET"
	}

	return dynamodb.SaveItem(ctx, dynamodb.SaveItemOptions{Item: dynamodb.Item{
		ID:               opt.Item.ControlAltRepeatID,
		Shelf:            shelf,
		EbayListingID:    opt.Item.EbayListingID,
		FreeagentOwnerID: opt.Item.FreeagentOwnerID,
		OwnerDisplayName: opt.Item.OwnerDisplayName,
		CreatedAt:        opt.Item.AddedTime.Unix(),
		UpdatedAt:        time.Now().Unix(),
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

func IterateItems(ctx context.Context, f func(context.Context, string) error) error {
	return s3.IterateS3Objects(ctx, f)
}
