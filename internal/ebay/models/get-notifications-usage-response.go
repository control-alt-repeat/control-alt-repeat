package models

import "encoding/xml"

type GetNotificationsUsageResponse struct {
	XMLName               xml.Name `xml:"GetNotificationsUsageResponse"`
	Text                  string   `xml:",chardata"`
	Xmlns                 string   `xml:"xmlns,attr"`
	EndTime               string   `xml:"EndTime"`
	MarkUpMarkDownHistory struct {
		Text                string `xml:",chardata"`
		MarkUpMarkDownEvent struct {
			Text   string `xml:",chardata"`
			Reason string `xml:"Reason"`
			Time   string `xml:"Time"`
			Type   string `xml:"Type"`
		} `xml:"MarkUpMarkDownEvent"`
	} `xml:"MarkUpMarkDownHistory"`
	NotificationDetailsArray struct {
		Text                string `xml:",chardata"`
		NotificationDetails struct {
			Text            string `xml:",chardata"`
			DeliveryStatus  string `xml:"DeliveryStatus"`
			DeliveryTime    string `xml:"DeliveryTime"`
			DeliveryURL     string `xml:"DeliveryURL"`
			DeliveryURLName string `xml:"DeliveryURLName"`
			ErrorMessage    string `xml:"ErrorMessage"`
			ExpirationTime  string `xml:"ExpirationTime"`
			NextRetryTime   string `xml:"NextRetryTime"`
			Retries         string `xml:"Retries"`
			Type            string `xml:"Type"`
		} `xml:"NotificationDetails"`
	} `xml:"NotificationDetailsArray"`
	NotificationStatistics struct {
		Text               string `xml:",chardata"`
		DeliveredCount     string `xml:"DeliveredCount"`
		ErrorCount         string `xml:"ErrorCount"`
		ExpiredCount       string `xml:"ExpiredCount"`
		QueuedNewCount     string `xml:"QueuedNewCount"`
		QueuedPendingCount string `xml:"QueuedPendingCount"`
	} `xml:"NotificationStatistics"`
	StartTime     string `xml:"StartTime"`
	Ack           string `xml:"Ack"`
	Build         string `xml:"Build"`
	CorrelationID string `xml:"CorrelationID"`
	Errors        struct {
		Text                string `xml:",chardata"`
		ErrorClassification string `xml:"ErrorClassification"`
		ErrorCode           string `xml:"ErrorCode"`
		ErrorParameters     struct {
			Text    string `xml:",chardata"`
			ParamID string `xml:"ParamID,attr"`
			Value   string `xml:"Value"`
		} `xml:"ErrorParameters"`
		LongMessage  string `xml:"LongMessage"`
		SeverityCode string `xml:"SeverityCode"`
		ShortMessage string `xml:"ShortMessage"`
	} `xml:"Errors"`
	HardExpirationWarning string `xml:"HardExpirationWarning"`
	Timestamp             string `xml:"Timestamp"`
	Version               string `xml:"Version"`
}
