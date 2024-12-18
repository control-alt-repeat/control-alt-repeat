package dynamodb

import (
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
)

const UnsetShelfDefault = "NOT_SET_YET"

type Item struct {
	ID               string `dynamodbav:"ID"`
	Shelf            string `dynamodbav:"Shelf"`
	Title            string `dynamodbav:"Title"`
	PictureURL       string `dynamodbav:"PictureURL"`
	EbayListingID    string `dynamodbav:"EbayListingID"`
	FreeagentOwnerID string `dynamodbav:"FreeagentOwnerID"`
	OwnerDisplayName string `dynamodbav:"OwnerDisplayName"`
	CreatedAt        int64  `dynamodbav:"CreatedAt"`
	UpdatedAt        int64  `dynamodbav:"UpdatedAt"`
}

func (i Item) Map() models.WarehouseItem {
	var shelf string
	if i.Shelf == UnsetShelfDefault {
		shelf = ""
	} else {
		shelf = i.Shelf
	}

	return models.WarehouseItem{
		ControlAltRepeatID: i.ID,
		Title:              i.Title,
		Shelf:              shelf,
		AddedTime:          time.Unix(i.CreatedAt, 0),
		EbayListingID:      i.EbayListingID,
		FreeagentOwnerID:   i.FreeagentOwnerID,
		OwnerDisplayName:   i.OwnerDisplayName,
	}
}

func FromWarehouseItem(item models.WarehouseItem) Item {
	return Item{
		ID:               item.ControlAltRepeatID,
		Shelf:            item.Shelf,
		Title:            item.Title,
		CreatedAt:        item.AddedTime.Unix(),
		UpdatedAt:        time.Now().Unix(),
		PictureURL:       item.PictureURL,
		EbayListingID:    item.EbayListingID,
		FreeagentOwnerID: item.FreeagentOwnerID,
		OwnerDisplayName: item.OwnerDisplayName,
	}
}
