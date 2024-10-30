package freeagent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/internal/freeagent/access-tokens"
)

type GetOpts struct {
	Path            string
	QueryParameters map[string]string
}

const BaseUrl = "https://api.freeagent.com/v2/"

func FreeagentApiGet[T any](ctx context.Context, opts GetOpts, target *T) error {
	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return err
	}

	newURL, err := opts.buildUrl()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", newURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+access_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}

func (o GetOpts) buildUrl() (string, error) {
	values := url.Values{}
	if o.QueryParameters != nil {
		for key, value := range o.QueryParameters {
			values.Add(key, value)
		}
	}

	return url.JoinPath(BaseUrl, o.Path, values.Encode())
}
