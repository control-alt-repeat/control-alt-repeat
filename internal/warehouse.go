package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
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

func generateControlAltRepeatID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers := []rune("0123456789")

	var result strings.Builder
	for i := 0; i < 3; i++ {
		result.WriteRune(letters[r.Intn(len(letters))])
	}
	result.WriteRune('-')
	for i := 0; i < 3; i++ {
		result.WriteRune(numbers[r.Intn(len(numbers))])
	}

	return result.String()
}

func (wi *WarehouseItem) initialiseFromSKU(sku string) {
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

func validateSKU(sku string) error {
	re := regexp.MustCompile(SKURegex)

	if !re.MatchString(sku) {
		return fmt.Errorf("the SKU from eBay '%s' does not match regular expression - fix before retrying", sku)
	}

	return nil
}

func validateListingID(ebayListingID string) error {
	id, err := strconv.Atoi(ebayListingID)
	if err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("ebay listing ID does not look valid - should be a biggish number")
	}

	return nil
}

func (i WarehouseItem) toEbaySKU() string {
	return fmt.Sprintf("(%s) %s", i.Shelf, i.ControlAltRepeatID)
}
