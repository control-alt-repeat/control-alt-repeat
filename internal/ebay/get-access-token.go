package ebay

import ssmconfig "github.com/ianlopshire/go-ssm-config"

type EbayAccessToken struct {
	AccessToken string `ssm:"access_token" required:"false"`
	ExpiresIn   int    `ssm:"expires_in" default:"0"`
}

func getAccessToken() (string, error) {
	var ebayAccessToken EbayAccessToken

	err := ssmconfig.Process("/control_alt_repeat/ebay/live/", &ebayAccessToken)
	if err != nil {
		return "", err
	}

	if ebayAccessToken.ExpiresIn > 300 {
		return ebayAccessToken.AccessToken, nil
	}

	return refreshOAuthToken()
}
