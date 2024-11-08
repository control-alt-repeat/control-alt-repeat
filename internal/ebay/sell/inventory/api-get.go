package inventory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/internal/ebay/access-tokens"
)

type requestOptions struct {
	Path            string
	QueryParameters map[string]string
}

const BaseUrl = "https://api.ebay.com/sell/inventory/v1"

func apiGet[T any](ctx context.Context, opts requestOptions, target *T) error {
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

	req, err := http.NewRequestWithContext(ctx, "GET", newURL, nil)
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

func apiPost[T any](ctx context.Context, opts requestOptions, obj T) error {
	newURL, err := opts.buildUrl()
	if err != nil {
		return err
	}

	fmt.Println(newURL)

	jsonPayload, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonPayload))

	req, err := http.NewRequestWithContext(ctx, "POST", newURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return err
	}

	fmt.Println(access_token[0:20])

	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err

	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("client error %d occurred with status: %s", resp.StatusCode, body)
	case resp.StatusCode >= 500 && resp.StatusCode < 600:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("server error %d occurred with status: %s", resp.StatusCode, body)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func (o requestOptions) buildUrl() (string, error) {
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

type ErrorResponse struct {
	Errors []struct {
		ErrorID     int    `json:"errorId"`
		Domain      string `json:"domain"`
		Category    string `json:"category"`
		Message     string `json:"message"`
		LongMessage string `json:"longMessage"`
		Parameters  []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters"`
	} `json:"errors"`
}
