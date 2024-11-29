package persistence

import (
	"context"
	"time"

	models "github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/dynamodb"
)

func OverwriteItem(ctx context.Context, item models.WarehouseItem) error {
	return dynamodb.OverwriteItem(ctx, dynamodb.FromWarehouseItem(item))
}

func GetItem(ctx context.Context, id string) (models.WarehouseItem, error) {
	result, err := dynamodb.GetItem(ctx, id)
	if err != nil {
		return models.WarehouseItem{}, err
	}

	return result.Map(), nil
}

func GetUnshelvedItems(ctx context.Context) ([]models.WarehouseItem, error) {
	result, err := dynamodb.GetUnshelvedItems(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]models.WarehouseItem, len(result))
	for i, item := range result {
		items[i] = item.Map()
	}
	return items, nil
}

func GetItemsUpdatedSince(ctx context.Context, since time.Time) ([]models.WarehouseItem, error) {
	result, err := dynamodb.GetItemsUpdatedSince(ctx, since)
	if err != nil {
		return nil, err
	}

	items := make([]models.WarehouseItem, len(result))
	for i, item := range result {
		items[i] = item.Map()
	}
	return items, nil
}

func UpdateOwner(ctx context.Context, itemID, newOwnerID, newOwnerName string) error {
	return dynamodb.UpdateOwner(ctx, itemID, newOwnerID, newOwnerName)
}

func UpdateShelf(ctx context.Context, itemID, newShelf string) error {
	return dynamodb.UpdateShelf(ctx, itemID, newShelf)
}
