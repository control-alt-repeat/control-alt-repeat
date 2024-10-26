package models

import "encoding/xml"

type SetNotificationPreferencesRequest struct {
	XMLName              xml.Name             `xml:"SetNotificationPreferencesRequest"`
	Text                 string               `xml:",chardata"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	// DeliveryURLName                string                         `xml:"DeliveryURLName"`
	ApplicationDeliveryPreferences ApplicationDeliveryPreferences `xml:"ApplicationDeliveryPreferences"`
	UserDeliveryPreferenceArray    UserDeliveryPreferenceArray    `xml:"UserDeliveryPreferenceArray"`

	// Not interested in this right now
	// EventProperty   struct {
	// 	Text      string `xml:",chardata"`
	// 	EventType string `xml:"EventType"`
	// 	Name      string `xml:"Name"`
	// 	Value     string `xml:"Value"`
	// } `xml:"EventProperty"`

	// We can use this to pass UserData so that we can reconcile notifications with a user on our side
	// UserData struct {
	// 	Text             string `xml:",chardata"`
	// 	ExternalUserData string `xml:"ExternalUserData"`
	// } `xml:"UserData"`

	// Optional. English is OK here
	// ErrorLanguage string `xml:"ErrorLanguage"`

	// Returned back to us as CorrelationID (for tracing requests)
	// MessageID     string `xml:"MessageID"`

	// Optional. I'd rather receive latest
	// Version       string `xml:"Version"`
	WarningLevel string `xml:"WarningLevel"`
}

type ApplicationDeliveryPreferences struct {
	Text string `xml:",chardata"`
	// AlertEmail        string `xml:"AlertEmail"`
	// AlertEnable       string `xml:"AlertEnable"`
	ApplicationEnable string `xml:"ApplicationEnable"`
	ApplicationURL    string `xml:"ApplicationURL"`

	// Not sure when we'll need this
	// DeliveryURLDetails struct {
	// 	Text            string `xml:",chardata"`
	// 	DeliveryURL     string `xml:"DeliveryURL"`
	// 	DeliveryURLName string `xml:"DeliveryURLName"`
	// 	Status          string `xml:"Status"`
	// } `xml:"DeliveryURLDetails"`
	DeviceType string `xml:"DeviceType"`

	// Optional. I'd rather receive latest
	// PayloadVersion string `xml:"PayloadVersion"`
}

type UserDeliveryPreferenceArray struct {
	Text               string               `xml:",chardata"`
	NotificationEnable []NotificationEnable `xml:"NotificationEnable"`
}

type NotificationEnable struct {
	Text        string `xml:",chardata"`
	EventEnable string `xml:"EventEnable"`
	EventType   string `xml:"EventType"`
}

type SetNotificationPreferencesResponse struct {
	XMLName   xml.Name `xml:"SetNotificationPreferencesResponse"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Timestamp string   `xml:"Timestamp"`
	Ack       string   `xml:"Ack"`
	Version   string   `xml:"Version"`
	Build     string   `xml:"Build"`
	Errors    struct {
		Text                string `xml:",chardata"`
		ShortMessage        string `xml:"ShortMessage"`
		LongMessage         string `xml:"LongMessage"`
		ErrorCode           string `xml:"ErrorCode"`
		SeverityCode        string `xml:"SeverityCode"`
		ErrorClassification string `xml:"ErrorClassification"`
	} `xml:"Errors"`
}
