package models

import (
	"fmt"
	"regexp"
	"time"
)

type WarehouseItem struct {
	ControlAltRepeatID string
	Title              string
	Shelf              string
	AddedTime          time.Time
	EbayListingID      string
	PictureURL         string
	FreeagentOwnerID   string
	OwnerDisplayName   string
}

const SKURegex = `(?m)^\(?(?P<Shelf>\d+[A-Z]+)?\)?( )?(?P<ID>[A-Z]{3}\-[0-9]{3})?`

func (i WarehouseItem) ToEbaySKU() string {
	return fmt.Sprintf("(%s) %s", i.Shelf, i.ControlAltRepeatID)
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
