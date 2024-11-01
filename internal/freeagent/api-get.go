package freeagent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/internal/freeagent/access-tokens"
)

type ApiGetOptions struct {
	Path            string
	QueryParameters map[string]string
}

const BaseUrl = "https://api.freeagent.com/v2/"

func FreeagentApiGet[T any](ctx context.Context, opts ApiGetOptions, target *T) error {
	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return err
	}

	fmt.Println(access_token[0:20])

	newURL, err := opts.buildUrl()
	if err != nil {
		return err
	}

	fmt.Println(newURL)

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

func (o ApiGetOptions) buildUrl() (string, error) {
	u, err := url.Parse(BaseUrl)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	if o.QueryParameters != nil {
		for key, value := range o.QueryParameters {
			values.Add(key, value)
		}
	}

	u.Path += o.Path
	u.RawQuery = values.Encode()
	return u.String(), nil
}
