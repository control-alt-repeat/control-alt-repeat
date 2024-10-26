package ebay

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func GetNotificationUsage(ctx context.Context, itemID string) error {
	request, requesterCredentials, err := newTraditionalAPIRequest("GetNotificationsUsage")
	if err != nil {
		return err
	}

	payload := models.GetNotificationsUsageRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		ItemID:               itemID,
	}

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return err
	}

	var response models.GetNotificationsUsageResponse
	err = xml.Unmarshal(resp, &response)

	if err != nil {
		return err
	}

	if response.Ack != "Success" {
		return errors.New(response.Errors.LongMessage)
	}

	fmt.Println("DeliveredCount: ", response.NotificationStatistics.DeliveredCount)
	fmt.Println("QueuedNewCount: ", response.NotificationStatistics.QueuedNewCount)
	fmt.Println("QueuedPendingCount: ", response.NotificationStatistics.QueuedPendingCount)
	fmt.Println("ExpiredCount: ", response.NotificationStatistics.ExpiredCount)
	fmt.Println("ErrorCount: ", response.NotificationStatistics.ErrorCount)

	return nil
}
