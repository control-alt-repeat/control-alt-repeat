package freeagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/internal/freeagent/access-tokens"
)

type requestOptions struct {
	Path            string
	QueryParameters map[string]string
}

func apiPost[T any](ctx context.Context, opts requestOptions, obj T) error {
	newURL, err := opts.buildUrl()
	if err != nil {
		return err
	}

	jsonPayload, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", newURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err

	}
	defer resp.Body.Close()

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
