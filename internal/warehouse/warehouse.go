package warehouse

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

const (
	SKURegex                 = `(?m)^\(?(?P<Shelf>\d+[A-Z]+)?\)?( )?(?P<ID>[A-Z]{3}\-[0-9]{3})?`
	WarehouseItemsBucketName = "control-alt-repeat-warehouse"
)

type WarehouseItem struct {
	ControlAltRepeatID string    `json:"controlAltRepeatID"`
	Shelf              string    `json:"shelf"`
	AddedTime          time.Time `json:"addedTime"`
	EbayListingIDs     []string  `json:"ebayListingIDs"`
}

func GetWarehouseItem(ctx context.Context, itemID string) (WarehouseItem, error) {
	var warehouseItem WarehouseItem

	err := aws.LoadJsonObjectS3(ctx, WarehouseItemsBucketName, itemID, &warehouseItem)
	if err != nil {
		return warehouseItem, err
	}

	return warehouseItem, nil
}

func GenerateControlAltRepeatID() (string, error) {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers := []rune("0123456789")

	var result strings.Builder
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			fmt.Println("Error generating secure random number:", err)
			return "", err
		}
		result.WriteRune(letters[n.Int64()])
	}
	result.WriteRune('-')
	for i := 0; i < 3; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(numbers))))
		if err != nil {
			fmt.Println("Error generating secure random number:", err)
			return "", err
		}
		result.WriteRune(numbers[n.Int64()])
	}

	return result.String(), nil
}

func (wi *WarehouseItem) InitialiseFromSKU(sku string) {
	re := regexp.MustCompile(SKURegex)

	match := re.FindStringSubmatch(sku)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			if name == "Shelf" {
				wi.Shelf = match[i]
			}
			if name == "ID" {
				wi.ControlAltRepeatID = match[i]
			}
		}
	}
}

func ValidateSKU(sku string) error {
	re := regexp.MustCompile(SKURegex)

	if !re.MatchString(sku) {
		return fmt.Errorf("the SKU from eBay '%s' does not match regular expression - fix before retrying", sku)
	}

	return nil
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

func (i WarehouseItem) ToEbaySKU() string {
	return fmt.Sprintf("(%s) %s", i.Shelf, i.ControlAltRepeatID)
}
