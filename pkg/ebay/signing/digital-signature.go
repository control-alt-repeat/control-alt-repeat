package signing

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const beginPrivateKey = "-----BEGIN PRIVATE KEY-----"
const endPrivateKey = "-----END PRIVATE KEY-----"

func getSignatureParams(payload []byte) []string {
	if payload != nil {
		return []string{"content-digest", "x-ebay-signature-key", "@method", "@path", "@authority"}
	}
	return []string{"x-ebay-signature-key", "@method", "@path", "@authority"}
}

func getSignatureParamsValue(payload []byte) string {
	params := getSignatureParams(payload)

	signatureParamsValue := ""

	for _, param := range params {
		signatureParamsValue += fmt.Sprintf(`"%s" `, param)
	}

	return strings.TrimSpace(signatureParamsValue)
}

func GenerateContentDigestValue(payload []byte, cipher string) string {
	payloadBuffer := bytes.NewBuffer(payload)

	hash := sha256.New()
	hash.Write(payloadBuffer.Bytes())
	digest := hash.Sum(nil)
	algo := "sha-256"
	if cipher == "sha512" {
		algo = "sha-512"
	}
	return fmt.Sprintf("%s=:%s:", algo, digest)
}

type SignatureComponents struct {
	Method    string
	Authority string
	Path      string
}

func GenerateBaseString(headers map[string]string, components SignatureComponents, payload []byte, timestamp time.Time) (string, error) {
	var baseString string

	signatureParams := getSignatureParams(payload)

	for _, param := range signatureParams {
		baseString += fmt.Sprintf(`"%s": `, strings.ToLower(param))

		if param[0] == '@' {
			switch strings.ToLower(param) {
			case "@method":
				baseString += components.Method
			case "@path":
				baseString += components.Path
			case "@authority":
				baseString += components.Authority
			default:
				return "", fmt.Errorf("invalid signature param: %s", param)
			}
		} else {
			if _, ok := headers[param]; !ok {
				return "", fmt.Errorf("header %s not included in message", param)
			}
			baseString += headers[param]
		}

		baseString += "\n"
	}

	baseString += fmt.Sprintf(`"@signature-params": (%s);created=%d`, getSignatureParamsValue(payload), timestamp.Unix())

	return baseString, nil
}

func GenerateSignatureInput(payload []byte, timestamp time.Time) string {
	return fmt.Sprintf(`sig1=(%s);created=%d`, getSignatureParamsValue(payload), timestamp.Unix())
}

func GenerateSignature(headers http.Header, privateKey string, components SignatureComponents, payload []byte, timestamp time.Time) (string, error) {
	digitalSignatureHeaders := map[string]string{
		"x-ebay-enforce-signature": "true",
		"x-ebay-signature-key":     headers.Get("x-ebay-signature-key"),
		"signature-input":          GenerateSignatureInput(payload, timestamp),
	}

	baseString, err := GenerateBaseString(digitalSignatureHeaders, components, payload, timestamp)
	if err != nil {
		return "", err
	}
	privateKey = strings.TrimSpace(privateKey)

	if !strings.HasPrefix(privateKey, beginPrivateKey) {
		privateKey = beginPrivateKey + "\n" + privateKey + "\n" + endPrivateKey
	}

	key, err := parsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	signatureBuffer, err := sign(baseString, key)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(signatureBuffer)

	return fmt.Sprintf(`sig1=:%s:`, signature), nil
}

func sign(payload string, privateKey crypto.PrivateKey) ([]byte, error) {
	payloadBytes := []byte(payload)

	hashed := sha256.Sum256(payloadBytes)

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
	case *ecdsa.PrivateKey:
		r, s, err := ecdsa.Sign(rand.Reader, key, hashed[:])
		if err != nil {
			return nil, err
		}
		return append(r.Bytes(), s.Bytes()...), nil
	case ed25519.PrivateKey:
		return ed25519.Sign(key, payloadBytes), nil
	}

	return nil, fmt.Errorf("unsupported private key type %T", privateKey)
}

func parsePrivateKey(privateKey string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("failed to parse private key")
}
