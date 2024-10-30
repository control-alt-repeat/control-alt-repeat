package access_tokens

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	ssmconfig "github.com/ianlopshire/go-ssm-config"

	"github.com/control-alt-repeat/control-alt-repeat/internal/aws"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type FreeagentAuthSecretsConfig struct {
	TokenUrl     string `ssm:"token_url" required:"true"`
	ClientID     string `ssm:"client_id" required:"true"`
	ClientSecret string `ssm:"client_secret" required:"true"`
	RefreshToken string `ssm:"refresh_token" required:"true"`
}

// Function to refresh the OAuth token using the refresh token
func refreshOAuthToken() (string, error) {
	var freeagentAuthSecrets FreeagentAuthSecretsConfig

	if err := os.Setenv("AWS_REGION", "eu-west-2"); err != nil {
		return "", err
	}
	if err := os.Setenv("AWS_DEFAULT_REGION", "eu-west-2"); err != nil {
		return "", err
	}
	if err := ssmconfig.Process("/control_alt_repeat/freeagent/live/", &freeagentAuthSecrets); err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(freeagentAuthSecrets.ClientID + ":" + freeagentAuthSecrets.ClientSecret))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", freeagentAuthSecrets.RefreshToken)

	req, err := http.NewRequest("POST", freeagentAuthSecrets.TokenUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var tokenResp TokenResponse
	if err = json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if err = aws.CreateOrUpdateSSMParameter(map[string]string{
		"/control_alt_repeat/freeagent/live/access_token": tokenResp.AccessToken,
		"/control_alt_repeat/freeagent/live/expires_in":   strconv.Itoa(tokenResp.ExpiresIn),
		"/control_alt_repeat/freeagent/live/timestamp":    strconv.Itoa(tokenResp.ExpiresIn),
	}); err != nil {
		return "", err
	}

	// Return the new access token
	return tokenResp.AccessToken, nil
}
