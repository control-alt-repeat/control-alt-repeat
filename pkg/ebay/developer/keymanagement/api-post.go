package keymanagement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/access-tokens"
)

type requestOptions struct {
	Path            string
	QueryParameters map[string]string
}

const BaseUrl = "https://apiz.ebay.com/developer/key_management/v1"

func apiPost[REQ any, RES any](ctx context.Context, opts requestOptions, obj REQ) (RES, error) {
	var responseObject RES

	newURL, err := opts.buildUrl()
	if err != nil {
		return responseObject, err
	}

	fmt.Println(newURL)

	jsonPayload, err := json.Marshal(obj)
	if err != nil {
		return responseObject, err
	}

	fmt.Println(string(jsonPayload))

	req, err := http.NewRequestWithContext(ctx, "POST", newURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return responseObject, err
	}

	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return responseObject, err
	}

	fmt.Println(access_token)

	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseObject, err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return responseObject, err
	}

	// Print the response status and headers
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)

	// Print the body (this might be JSON, HTML, etc.)
	fmt.Println("Response Body:", string(body))

	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return responseObject, err
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return responseObject, nil
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseObject, err
		}
		return responseObject, fmt.Errorf("client error %d occurred with status: %s", resp.StatusCode, body)
	case resp.StatusCode >= 500 && resp.StatusCode < 600:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return responseObject, err
		}
		return responseObject, fmt.Errorf("server error %d occurred with status: %s", resp.StatusCode, body)
	default:
		return responseObject, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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
