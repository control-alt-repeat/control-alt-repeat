package traditionalapi

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
)

func GetItemTransactions(ctx context.Context, ebayListingID string) (*GetItemTransactionsResponse, error) {
	request, requesterCredentials, err := newTraditionalAPIRequest("GetItemTransactions")
	if err != nil {
		fmt.Println(err)
		return &GetItemTransactionsResponse{}, err
	}

	payload := GetItemTransactionsRequest{
		Xmlns:                "urn:ebay:apis:eBLBaseComponents",
		RequesterCredentials: *requesterCredentials,
		ItemID:               ebayListingID,
		IncludeFinalValueFee: true,
	}

	resp, err := request.Post(ctx, payload)
	if err != nil {
		return &GetItemTransactionsResponse{}, err
	}

	var getItemTransactionsResponse GetItemTransactionsResponse
	if err = xml.Unmarshal(resp, &getItemTransactionsResponse); err != nil {
		return &GetItemTransactionsResponse{}, err
	}

	if getItemTransactionsResponse.Ack != "Success" {
		err = errors.New(getItemTransactionsResponse.Errors.LongMessage)
	}

	return &getItemTransactionsResponse, err
}

type GetItemTransactionsRequest struct {
	XMLName              xml.Name             `xml:"GetItemTransactionsRequest"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	ItemID               string               `xml:"ItemID"`
	IncludeFinalValueFee bool                 `xml:"IncludeFinalValueFee"`
}

// GetItemTransactionsResponse was generated 2024-12-04 16:50:04 by https://xml-to-go.github.io/ in Ukraine.
type GetItemTransactionsResponse struct {
	XMLName   xml.Name `xml:"GetItemTransactionsResponse"`
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
	PaginationResult struct {
		Text                 string `xml:",chardata"`
		TotalNumberOfPages   string `xml:"TotalNumberOfPages"`
		TotalNumberOfEntries string `xml:"TotalNumberOfEntries"`
	} `xml:"PaginationResult"`
	HasMoreTransactions            string `xml:"HasMoreTransactions"`
	TransactionsPerPage            string `xml:"TransactionsPerPage"`
	PageNumber                     string `xml:"PageNumber"`
	ReturnedTransactionCountActual string `xml:"ReturnedTransactionCountActual"`
	Item                           struct {
		Text           string `xml:",chardata"`
		AutoPay        string `xml:"AutoPay"`
		Currency       string `xml:"Currency"`
		ItemID         string `xml:"ItemID"`
		ListingDetails struct {
			Text                        string `xml:",chardata"`
			ViewItemURL                 string `xml:"ViewItemURL"`
			ViewItemURLForNaturalSearch string `xml:"ViewItemURLForNaturalSearch"`
		} `xml:"ListingDetails"`
		ListingType    string `xml:"ListingType"`
		PaymentMethods string `xml:"PaymentMethods"`
		PrivateListing string `xml:"PrivateListing"`
		Quantity       string `xml:"Quantity"`
		Seller         struct {
			Text                    string `xml:",chardata"`
			AboutMePage             string `xml:"AboutMePage"`
			EIASToken               string `xml:"EIASToken"`
			Email                   string `xml:"Email"`
			FeedbackScore           string `xml:"FeedbackScore"`
			PositiveFeedbackPercent string `xml:"PositiveFeedbackPercent"`
			FeedbackPrivate         string `xml:"FeedbackPrivate"`
			FeedbackRatingStar      string `xml:"FeedbackRatingStar"`
			IDVerified              string `xml:"IDVerified"`
			EBayGoodStanding        string `xml:"eBayGoodStanding"`
			NewUser                 string `xml:"NewUser"`
			RegistrationDate        string `xml:"RegistrationDate"`
			Site                    string `xml:"Site"`
			Status                  string `xml:"Status"`
			UserID                  string `xml:"UserID"`
			UserIDChanged           string `xml:"UserIDChanged"`
			UserIDLastChanged       string `xml:"UserIDLastChanged"`
			VATStatus               string `xml:"VATStatus"`
			SellerInfo              struct {
				Text                 string `xml:",chardata"`
				AllowPaymentEdit     string `xml:"AllowPaymentEdit"`
				CheckoutEnabled      string `xml:"CheckoutEnabled"`
				CIPBankAccountStored string `xml:"CIPBankAccountStored"`
				GoodStanding         string `xml:"GoodStanding"`
				QualifiesForB2BVAT   string `xml:"QualifiesForB2BVAT"`
				StoreOwner           string `xml:"StoreOwner"`
				StoreURL             string `xml:"StoreURL"`
				SafePaymentExempt    string `xml:"SafePaymentExempt"`
			} `xml:"SellerInfo"`
		} `xml:"Seller"`
		SellingStatus struct {
			Text                  string `xml:",chardata"`
			BidCount              string `xml:"BidCount"`
			ConvertedCurrentPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedCurrentPrice"`
			CurrentPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"CurrentPrice"`
			FinalValueFee struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"FinalValueFee"`
			QuantitySold  string `xml:"QuantitySold"`
			ListingStatus string `xml:"ListingStatus"`
		} `xml:"SellingStatus"`
		Site                                string `xml:"Site"`
		Title                               string `xml:"Title"`
		GetItFast                           string `xml:"GetItFast"`
		IntegratedMerchantCreditCardEnabled string `xml:"IntegratedMerchantCreditCardEnabled"`
	} `xml:"Item"`
	TransactionArray struct {
		Text        string `xml:",chardata"`
		Transaction struct {
			Text       string `xml:",chardata"`
			AmountPaid struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"AmountPaid"`
			AdjustmentAmount struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"AdjustmentAmount"`
			ConvertedAdjustmentAmount struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedAdjustmentAmount"`
			Buyer struct {
				Text                    string `xml:",chardata"`
				AboutMePage             string `xml:"AboutMePage"`
				EIASToken               string `xml:"EIASToken"`
				Email                   string `xml:"Email"`
				FeedbackScore           string `xml:"FeedbackScore"`
				PositiveFeedbackPercent string `xml:"PositiveFeedbackPercent"`
				FeedbackPrivate         string `xml:"FeedbackPrivate"`
				FeedbackRatingStar      string `xml:"FeedbackRatingStar"`
				IDVerified              string `xml:"IDVerified"`
				EBayGoodStanding        string `xml:"eBayGoodStanding"`
				NewUser                 string `xml:"NewUser"`
				RegistrationDate        string `xml:"RegistrationDate"`
				Site                    string `xml:"Site"`
				Status                  string `xml:"Status"`
				UserID                  string `xml:"UserID"`
				UserIDChanged           string `xml:"UserIDChanged"`
				UserIDLastChanged       string `xml:"UserIDLastChanged"`
				VATStatus               string `xml:"VATStatus"`
				BuyerInfo               struct {
					Text            string `xml:",chardata"`
					ShippingAddress struct {
						Text            string `xml:",chardata"`
						Street1         string `xml:"Street1"`
						CityName        string `xml:"CityName"`
						StateOrProvince string `xml:"StateOrProvince"`
						Country         string `xml:"Country"`
						CountryName     string `xml:"CountryName"`
						Phone           string `xml:"Phone"`
						PostalCode      string `xml:"PostalCode"`
						AddressOwner    string `xml:"AddressOwner"`
					} `xml:"ShippingAddress"`
				} `xml:"BuyerInfo"`
				UserAnonymized string `xml:"UserAnonymized"`
			} `xml:"Buyer"`
			ShippingDetails struct {
				Text                      string `xml:",chardata"`
				ChangePaymentInstructions string `xml:"ChangePaymentInstructions"`
				PaymentEdited             string `xml:"PaymentEdited"`
				SalesTax                  struct {
					Text                  string `xml:",chardata"`
					SalesTaxPercent       string `xml:"SalesTaxPercent"`
					ShippingIncludedInTax string `xml:"ShippingIncludedInTax"`
					SalesTaxAmount        struct {
						Text       string `xml:",chardata"`
						CurrencyID string `xml:"currencyID,attr"`
					} `xml:"SalesTaxAmount"`
				} `xml:"SalesTax"`
				ShippingServiceOptions struct {
					Text                string `xml:",chardata"`
					ShippingService     string `xml:"ShippingService"`
					ShippingServiceCost struct {
						Text       string `xml:",chardata"`
						CurrencyID string `xml:"currencyID,attr"`
					} `xml:"ShippingServiceCost"`
					ShippingServicePriority string `xml:"ShippingServicePriority"`
					ExpeditedService        string `xml:"ExpeditedService"`
					ShippingTimeMin         string `xml:"ShippingTimeMin"`
					ShippingTimeMax         string `xml:"ShippingTimeMax"`
				} `xml:"ShippingServiceOptions"`
				ShippingType                    string `xml:"ShippingType"`
				SellingManagerSalesRecordNumber string `xml:"SellingManagerSalesRecordNumber"`
				ThirdPartyCheckout              string `xml:"ThirdPartyCheckout"`
				TaxTable                        string `xml:"TaxTable"`
				GetItFast                       string `xml:"GetItFast"`
			} `xml:"ShippingDetails"`
			ConvertedAmountPaid struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedAmountPaid"`
			ConvertedTransactionPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedTransactionPrice"`
			CreatedDate       string `xml:"CreatedDate"`
			DepositType       string `xml:"DepositType"`
			QuantityPurchased string `xml:"QuantityPurchased"`
			Status            struct {
				Text                                string `xml:",chardata"`
				EBayPaymentStatus                   string `xml:"eBayPaymentStatus"`
				CheckoutStatus                      string `xml:"CheckoutStatus"`
				LastTimeModified                    string `xml:"LastTimeModified"`
				PaymentMethodUsed                   string `xml:"PaymentMethodUsed"`
				CompleteStatus                      string `xml:"CompleteStatus"`
				BuyerSelectedShipping               string `xml:"BuyerSelectedShipping"`
				PaymentHoldStatus                   string `xml:"PaymentHoldStatus"`
				IntegratedMerchantCreditCardEnabled string `xml:"IntegratedMerchantCreditCardEnabled"`
			} `xml:"Status"`
			TransactionID    string `xml:"TransactionID"`
			TransactionPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"TransactionPrice"`
			BestOfferSale           string `xml:"BestOfferSale"`
			ShippingServiceSelected struct {
				Text                string `xml:",chardata"`
				ShippingService     string `xml:"ShippingService"`
				ShippingServiceCost struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
				} `xml:"ShippingServiceCost"`
			} `xml:"ShippingServiceSelected"`
			PaidTime            string `xml:"PaidTime"`
			ShippedTime         string `xml:"ShippedTime"`
			TransactionSiteID   string `xml:"TransactionSiteID"`
			Platform            string `xml:"Platform"`
			PayPalEmailAddress  string `xml:"PayPalEmailAddress"`
			BuyerGuaranteePrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"BuyerGuaranteePrice"`
			IntangibleItem string `xml:"IntangibleItem"`
		} `xml:"Transaction"`
	} `xml:"TransactionArray"`
	PayPalPreferred string `xml:"PayPalPreferred"`
}
