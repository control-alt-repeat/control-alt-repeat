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

const UnsetShelfDefault = "NOT_SET_YET"

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
		shelf = UnsetShelfDefault
	}

	return dynamodb.SaveItem(ctx, dynamodb.SaveItemOptions{Item: dynamodb.Item{
		ID:               opt.Item.ControlAltRepeatID,
		Shelf:            shelf,
		Title:            opt.Item.Title,
		PictureURL:       opt.Item.PictureURL,
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

type QueryItemsOptions struct {
	IndexName              string
	KeyConditionExpression string
	StringExpressionValues map[string]string
}

var UnshelvedItemsQuery = QueryItemsOptions{
	IndexName:              "GSI_Shelf",
	KeyConditionExpression: "Shelf = :shelf",
	StringExpressionValues: map[string]string{
		":shelf": UnsetShelfDefault,
	},
}

func ItemByIDQuery(id string) QueryItemsOptions {
	return QueryItemsOptions{
		KeyConditionExpression: "ID = :id",
		StringExpressionValues: map[string]string{
			":id": id,
		},
	}
}

func QueryItems(ctx context.Context, opt QueryItemsOptions) ([]models.WarehouseItem, error) {
	result, err := dynamodb.QueryItems(ctx, dynamodb.QueryItemsOptions{
		IndexName:              opt.IndexName,
		KeyConditionExpression: opt.KeyConditionExpression,
		StringExpressionValues: opt.StringExpressionValues,
	})
	if err != nil {
		return nil, err
	}

	items := []models.WarehouseItem{}

	for _, item := range result {
		items = append(items, item.Map())
	}

	return items, nil
}
