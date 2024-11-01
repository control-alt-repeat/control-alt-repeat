package freeagent

import (
	"context"
)

var defaultSort = "name"

func ListContacts(ctx context.Context, opts ListContactsOptions) ([]Contact, error) {
	sort := opts.Sort
	if sort == "" {
		sort = defaultSort
	}
	apiopts := ApiGetOptions{
		Path: "contacts",
		QueryParameters: map[string]string{
			"sort": sort,
		},
	}
	var response ListContactsResponse
	if err := FreeagentApiGet(ctx, apiopts, &response); err != nil {
		return []Contact{}, err
	}

	return response.Contacts, nil
}

type ListContactsOptions struct {
	Sort string
}

type ListContactsResponse struct {
	Contacts []Contact `json:"contacts"`
}
