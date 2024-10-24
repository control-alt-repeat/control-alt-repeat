package warehouse

import (
	"fmt"
	"time"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

const EbayListingsBucketName = "control-alt-repeat-ebay-listings"

type EbayItemInternal struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	PictureURL  string    `json:"pictureURL"`
	ViewItemURL string    `json:"viewItemURL"`
	StartTime   time.Time `json:"startTime"`
}

func GetEbayInternalItems(ebayListingIDs []string) ([]EbayItemInternal, error) {
	var ebayInternalItems []EbayItemInternal

	for _, ebayListingID := range ebayListingIDs {
		var ebayItem EbayItemInternal
		fmt.Printf("Loading item '%s' from ebay listings\n", ebayListingID)
		err := aws.LoadJsonObjectS3(EbayListingsBucketName, ebayListingID, &ebayItem)
		if err != nil {
			return []EbayItemInternal{}, err
		}

		ebayInternalItems = append(ebayInternalItems, ebayItem)
	}

	return ebayInternalItems, nil
}
