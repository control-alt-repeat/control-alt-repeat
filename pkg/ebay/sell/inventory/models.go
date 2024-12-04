package inventory

type InventoryLocation struct {
	Location                        Location                        `json:"location"`
	LocationAdditionalInformation   string                          `json:"locationAdditionalInformation,omitempty"`
	LocationInstructions            string                          `json:"locationInstructions,omitempty"`
	LocationTypes                   []string                        `json:"locationTypes"`
	LocationWebURL                  string                          `json:"locationWebUrl,omitempty"`
	MerchantLocationStatus          string                          `json:"merchantLocationStatus"`
	Name                            string                          `json:"name"`
	OperatingHours                  []OperatingHours                `json:"operatingHours,omitempty"`
	Phone                           string                          `json:"phone,omitempty"`
	SpecialHours                    []SpecialHours                  `json:"specialHours,omitempty"`
	TimeZoneID                      string                          `json:"timeZoneId,omitempty"`
	FulfillmentCenterSpecifications FulfillmentCenterSpecifications `json:"fulfillmentCenterSpecifications,omitempty"`
}

type Address struct {
	AddressLine1    string `json:"addressLine1,omitempty"`
	AddressLine2    string `json:"addressLine2,omitempty"`
	City            string `json:"city,omitempty"`
	Country         string `json:"country,omitempty"`
	County          string `json:"county,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	StateOrProvince string `json:"stateOrProvince,omitempty"`
}

type GeoCoordinates struct {
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
}

type Location struct {
	Address        Address        `json:"address"`
	GeoCoordinates GeoCoordinates `json:"geoCoordinates,omitempty"`
	LocationID     string         `json:"locationId,omitempty"`
}

type Intervals struct {
	Close string `json:"close"`
	Open  string `json:"open"`
}

type OperatingHours struct {
	DayOfWeekEnum string      `json:"dayOfWeekEnum"`
	Intervals     []Intervals `json:"intervals"`
}

type SpecialHours struct {
	Date      string      `json:"date"`
	Intervals []Intervals `json:"intervals"`
}

type Overrides struct {
	CutOffTime string `json:"cutOffTime"`
	EndDate    string `json:"endDate"`
	StartDate  string `json:"startDate"`
}

type WeeklySchedule struct {
	CutOffTime    string   `json:"cutOffTime"`
	DayOfWeekEnum []string `json:"dayOfWeekEnum"`
}

type SameDayShippingCutOffTimes struct {
	Overrides      []Overrides      `json:"overrides,omitempty"`
	WeeklySchedule []WeeklySchedule `json:"weeklySchedule,omitempty"`
}

type FulfillmentCenterSpecifications struct {
	SameDayShippingCutOffTimes SameDayShippingCutOffTimes `json:"sameDayShippingCutOffTimes,omitempty"`
}
