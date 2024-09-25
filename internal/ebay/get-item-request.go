package ebay

import "encoding/xml"

type GetItemRequest struct {
	XMLName              xml.Name             `xml:"GetItemRequest"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	ItemID               string               `xml:"ItemID"`
}

type RequesterCredentials struct {
	EBayAuthToken string `xml:"eBayAuthToken"`
}
