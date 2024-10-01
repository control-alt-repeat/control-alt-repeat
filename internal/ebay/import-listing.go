package ebay

func ImportListing() error {
	_, err := GetItem("387372844761")
	return err
}
