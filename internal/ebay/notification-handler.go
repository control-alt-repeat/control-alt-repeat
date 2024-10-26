package ebay

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

type Notification struct {
	NotificationEventName string `xml:"NotificationEventName"`
	Timestamp             string `xml:"Timestamp"`
}

const notificationBucket = "control-alt-repeat-live-ebay-incoming-notifications"

func HandleNotification(ctx context.Context, notificationXml string) error {
	notificationBytes := []byte(notificationXml)

	var notification models.ItemNotificationEnvelope
	err := xml.Unmarshal(notificationBytes, &notification)
	if err != nil {
		return saveRawXml(ctx, notificationBytes)
	}

	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Format the time as a human-readable string
	timestamp := currentTime.Format("2006-01-02T15:04:05Z") // ISO 8601 format

	key := fmt.Sprintf("%s-%s.xml", timestamp, notification.Body.GetItemResponse.NotificationEventName)

	return aws.SaveBytesToS3(
		ctx,
		notificationBucket,
		key,
		notificationBytes,
		"application/xml",
	)
}

func saveRawXml(ctx context.Context, rawXml []byte) error {
	// Get the current time in UTC
	currentTime := time.Now().UTC()

	// Format the time as a human-readable string
	timestamp := currentTime.Format("2006-01-02T15:04:05Z.xml") // ISO 8601 format

	return aws.SaveBytesToS3(
		ctx,
		notificationBucket,
		timestamp,
		[]byte(rawXml),
		"application/xml",
	)
}
