// https://developer.ebay.com/api-docs/static/platform-notifications-landing.html
// https://developer.ebay.com/Devzone/XML/docs/Reference/eBay/SetNotificationPreferences.html#Request.UserDeliveryPreferenceArray

package ebay

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/control-alt-repeat/control-alt-repeat/internal/aws"
)

var applicationDeliveryPreferences = ApplicationDeliveryPreferences{
	ApplicationEnable: "Enable",
	DeviceType:        "Platform",
}

var userDeliveryPreferenceArray = UserDeliveryPreferenceArray{
	NotificationEnable: []NotificationEnable{
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
	request, requesterCredentials, err := newTraditionalAPIRequest("SetNotificationPreferences")
	if err != nil {
		return err
	}

	ebayNotificationEndpoint, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/notification/endpoint")
	if err != nil {
		return err
	}

	applicationDeliveryPreferences.ApplicationURL = ebayNotificationEndpoint

	payload := SetNotificationPreferencesRequest{
		Xmlns:                          "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials:           *requesterCredentials,
		ApplicationDeliveryPreferences: applicationDeliveryPreferences,
		UserDeliveryPreferenceArray:    userDeliveryPreferenceArray,
		WarningLevel:                   "High",
	}

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return err
	}

	var response SetNotificationPreferencesResponse
	err = xml.Unmarshal(resp, &response)

	fmt.Println(string(resp))

	if err != nil {
		return err
	}

	if response.Ack != "Success" {
		return errors.New(response.Errors.LongMessage)
	}

	return nil
}

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
