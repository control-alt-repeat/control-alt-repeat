package finances

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	access_tokens "github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/access-tokens"
	"github.com/control-alt-repeat/control-alt-repeat/pkg/ebay/signing"
)

type requestOptions struct {
	Path            string
	QueryParameters map[string]string
}

const BaseUrl = "https://apiz.ebay.com/sell/finances/v1"

func apiGet[T any](ctx context.Context, opts requestOptions, target *T) error {
	method := "GET"
	access_token, err := access_tokens.GetAccessToken()
	if err != nil {
		return err
	}

	newURL, err := opts.buildUrl()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, newURL, nil)
	if err != nil {
		return err
	}

	privateKey := "TODO"
	jwe := "TODO"

	req.Header.Set("Authorization", "Bearer "+access_token)
	req.Header.Set("x-ebay-c-marketplace-id", "EBAY_GB")
	req.Header.Set("x-ebay-signature-key", jwe)

	timestamp := time.Now()

	signatureInput := signing.GenerateSignatureInput(nil, timestamp)
	req.Header.Set("Signature-Input", signatureInput)

	parsedURL, err := url.Parse(newURL)
	if err != nil {
		return err
	}

	signature, err := signing.GenerateSignature(req.Header, privateKey, signing.SignatureComponents{
		Method:    method,
		Path:      parsedURL.Path,
		Authority: parsedURL.Host,
	}, nil, timestamp)
	if err != nil {
		return err
	}

	req.Header.Set("Signature", signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
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
