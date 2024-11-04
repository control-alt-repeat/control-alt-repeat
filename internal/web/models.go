package web

import (
	"fmt"

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
}

type EbayReference struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL"`
	ListingURL  string `json:"listingURL"`
}

func Map(i models.WarehouseItem) Item {
	return Item{
		ID:               i.ControlAltRepeatID,
		Title:            i.Title,
		Shelf:            i.Shelf,
		FreeagentOwnerID: i.FreeagentOwnerID,
		OwnerDisplayName: i.OwnerDisplayName,
		EbayListingURL:   fmt.Sprintf("https://www.ebay.co.uk/itm/%s", i.EbayListingID),
		EbayReferences:   []EbayReference{},
	}
}
