package warehouse

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence"
)

func GetWarehouseItem(ctx context.Context, itemID string) (models.WarehouseItem, error) {
	item, err := persistence.LoadItem(ctx, persistence.LoadItemOptions{ID: itemID})
	if err != nil {
		return models.WarehouseItem{}, err
	}

	return item, nil
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

func ValidateListingID(ebayListingID string) error {
	id, err := strconv.Atoi(ebayListingID)
	if err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("ebay listing ID does not look valid - should be a biggish number")
	}

	return nil
}
