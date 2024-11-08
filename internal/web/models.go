package web

import (
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/freeagent"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
)

type Item struct {
	ID               string          `json:"id"`
	Title            string          `json:"title"`
	Shelf            string          `json:"shelf"`
	EbayReferences   []EbayReference `json:"ebayReferences"`
	FreeagentOwnerID string          `json:"freeagentOwnerID"`
	OwnerDisplayName string          `json:"ownerDisplayName"`
	EbayListingURL   string          `json:"ebayListingURL"`
	ImageURL         string          `json:"imageURL"`
	AddedTime        time.Time       `json:"addedTime"`
}

type EbayReference struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL"`
	ListingURL  string `json:"listingURL"`
}

type Contact struct {
	Name        string `json:"name"`
	FreeagentID string `json:"id"`
}

func MapToWebItem(i models.WarehouseItem) Item {
	return Item{
		ID:               i.ControlAltRepeatID,
		Title:            i.Title,
		Shelf:            i.Shelf,
		FreeagentOwnerID: i.FreeagentOwnerID,
		OwnerDisplayName: i.OwnerDisplayName,
		EbayListingURL:   fmt.Sprintf("https://www.ebay.co.uk/itm/%s", i.EbayListingID),
		EbayReferences:   []EbayReference{},
		AddedTime:        i.AddedTime,
	}
}

func MapToWebContact(fc freeagent.Contact) Contact {
	return Contact{
		Name:        fc.DisplayName(),
		FreeagentID: fc.ID(),
	}
}

func MapSlice[T any, U any](input []T, mapper func(T) U) []U {
	result := make([]U, len(input))
	for i, item := range input {
		result[i] = mapper(item)
	}
	return result
}
