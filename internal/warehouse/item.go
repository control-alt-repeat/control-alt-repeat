package warehouse

import (
	"context"
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func SaveItem(ctx context.Context, item models.WarehouseItem) error {
	return persistence.SaveItem(ctx, persistence.SaveItemOptions{Item: item})
}

func LoadItem(ctx context.Context, itemID string) (models.WarehouseItem, error) {
	items, err := persistence.QueryItems(ctx, persistence.ItemByIDQuery(itemID))
	if err != nil {
		return models.WarehouseItem{}, err
	}

	return items[0], err
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
