package models

import "encoding/xml"

type ReviseItemRequest struct {
	XMLName              xml.Name             `xml:"ReviseItemRequest"`
	Text                 string               `xml:",chardata"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	Item                 struct {
		Text   string `xml:",chardata"`
		ItemID string `xml:"ItemID"`
		SKU    string `xml:"SKU"`
	} `xml:"Item"`
}
