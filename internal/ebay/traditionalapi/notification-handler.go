package traditionalapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/control-alt-repeat/control-alt-repeat/internal/aws"
)

type Notification struct {
	NotificationEventName string `xml:"NotificationEventName"`
	Timestamp             string `xml:"Timestamp"`
}

const notificationBucket = "control-alt-repeat-live-ebay-incoming-notifications"

func HandleNotification(ctx context.Context, notificationXml string) error {
	notificationBytes := []byte(notificationXml)

	var notification ItemNotificationEnvelope
	if err := xml.Unmarshal(notificationBytes, &notification); err != nil {
		return saveRawXml(ctx, notificationBytes)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")

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
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z.xml")

	return aws.SaveBytesToS3(
		ctx,
		notificationBucket,
		timestamp,
		[]byte(rawXml),
		"application/xml",
	)
}
