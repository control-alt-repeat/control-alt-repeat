package s3

import (
	"context"
	"encoding/json"
	"time"

	"github.com/minio/minio-go/v7"
)

type Item struct {
	ControlAltRepeatID string    `json:"controlAltRepeatID"`
	Shelf              string    `json:"shelf"`
	AddedTime          time.Time `json:"addedTime"`
	EbayListingIDs     []string  `json:"ebayListingIDs"`
}

type LoadItemOptions struct {
	ID string
}

func LoadItem(ctx context.Context, opt LoadItemOptions) (Item, error) {
	s3, err := getClient()
	if err != nil {
		return Item{}, err
	}

	obj, err := s3.client.GetObject(ctx, "control-alt-repeat-warehouse", opt.ID, minio.GetObjectOptions{})
	if err != nil {
		return Item{}, err
	}
	defer obj.Close()

	var item Item
	if err := json.NewDecoder(obj).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
