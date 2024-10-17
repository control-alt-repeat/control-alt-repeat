package ebay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
	ssmconfig "github.com/ianlopshire/go-ssm-config"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type EbayAuthSecretsConfig struct {
	TokenUrl     string `ssm:"token_url" required:"true"`
	ClientID     string `ssm:"client_id" required:"true"`
	ClientSecret string `ssm:"client_secret" required:"true"`
	RefreshToken string `ssm:"refresh_token" required:"true"`
}

// Function to refresh the OAuth token using the refresh token
func refreshOAuthToken() (string, error) {
	var ebayAuthSecrets EbayAuthSecretsConfig

	os.Setenv("AWS_REGION", "eu-west-2")
	os.Setenv("AWS_DEFAULT_REGION", "eu-west-2")

	err := ssmconfig.Process("/control_alt_repeat/ebay/live/", &ebayAuthSecrets)
	if err != nil {
		return "", err
	}

	// Prepare basic authentication (base64 encoded clientID:clientSecret)
	auth := base64.StdEncoding.EncodeToString([]byte(ebayAuthSecrets.ClientID + ":" + ebayAuthSecrets.ClientSecret))

	// Create form data for the token refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", ebayAuthSecrets.RefreshToken)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", ebayAuthSecrets.TokenUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Add necessary headers
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+auth)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to refresh token: %s", string(body))
	}

	// Parse the token response
	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	aws.CreateOrUpdateSSMParameter(map[string]string{
		"/control_alt_repeat/ebay/live/access_token": tokenResp.AccessToken,
		"/control_alt_repeat/ebay/live/expires_in":   strconv.Itoa(tokenResp.ExpiresIn),
		"/control_alt_repeat/ebay/live/timestamp":    strconv.Itoa(tokenResp.ExpiresIn),
	})

	// Return the new access token
	return tokenResp.AccessToken, nil
}
