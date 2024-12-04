package ebay

import (
	"context"
	"errors"
	"fmt"

	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/sell/inventory"
)

const merchantLocationKey = "CONTROLALTREPEAT501A"

var controlAltRepeatWorkshop = inventory.InventoryLocation{
	Location: inventory.Location{
		Address: inventory.Address{
			AddressLine1: "White Cross Business Park",
			AddressLine2: "South Road",
			City:         "Lancaster",
			Country:      "GB",
			PostalCode:   "LA14XQ",
		},
		GeoCoordinates: inventory.GeoCoordinates{
			Latitude:  "54.044994",
			Longitude: "-2.797446",
		},
	},
	Name:                   "Control Alt Repeat",
	LocationTypes:          []string{"WAREHOUSE"},
	LocationWebURL:         "https://controlaltrepeat.net",
	MerchantLocationStatus: "ENABLED",
	Phone:                  "07984599771",
	TimeZoneID:             "Europe/London",
	OperatingHours: []inventory.OperatingHours{
		{
			DayOfWeekEnum: "MONDAY",
			Intervals: []inventory.Intervals{
				{
					Open:  "11:00:00",
					Close: "15:00:00",
				},
			},
		},
		{
			DayOfWeekEnum: "WEDNESDAY",
			Intervals: []inventory.Intervals{
				{
					Open:  "11:00:00",
					Close: "15:00:00",
				},
			},
		},
		{
			DayOfWeekEnum: "THURSDAY",
			Intervals: []inventory.Intervals{
				{
					Open:  "11:00:00",
					Close: "15:00:00",
				},
			},
		},
		{
			DayOfWeekEnum: "FRIDAY",
			Intervals: []inventory.Intervals{
				{
					Open:  "11:00:00",
					Close: "15:00:00",
				},
			},
		},
	},
}

func SetControlAltRepeatWorkshopLocation(ctx context.Context) error {
	response, err := inventory.GetLocations(ctx, inventory.GetLocationsOptions{Limit: 2})
	if err != nil {
		return err
	}

	fmt.Println(len(response.Locations))

	if len(response.Locations) > 1 {
		return errors.New("there should only be 1 or no locations")
	}

	if len(response.Locations) == 0 {
		return createWorkshop(ctx)
	}

	fmt.Println("Updating location")

	return updateWorkshop(ctx)
}

func createWorkshop(ctx context.Context) error {
	return CreateLocation(ctx, inventory.CreateLocationOptions{
		MerchantLocationKey: merchantLocationKey,
		InventoryLocation:   controlAltRepeatWorkshop,
	})
}

func updateWorkshop(ctx context.Context) error {
	return UpdateLocation(ctx, inventory.UpdateLocationOptions{
		MerchantLocationKey: merchantLocationKey,
		InventoryLocation:   controlAltRepeatWorkshop,
	})
}
