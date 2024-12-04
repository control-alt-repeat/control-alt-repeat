package inventory

import (
	"context"
	"fmt"
	"strconv"
)

var defaultlimit = 10

func GetLocations(ctx context.Context, opts GetLocationsOptions) (GetLocationsResponse, error) {
	limit := opts.Limit
	if limit == 0 {
		limit = defaultlimit
	}
	apiopts := requestOptions{
		Path: "/location",
		QueryParameters: map[string]string{
			"limit": strconv.Itoa(limit),
		},
	}
	var response GetLocationsResponse
	if err := apiGet(ctx, apiopts, &response); err != nil {
		return GetLocationsResponse{}, err
	}

	return response, nil
}

type GetLocationsOptions struct {
	Limit int
}

type GetLocationsResponse struct {
	Href      string     `json:"href"`
	Limit     string     `json:"limit"`
	Next      string     `json:"next"`
	Offset    string     `json:"offset"`
	Prev      string     `json:"prev"`
	Total     int        `json:"total"`
	Locations []Location `json:"locations"`
}

type CreateLocationOptions struct {
	InventoryLocation   InventoryLocation
	MerchantLocationKey string
}

func CreateLocation(ctx context.Context, opts CreateLocationOptions) error {
	apiopts := requestOptions{
		Path: fmt.Sprintf("/location/\"%s\"", opts.MerchantLocationKey),
	}
	return apiPost(ctx, apiopts, opts.InventoryLocation)
}

type UpdateLocationOptions struct {
	InventoryLocation   InventoryLocation
	MerchantLocationKey string
}

func UpdateLocation(ctx context.Context, opts UpdateLocationOptions) error {
	apiopts := requestOptions{
		Path: fmt.Sprintf("/location/\"%s\"/update_location_details", opts.MerchantLocationKey),
	}
	return apiPost(ctx, apiopts, opts.InventoryLocation)
}
