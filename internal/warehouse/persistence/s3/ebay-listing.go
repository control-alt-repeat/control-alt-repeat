package s3

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type EbayListing struct {
	ID                 string    `json:"id"`
	Title              string    `json:"title"`
	PictureURL         string    `json:"pictureURL"`
	ViewEbayListingURL string    `json:"viewEbayListingURL"`
	StartTime          time.Time `json:"startTime"`
}

type SaveEbayListingOptions struct {
	EbayListing EbayListing
}

func SaveEbayListing(ctx context.Context, opt SaveEbayListingOptions) error {
	s3, err := getClient()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(opt.EbayListing)
	if err != nil {
		return err
	}

	_, err = s3.client.PutObject(
		ctx,
		"control-alt-repeat-ebay-listings",
		opt.EbayListing.ID,
		strings.NewReader(string(jsonData)),
		int64(len(jsonData)),
		minio.PutObjectOptions{})

	return err
}

type LoadEbayListingOptions struct {
	ID string
}

func LoadEbayListing(ctx context.Context, opt LoadEbayListingOptions) (EbayListing, error) {
	s3, err := getClient()
	if err != nil {
		return EbayListing{}, err
	}

	obj, err := s3.client.GetObject(ctx, "control-alt-repeat-ebay-listings", opt.ID, minio.GetObjectOptions{})
	if err != nil {
		return EbayListing{}, err
	}
	defer obj.Close()

	var ebayListing EbayListing
	if err := json.NewDecoder(obj).Decode(&ebayListing); err != nil {
		return EbayListing{}, err
	}

	return ebayListing, nil
}
