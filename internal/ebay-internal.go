package internal

const EbayListingsBucketName = "control-alt-repeat-ebay-listings"

type EbayItemInternal struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	PictureURL  string `json:"pictureURL"`
	ViewItemURL string `json:"viewItemURL"`
}
