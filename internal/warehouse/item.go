package warehouse

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func SaveItem(ctx context.Context, item models.WarehouseItem) error {
	return persistence.SaveItem(ctx, persistence.SaveItemOptions{Item: item})
}

func UpdateOwner(ctx context.Context, itemID, newOwnerID, newOwnerName string) error {
	return persistence.UpdateItem(ctx, persistence.UpdateItemOptions{
		ItemID: itemID,
		UpdateItemAttributes: []persistence.UpdateItemAttributes{
			{
				Name:  "FreeagentOwnerID",
				Value: newOwnerID,
			},
			{
				Name:  "OwnerDisplayName",
				Value: newOwnerName,
			},
		},
	})
}

func UpdateShelf(ctx context.Context, itemID, newShelf string) error {
	return persistence.UpdateItem(ctx, persistence.UpdateItemOptions{
		ItemID: itemID,
		UpdateItemAttributes: []persistence.UpdateItemAttributes{
			{
				Name:  "Shelf",
				Value: newShelf,
			},
		},
	})
}

func LoadItem(ctx context.Context, itemID string) (models.WarehouseItem, bool, error) {
	items, err := persistence.QueryItems(ctx, persistence.ItemByIDQuery(itemID))
	if err != nil {
		return models.WarehouseItem{}, false, err
	}

	if len(items) == 0 {
		return models.WarehouseItem{}, false, nil
	}

	return items[0], true, err
}

func LoadAllItems(ctx context.Context) ([]models.WarehouseItem, error) {
	return persistence.ScanItems(ctx, persistence.ItemsUpdatedSince(time.Now().Add(-365*24*time.Hour)))
}

func LoadUnshelvedItems(ctx context.Context) ([]models.WarehouseItem, error) {
	return persistence.QueryItems(ctx, persistence.UnshelvedItemsQuery)
}

func GenerateControlAltRepeatID() (string, error) {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers := []rune("0123456789")

	var result strings.Builder
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result.WriteRune(letters[n.Int64()])
	}
	result.WriteRune('-')
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbers))))
		if err != nil {
			return "", err
		}
		result.WriteRune(numbers[n.Int64()])
	}

	return result.String(), nil
}
