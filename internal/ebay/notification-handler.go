package ebay

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

type Notification struct {
	NotificationEventName string `xml:"NotificationEventName"`
	Timestamp             string `xml:"Timestamp"`
}

const notificationBucket = "control-alt-repeat-live-ebay-incoming-notifications"

func HandleNotification(notificationXml string) error {
	notificationBytes := []byte(notificationXml)

	var notification Notification

	if err := xml.Unmarshal(notificationBytes, &notification); err != nil {
		log.Fatalf("Error unmarshalling XML: %v", err)
	}

	key := fmt.Sprintf("%s %s", notification.Timestamp, notification.NotificationEventName)

	fmt.Printf("Saving %s to %s\r", key, notificationBucket)

	return aws.SaveBytesToS3(
		notificationBucket,
		key,
		notificationBytes,
	)
}
