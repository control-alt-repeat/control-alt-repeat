package models

import (
	access_tokens "github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay/access-tokens"
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
