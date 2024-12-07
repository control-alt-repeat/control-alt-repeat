package keymanagement

import "context"

func CreateSigningKey(ctx context.Context) (CreateSigningKeyResponse, error) {
	apiopts := requestOptions{
		Path: "/signing_key",
	}

	response, err := apiPost[CreateSigningKeyRequest, CreateSigningKeyResponse](ctx, apiopts, CreateSigningKeyRequest{SigningKeyCipher: "RSA"})
	if err != nil {
		return CreateSigningKeyResponse{}, err
	}

	return response, nil
}

type CreateSigningKeyRequest struct {
	SigningKeyCipher string `json:"signingKeyCipher"`
}

type CreateSigningKeyResponse struct {
	CreationTime     string `json:"creationTime"`
	ExpirationTime   string `json:"expirationTime"`
	Jwe              string `json:"jwe"`
	PrivateKey       string `json:"privateKey"`
	PublicKey        string `json:"publicKey"`
	SigningKeyCipher string `json:"signingKeyCipher"`
	SigningKeyID     string `json:"signingKeyId"`
}
