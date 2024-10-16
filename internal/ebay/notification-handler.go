package ebay

import (
	"encoding/xml"
	"log"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

type Notification struct {
	NotificationEventName string `xml:"NotificationEventName"`
	Timestamp             string `xml:"Timestamp"`
}

func HandleNotification(notificationXml string) error {
	notificationBytes := []byte(notificationXml)

	var notification Notification

	if err := xml.Unmarshal(notificationBytes, &notification); err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	return aws.SaveBytesToS3(
		"control-alt-repeat-live-ebay-incoming-notifications",
		notification.Timestamp+" "+notification.NotificationEventName,
		notificationBytes,
	)
}
