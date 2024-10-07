package internal

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	aws "github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
)

const (
	SKURegex                 = `(?m)^\(?(?P<Shelf>\d+[A-Z]+)?\)?( )?(?P<ID>[A-Z]{3}\-[0-9]{3})?`
	EbayListingsBucketName   = "control-alt-repeat-ebay-listings"
	WarehouseItemsBucketName = "control-alt-repeat-warehouse"
)

type WarehouseItem struct {
	ControlAltRepeatID string `json:"controlAltRepeatID"`
	Shelf              string `json:"shelf"`
	Ebay               struct {
		ID string `json:"id"`
	} `json:"ebay"`
}

type EbayListingItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func ImportEbayListing(ebayListingID string) error {
	err := validateListingID(ebayListingID)
	if err != nil {
		return err
	}

	// err = checkListingAlreadyImported(ebayListingID)
	// if err != nil {
	// 	return err
	// }

	fmt.Printf("Importing eBay listing with ID: %s\n", ebayListingID)

	ebayListing, err := ebay.GetItem(ebayListingID)
	if err != nil {
		return err
	}

	if ebayListing.SKU != "" {
		err = validateSKU(ebayListing.SKU)
		if err != nil {
			return err
		}
	}

	warehouseItem := &WarehouseItem{}
	warehouseItem.Ebay.ID = ebayListing.ItemID

	warehouseItem.initialiseFromSKU(ebayListing.SKU)

	if warehouseItem.ControlAltRepeatID == "" {
		warehouseItem.ControlAltRepeatID = generateControlAltRepeatID()

		newSKU := warehouseItem.toEbaySKU()

		ebay.ReviseSKU(ebayListing.ItemID, newSKU)
	}

	ebayListingItem := &EbayListingItem{
		ID:    ebayListing.ItemID,
		Title: ebayListing.Title,
	}

	err = aws.SaveJsonObjectS3(
		EbayListingsBucketName,
		ebayListingItem.ID,
		ebayListingItem,
	)
	if err != nil {
		fmt.Printf("Failed to save eBay listing '%s'\n", ebayListingItem.ID)
		return err
	}

	err = aws.SaveJsonObjectS3(
		WarehouseItemsBucketName,
		warehouseItem.ControlAltRepeatID,
		warehouseItem,
	)
	if err != nil {
		fmt.Printf("Failed to save warehouse item '%s'\n", warehouseItem.ControlAltRepeatID)
		return err
	}

	fmt.Printf("Successfully imported eBay listing %s with ID %s\n", ebayListingID, warehouseItem.ControlAltRepeatID)

	return err
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
