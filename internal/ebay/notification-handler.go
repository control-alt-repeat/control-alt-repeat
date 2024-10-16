package ebay

import (
	"time"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

type Notification struct {
	NotificationEventName string `xml:"NotificationEventName"`
	Timestamp             string `xml:"Timestamp"`
}

const notificationBucket = "control-alt-repeat-live-ebay-incoming-notifications"

func HandleNotification(notificationXml string) error {
	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Format the time as a human-readable string
	timestamp := currentTime.Format("2006-01-02T15:04:05Z") // ISO 8601 format

	return aws.SaveBytesToS3(
		notificationBucket,
		timestamp,
		[]byte(notificationXml),
	)
}
