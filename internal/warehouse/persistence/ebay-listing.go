package persistence

import (
	"context"

	models "github.com/control-alt-repeat/control-alt-repeat/internal/models"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse/persistence/s3"
)

type SaveEbayListingOptions struct {
	EbayListing models.WarehouseEbayListing
}

func SaveEbayListing(ctx context.Context, opt SaveEbayListingOptions) error {
	return s3.SaveEbayListing(ctx, s3.SaveEbayListingOptions{EbayListing: s3.EbayListing{
		ID:                 opt.EbayListing.ID,
		Title:              opt.EbayListing.Title,
		PictureURL:         opt.EbayListing.PictureURL,
		ViewEbayListingURL: opt.EbayListing.ViewItemURL,
		StartTime:          opt.EbayListing.StartTime,
	}})
}

type LoadEbayListingOptions struct {
	ID string
}

func LoadEbayListing(ctx context.Context, opt LoadEbayListingOptions) (models.WarehouseEbayListing, error) {
	result, err := s3.LoadEbayListing(ctx, s3.LoadEbayListingOptions{ID: opt.ID})

	return models.WarehouseEbayListing{
		ID:          result.ID,
		Title:       result.Title,
		PictureURL:  result.PictureURL,
		ViewItemURL: result.ViewEbayListingURL,
		StartTime:   result.StartTime,
	}, err
}
