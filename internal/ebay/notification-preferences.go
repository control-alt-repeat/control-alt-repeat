// https://developer.ebay.com/api-docs/static/platform-notifications-landing.html
// https://developer.ebay.com/Devzone/XML/docs/Reference/eBay/SetNotificationPreferences.html#Request.UserDeliveryPreferenceArray

package ebay

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/models"
)

var applicationDeliveryPreferences = models.ApplicationDeliveryPreferences{
	ApplicationEnable: "Enable",
	DeviceType:        "Platform",
}

var userDeliveryPreferenceArray = models.UserDeliveryPreferenceArray{
	NotificationEnable: []models.NotificationEnable{
		// https://developer.ebay.com/Devzone/XML/docs/Reference/eBay/types/NotificationEventTypeCodeType.html
		// {EventEnable: "Enable", EventType: "AccountSummary"},
		// {EventEnable: "Enable", EventType: "AccountSuspended"},
		// {EventEnable: "Enable", EventType: "AskSellerQuestion"},
		// {EventEnable: "Enable", EventType: "AuctionCheckoutComplete"},
		// {EventEnable: "Enable", EventType: "BestOffer"},
		// {EventEnable: "Enable", EventType: "BestOfferDeclined"},
		// {EventEnable: "Enable", EventType: "BestOfferPlaced"},
		// {EventEnable: "Enable", EventType: "BidItemEndingSoon"},
		// {EventEnable: "Enable", EventType: "BidPlaced"},
		// {EventEnable: "Enable", EventType: "BidReceived"},
		// {EventEnable: "Enable", EventType: "BulkDataExchangeJobCompleted"},
		// {EventEnable: "Enable", EventType: "BuyerCancelRequested"},
		// {EventEnable: "Enable", EventType: "BuyerNoShow"},
		// {EventEnable: "Enable", EventType: "CheckoutBuyerRequestsTotal"},
		// {EventEnable: "Enable", EventType: "CounterOfferReceived"},
		// {EventEnable: "Enable", EventType: "CustomCode"},
		// {EventEnable: "Enable", EventType: "EBNOrderCanceled"},
		// {EventEnable: "Enable", EventType: "EBNOrderPickedUp"},
		// {EventEnable: "Enable", EventType: "EBPAppealedCase"},
		// {EventEnable: "Enable", EventType: "EBPClosedAppeal"},
		// {EventEnable: "Enable", EventType: "EBPClosedCase"},
		// {EventEnable: "Enable", EventType: "EBPEscalatedCase"},
		// {EventEnable: "Enable", EventType: "EBPMyPaymentDue"},
		// {EventEnable: "Enable", EventType: "EBPMyResponseDue"},
		// {EventEnable: "Enable", EventType: "EBPOnHoldCase"},
		// {EventEnable: "Enable", EventType: "EBPOtherPartyResponseDue"},
		// {EventEnable: "Enable", EventType: "EBPPaymentDone"},
		// {EventEnable: "Enable", EventType: "EmailAddressChanged"},
		// {EventEnable: "Enable", EventType: "EndOfAuction"},
		// {EventEnable: "Enable", EventType: "Feedback"},
		// {EventEnable: "Enable", EventType: "FeedbackLeft"},
		// {EventEnable: "Enable", EventType: "FeedbackReceived"},
		// {EventEnable: "Enable", EventType: "FeedbackStarChanged"},
		// {EventEnable: "Enable", EventType: "FixedPriceTransaction"},
		// {EventEnable: "Enable", EventType: "INRBuyerRespondedToDispute"},
		{EventEnable: "Enable", EventType: "ItemAddedToWatchList"},
		{EventEnable: "Enable", EventType: "ItemClosed"},
		{EventEnable: "Enable", EventType: "ItemExtended"},
		{EventEnable: "Enable", EventType: "ItemListed"},
		{EventEnable: "Enable", EventType: "ItemLost"},
		{EventEnable: "Enable", EventType: "ItemMarkedPaid"},
		{EventEnable: "Enable", EventType: "ItemMarkedShipped"},
		{EventEnable: "Enable", EventType: "ItemOutOfStock"},
		{EventEnable: "Enable", EventType: "ItemReadyForPickup"},
		{EventEnable: "Enable", EventType: "ItemRemovedFromWatchList"},
		{EventEnable: "Enable", EventType: "ItemRevised"},
		{EventEnable: "Enable", EventType: "ItemRevisedAddCharity"},
		{EventEnable: "Enable", EventType: "ItemSold"},
		{EventEnable: "Enable", EventType: "ItemSuspended"},
		{EventEnable: "Enable", EventType: "ItemUnsold"},
		{EventEnable: "Enable", EventType: "ItemWon"},
		// {EventEnable: "Enable", EventType: "M2MMessageStatusChange"},
		// {EventEnable: "Enable", EventType: "MyMessageseBayMessage"},
		// {EventEnable: "Enable", EventType: "MyMessagesHighPriorityMessage"},
		// {EventEnable: "Enable", EventType: "MyMessagesM2MMessage"},
		// {EventEnable: "Enable", EventType: "None"},
		// {EventEnable: "Enable", EventType: "OrderInquiryReminderForEscalation"},
		// {EventEnable: "Enable", EventType: "OutBid"},
		// {EventEnable: "Enable", EventType: "PasswordChanged"},
		// {EventEnable: "Enable", EventType: "PasswordHintChanged"},
		// {EventEnable: "Enable", EventType: "PaymentDetailChanged"},
		// {EventEnable: "Enable", EventType: "PaymentReminder"},
		// {EventEnable: "Enable", EventType: "ReturnClosed"},
		// {EventEnable: "Enable", EventType: "ReturnCreated"},
		// {EventEnable: "Enable", EventType: "ReturnDelivered"},
		// {EventEnable: "Enable", EventType: "ReturnEscalated"},
		// {EventEnable: "Enable", EventType: "ReturnRefundOverdue"},
		// {EventEnable: "Enable", EventType: "ReturnSellerInfoOverdue"},
		// {EventEnable: "Enable", EventType: "ReturnShipped"},
		// {EventEnable: "Enable", EventType: "ReturnWaitingForSellerInfo"},
		// {EventEnable: "Enable", EventType: "SecondChanceOffer"},
		// {EventEnable: "Enable", EventType: "ShoppingCartItemEndingSoon"},
		// {EventEnable: "Enable", EventType: "TokenRevocation"},
		// {EventEnable: "Enable", EventType: "WatchedItemEndingSoon"},
		// {EventEnable: "Enable", EventType: "WebnextMobilePhotoSync"},
	},
}

func SetNotificationPreferences(ctx context.Context) error {
	fmt.Println("Setting notification preferences")

	request, requesterCredentials, err := newTraditionalAPIRequest("SetNotificationPreferences")
	if err != nil {
		return err
	}

	ebayNotificationEndpoint, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/notification/endpoint")
	if err != nil {
		return err
	}
	fmt.Println(ebayNotificationEndpoint)

	applicationDeliveryPreferences.ApplicationURL = ebayNotificationEndpoint

	payload := models.SetNotificationPreferencesRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		// DeliveryURLName:                "AWS Lambda function - *-ebay-notifification-endpoint",
		ApplicationDeliveryPreferences: applicationDeliveryPreferences,
		UserDeliveryPreferenceArray:    userDeliveryPreferenceArray,
		WarningLevel:                   "High",
	}

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return err
	}

	var response models.SetNotificationPreferencesResponse
	err = xml.Unmarshal(resp, &response)

	fmt.Println(string(resp))

	if err != nil {
		return err
	}

	if response.Ack != "Success" {
		err = errors.New(response.Errors.LongMessage)
	}

	fmt.Println(response.Ack)

	fmt.Println()

	return err
}
