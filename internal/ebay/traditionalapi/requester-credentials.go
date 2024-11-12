package traditionalapi

import (
	access_tokens "github.com/control-alt-repeat/control-alt-repeat/internal/ebay/access-tokens"
)

type RequesterCredentials struct {
	EBayAuthToken string `xml:"eBayAuthToken"`
}

func (s *RequesterCredentials) SetEBayAuthToken() error {
	ebayAccessToken, err := access_tokens.GetAccessToken()

	if err != nil {
		return err
	}

	s.EBayAuthToken = ebayAccessToken

	return nil
}
