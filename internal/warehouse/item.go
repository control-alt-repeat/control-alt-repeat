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

func OverwriteItem(ctx context.Context, item models.WarehouseItem) error {
	return persistence.OverwriteItem(ctx, item)
}

func UpdateOwner(ctx context.Context, itemID, newOwnerID, newOwnerName string) error {
	return persistence.UpdateOwner(ctx, itemID, newOwnerID, newOwnerName)
}

func UpdateShelf(ctx context.Context, itemID, newShelf string) error {
	return persistence.UpdateShelf(ctx, itemID, newShelf)
}

func GetItem(ctx context.Context, itemID string) (models.WarehouseItem, error) {
	return persistence.GetItem(ctx, itemID)
}

func GetItemsUpdatedInLastYear(ctx context.Context) ([]models.WarehouseItem, error) {
	return persistence.GetItemsUpdatedSince(ctx, time.Now().Add(-365*24*time.Hour))
}

func LoadUnshelvedItems(ctx context.Context) ([]models.WarehouseItem, error) {
	return persistence.GetUnshelvedItems(ctx)
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
