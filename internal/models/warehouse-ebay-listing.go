package models

import "time"

type WarehouseEbayListing struct {
	ID          string
	Title       string
	PictureURL  string
	ViewItemURL string
	StartTime   time.Time
}
