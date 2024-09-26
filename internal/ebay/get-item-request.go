package ebay

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func GetItem(itemId string) error {
	// ebayAccessToken, err := getAccessToken()

	// if err != nil {
	// 	fmt.Println("Error getting eBay access token:", err)
	// 	return err
	// }

	payload := GetItemRequest{
		Xmlns: "urn:ebay:apis:eBLBaseComponents",
		// RequesterCredentials: RequesterCredentials{
		// 	EBayAuthToken: ebayAccessToken,
		// },
		ItemID: itemId,
	}

	request, err := newTraditionalAPIRequest("GetItem", payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = payload.RequesterCredentials.SetEBayAuthToken()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// xmlData, err := xml.MarshalIndent(payload, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error marshalling XML:", err)
	// 	return err
	// }

	// xmlData = append([]byte(xml.Header), xmlData...)

	// xmlString := string(xmlData)
	// fmt.Println("Marshalled XML:")
	// fmt.Println(xmlString)

	// reader := strings.NewReader(xmlString)

	// req, err := http.NewRequest("POST", "https://api.ebay.com/ws/api.dll", reader)

	client := &http.Client{}

	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := client.Do(&request.HttpRequest)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var getItemResponse GetItemResponse
	err = xml.Unmarshal(body, &getItemResponse)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// Models generated using https://blog.kowalczyk.info/tools/xmltogo/
type GetItemRequest struct {
	XMLName              xml.Name             `xml:"GetItemRequest"`
	Xmlns                string               `xml:"xmlns,attr"`
	RequesterCredentials RequesterCredentials `xml:"RequesterCredentials"`
	ItemID               string               `xml:"ItemID"`
}

type GetItemResponse struct {
	XMLName               xml.Name `xml:"GetItemResponse"`
	Text                  string   `xml:",chardata"`
	Xmlns                 string   `xml:"xmlns,attr"`
	Timestamp             string   `xml:"Timestamp"`
	Ack                   string   `xml:"Ack"`
	Version               string   `xml:"Version"`
	Build                 string   `xml:"Build"`
	HardExpirationWarning string   `xml:"HardExpirationWarning"`
	Item                  struct {
		Text            string `xml:",chardata"`
		AutoPay         string `xml:"AutoPay"`
		BuyerProtection string `xml:"BuyerProtection"`
		BuyItNowPrice   struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"BuyItNowPrice"`
		Country        string `xml:"Country"`
		Currency       string `xml:"Currency"`
		ItemID         string `xml:"ItemID"`
		ListingDetails struct {
			Text                   string `xml:",chardata"`
			Adult                  string `xml:"Adult"`
			BindingAuction         string `xml:"BindingAuction"`
			CheckoutEnabled        string `xml:"CheckoutEnabled"`
			ConvertedBuyItNowPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedBuyItNowPrice"`
			ConvertedStartPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedStartPrice"`
			ConvertedReservePrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedReservePrice"`
			HasReservePrice             string `xml:"HasReservePrice"`
			RelistedItemID              string `xml:"RelistedItemID"`
			StartTime                   string `xml:"StartTime"`
			EndTime                     string `xml:"EndTime"`
			ViewItemURL                 string `xml:"ViewItemURL"`
			HasUnansweredQuestions      string `xml:"HasUnansweredQuestions"`
			HasPublicMessages           string `xml:"HasPublicMessages"`
			ViewItemURLForNaturalSearch string `xml:"ViewItemURLForNaturalSearch"`
		} `xml:"ListingDetails"`
		ListingDuration string `xml:"ListingDuration"`
		ListingType     string `xml:"ListingType"`
		Location        string `xml:"Location"`
		PrimaryCategory struct {
			Text         string `xml:",chardata"`
			CategoryID   string `xml:"CategoryID"`
			CategoryName string `xml:"CategoryName"`
		} `xml:"PrimaryCategory"`
		PrivateListing    string `xml:"PrivateListing"`
		Quantity          string `xml:"Quantity"`
		IsItemEMSEligible string `xml:"IsItemEMSEligible"`
		ReservePrice      struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"ReservePrice"`
		ReviseStatus struct {
			Text        string `xml:",chardata"`
			ItemRevised string `xml:"ItemRevised"`
		} `xml:"ReviseStatus"`
		Seller struct {
			Text                    string `xml:",chardata"`
			AboutMePage             string `xml:"AboutMePage"`
			Email                   string `xml:"Email"`
			FeedbackScore           string `xml:"FeedbackScore"`
			PositiveFeedbackPercent string `xml:"PositiveFeedbackPercent"`
			FeedbackPrivate         string `xml:"FeedbackPrivate"`
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
				Text                  string `xml:",chardata"`
				AllowPaymentEdit      string `xml:"AllowPaymentEdit"`
				CheckoutEnabled       string `xml:"CheckoutEnabled"`
				CIPBankAccountStored  string `xml:"CIPBankAccountStored"`
				GoodStanding          string `xml:"GoodStanding"`
				LiveAuctionAuthorized string `xml:"LiveAuctionAuthorized"`
				MerchandizingPref     string `xml:"MerchandizingPref"`
				QualifiesForB2BVAT    string `xml:"QualifiesForB2BVAT"`
				StoreOwner            string `xml:"StoreOwner"`
				StoreURL              string `xml:"StoreURL"`
				SellerBusinessType    string `xml:"SellerBusinessType"`
				SafePaymentExempt     string `xml:"SafePaymentExempt"`
			} `xml:"SellerInfo"`
			MotorsDealer string `xml:"MotorsDealer"`
		} `xml:"Seller"`
		SellingStatus struct {
			Text         string `xml:",chardata"`
			BidCount     string `xml:"BidCount"`
			BidIncrement struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"BidIncrement"`
			ConvertedCurrentPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"ConvertedCurrentPrice"`
			CurrentPrice struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"CurrentPrice"`
			HighBidder struct {
				Text                    string `xml:",chardata"`
				AboutMePage             string `xml:"AboutMePage"`
				EIASToken               string `xml:"EIASToken"`
				Email                   string `xml:"Email"`
				FeedbackScore           string `xml:"FeedbackScore"`
				PositiveFeedbackPercent string `xml:"PositiveFeedbackPercent"`
				EBayGoodStanding        string `xml:"eBayGoodStanding"`
				NewUser                 string `xml:"NewUser"`
				RegistrationDate        string `xml:"RegistrationDate"`
				Site                    string `xml:"Site"`
				UserID                  string `xml:"UserID"`
				VATStatus               string `xml:"VATStatus"`
				UserAnonymized          string `xml:"UserAnonymized"`
			} `xml:"HighBidder"`
			LeadCount    string `xml:"LeadCount"`
			MinimumToBid struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"MinimumToBid"`
			QuantitySold                string `xml:"QuantitySold"`
			ReserveMet                  string `xml:"ReserveMet"`
			SecondChanceEligible        string `xml:"SecondChanceEligible"`
			ListingStatus               string `xml:"ListingStatus"`
			QuantitySoldByPickupInStore string `xml:"QuantitySoldByPickupInStore"`
		} `xml:"SellingStatus"`
		ShippingDetails struct {
			Text                   string `xml:",chardata"`
			ApplyShippingDiscount  string `xml:"ApplyShippingDiscount"`
			CalculatedShippingRate struct {
				Text        string `xml:",chardata"`
				WeightMajor struct {
					Text              string `xml:",chardata"`
					MeasurementSystem string `xml:"measurementSystem,attr"`
					Unit              string `xml:"unit,attr"`
				} `xml:"WeightMajor"`
				WeightMinor struct {
					Text              string `xml:",chardata"`
					MeasurementSystem string `xml:"measurementSystem,attr"`
					Unit              string `xml:"unit,attr"`
				} `xml:"WeightMinor"`
			} `xml:"CalculatedShippingRate"`
			SalesTax struct {
				Text                  string `xml:",chardata"`
				SalesTaxPercent       string `xml:"SalesTaxPercent"`
				ShippingIncludedInTax string `xml:"ShippingIncludedInTax"`
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
				FreeShipping            string `xml:"FreeShipping"`
			} `xml:"ShippingServiceOptions"`
			ShippingType                           string `xml:"ShippingType"`
			ThirdPartyCheckout                     string `xml:"ThirdPartyCheckout"`
			ShippingDiscountProfileID              string `xml:"ShippingDiscountProfileID"`
			InternationalShippingDiscountProfileID string `xml:"InternationalShippingDiscountProfileID"`
			SellerExcludeShipToLocationsPreference string `xml:"SellerExcludeShipToLocationsPreference"`
		} `xml:"ShippingDetails"`
		ShipToLocations string `xml:"ShipToLocations"`
		Site            string `xml:"Site"`
		StartPrice      struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"StartPrice"`
		Storefront struct {
			Text             string `xml:",chardata"`
			StoreCategoryID  string `xml:"StoreCategoryID"`
			StoreCategory2ID string `xml:"StoreCategory2ID"`
			StoreURL         string `xml:"StoreURL"`
		} `xml:"Storefront"`
		TimeLeft         string `xml:"TimeLeft"`
		Title            string `xml:"Title"`
		BestOfferDetails struct {
			Text             string `xml:",chardata"`
			BestOfferCount   string `xml:"BestOfferCount"`
			BestOfferEnabled string `xml:"BestOfferEnabled"`
			NewBestOffer     string `xml:"NewBestOffer"`
		} `xml:"BestOfferDetails"`
		GetItFast      string `xml:"GetItFast"`
		SKU            string `xml:"SKU"`
		PostalCode     string `xml:"PostalCode"`
		PictureDetails struct {
			Text        string   `xml:",chardata"`
			GalleryType string   `xml:"GalleryType"`
			PictureURL  []string `xml:"PictureURL"`
		} `xml:"PictureDetails"`
		DispatchTimeMax       string `xml:"DispatchTimeMax"`
		ProxyItem             string `xml:"ProxyItem"`
		BusinessSellerDetails struct {
			Text    string `xml:",chardata"`
			Address struct {
				Text            string `xml:",chardata"`
				Street1         string `xml:"Street1"`
				Street2         string `xml:"Street2"`
				CityName        string `xml:"CityName"`
				StateOrProvince string `xml:"StateOrProvince"`
				CountryName     string `xml:"CountryName"`
				Phone           string `xml:"Phone"`
				PostalCode      string `xml:"PostalCode"`
				CompanyName     string `xml:"CompanyName"`
				FirstName       string `xml:"FirstName"`
				LastName        string `xml:"LastName"`
			} `xml:"Address"`
		} `xml:"BusinessSellerDetails"`
		BuyerGuaranteePrice struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"BuyerGuaranteePrice"`
		ReturnPolicy struct {
			Text                                  string `xml:",chardata"`
			ReturnsWithinOption                   string `xml:"ReturnsWithinOption"`
			ReturnsWithin                         string `xml:"ReturnsWithin"`
			ReturnsAcceptedOption                 string `xml:"ReturnsAcceptedOption"`
			ReturnsAccepted                       string `xml:"ReturnsAccepted"`
			ShippingCostPaidByOption              string `xml:"ShippingCostPaidByOption"`
			ShippingCostPaidBy                    string `xml:"ShippingCostPaidBy"`
			InternationalReturnsAcceptedOption    string `xml:"InternationalReturnsAcceptedOption"`
			InternationalReturnsWithinOption      string `xml:"InternationalReturnsWithinOption"`
			InternationalShippingCostPaidByOption string `xml:"InternationalShippingCostPaidByOption"`
		} `xml:"ReturnPolicy"`
		ConditionID                   string `xml:"ConditionID"`
		ConditionDisplayName          string `xml:"ConditionDisplayName"`
		PostCheckoutExperienceEnabled string `xml:"PostCheckoutExperienceEnabled"`
		ShippingPackageDetails        struct {
			Text              string `xml:",chardata"`
			ShippingIrregular string `xml:"ShippingIrregular"`
			ShippingPackage   string `xml:"ShippingPackage"`
			WeightMajor       struct {
				Text              string `xml:",chardata"`
				MeasurementSystem string `xml:"measurementSystem,attr"`
				Unit              string `xml:"unit,attr"`
			} `xml:"WeightMajor"`
			WeightMinor struct {
				Text              string `xml:",chardata"`
				MeasurementSystem string `xml:"measurementSystem,attr"`
				Unit              string `xml:"unit,attr"`
			} `xml:"WeightMinor"`
		} `xml:"ShippingPackageDetails"`
		HideFromSearch       string `xml:"HideFromSearch"`
		PickupInStoreDetails struct {
			Text                      string `xml:",chardata"`
			AvailableForPickupInStore string `xml:"AvailableForPickupInStore"`
		} `xml:"PickupInStoreDetails"`
		EBayPlus            string `xml:"eBayPlus"`
		EBayPlusEligible    string `xml:"eBayPlusEligible"`
		IsSecureDescription string `xml:"IsSecureDescription"`
	} `xml:"Item"`
}
