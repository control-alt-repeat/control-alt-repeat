package ebay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

func GetNotificationUsage(itemID string) error {

	fmt.Println("Getting notification usage")

	request, requesterCredentials, err := newTraditionalAPIRequest("GetNotificationsUsage")
	if err != nil {
		return err
	}

	payload := models.GetNotificationsUsageRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		ItemID:               itemID,
	}

	resp, err := request.Post(payload)
	if err != nil {
		return err
	}

	var response models.GetNotificationsUsageResponse
	err = xml.Unmarshal(resp, &response)

	fmt.Println(string(resp))

	if err != nil {
		return err
	}

	fmt.Println(response.Ack)

	if response.Ack != "Success" {
		err = errors.New(response.Errors.LongMessage)
	}

	fmt.Println("DeliveredCount: ", response.NotificationStatistics.DeliveredCount)
	fmt.Println("QueuedNewCount: ", response.NotificationStatistics.QueuedNewCount)
	fmt.Println("QueuedPendingCount: ", response.NotificationStatistics.QueuedPendingCount)
	fmt.Println("ExpiredCount: ", response.NotificationStatistics.ExpiredCount)
	fmt.Println("ErrorCount: ", response.NotificationStatistics.ErrorCount)

	return err
}
