package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Constants for eBay OAuth
const (
	tokenURL = "https://api.ebay.com/identity/v1/oauth2/token"
)

// Struct to capture the token response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Function to refresh the OAuth token using the refresh token
func refreshOAuthToken() (string, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")

	// Prepare basic authentication (base64 encoded clientID:clientSecret)
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Create form data for the token refresh request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
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

	// Return the new access token
	return tokenResp.AccessToken, nil
}

func main() {
	// Refresh the OAuth token using the refresh token
	token, err := refreshOAuthToken()
	if err != nil {
		fmt.Printf("Error refreshing token: %v\n", err)
		return
	}

	// Print the new access token
	fmt.Printf("New Access Token: %s\n", token)

	// Now you can use this token to make authenticated requests to eBay API
}
