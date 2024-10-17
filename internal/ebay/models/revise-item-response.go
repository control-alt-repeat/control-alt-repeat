package models

import "encoding/xml"

type ReviseItemResponse struct {
	XMLName   xml.Name `xml:"ReviseItemResponse"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Timestamp string   `xml:"Timestamp"`
	Ack       string   `xml:"Ack"`
	Errors    struct {
		Text            string `xml:",chardata"`
		ShortMessage    string `xml:"ShortMessage"`
		LongMessage     string `xml:"LongMessage"`
		ErrorCode       string `xml:"ErrorCode"`
		SeverityCode    string `xml:"SeverityCode"`
		ErrorParameters struct {
			Text    string `xml:",chardata"`
			ParamID string `xml:"ParamID,attr"`
			Value   string `xml:"Value"`
		} `xml:"ErrorParameters"`
		ErrorClassification string `xml:"ErrorClassification"`
	} `xml:"Errors"`
	Version               string `xml:"Version"`
	Build                 string `xml:"Build"`
	HardExpirationWarning string `xml:"HardExpirationWarning"`
	ItemID                string `xml:"ItemID"`
	StartTime             string `xml:"StartTime"`
	EndTime               string `xml:"EndTime"`
	Fees                  struct {
		Text string `xml:",chardata"`
		Fee  []struct {
			Text string `xml:",chardata"`
			Name string `xml:"Name"`
			Fee  struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"Fee"`
		} `xml:"Fee"`
	} `xml:"Fees"`
	DiscountReason string `xml:"DiscountReason"`
}
