package s3

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type Item struct {
	ControlAltRepeatID string    `json:"controlAltRepeatID"`
	Shelf              string    `json:"shelf"`
	AddedTime          time.Time `json:"addedTime"`
	EbayListingIDs     []string  `json:"ebayListingIDs"`
	EbayListingID      string    `json:"ebayListingID"`
	FreeagentOwnerID   string    `json:"freeagentOwnerID"`
	OwnerDisplayName   string    `json:"ownerDisplayName"`
}

type SaveItemOptions struct {
	Item Item
}

func SaveItem(ctx context.Context, opt SaveItemOptions) error {
	s3, err := getClient()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(opt.Item)
	if err != nil {
		return err
	}

	_, err = s3.client.PutObject(
		ctx,
		"control-alt-repeat-warehouse",
		opt.Item.ControlAltRepeatID,
		strings.NewReader(string(jsonData)),
		int64(len(jsonData)),
		minio.PutObjectOptions{})

	return err
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

func IterateS3Objects(ctx context.Context, f func(context.Context, string) error) error {
	s3, err := getClient()
	if err != nil {
		return err
	}

	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
	}

	for object := range s3.client.ListObjects(ctx, "control-alt-repeat-warehouse", opts) {
		if object.Err != nil {
			return err
		}

		err := f(ctx, object.Key)
		if err != nil {
			return err
		}
	}

	return nil
}
