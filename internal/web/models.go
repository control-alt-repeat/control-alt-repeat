package web

type Item struct {
	ID             string          `json:"id"`
	Shelf          string          `json:"shelf"`
	EbayReferences []EbayReference `json:"ebayReferences"`
}

type EbayReference struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL"`
	ListingURL  string `json:"listingURL"`
}
