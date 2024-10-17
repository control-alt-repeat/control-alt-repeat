package models

import "encoding/xml"

type GetNotificationsUsageRequest struct {
	XMLName              xml.Name             `xml:"GetNotificationsUsageRequest"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	Text                 string               `xml:",chardata"`
	Xmlns                string               `xml:"xmlns,attr"`
	// EndTime              string               `xml:"EndTime"`
	ItemID string `xml:"ItemID"`
	// StartTime            string               `xml:"StartTime"`
	// ErrorLanguage        string               `xml:"ErrorLanguage"`
	// MessageID            string               `xml:"MessageID"`
	// Version              string               `xml:"Version"`
	// WarningLevel         string               `xml:"WarningLevel"`
}
