package signing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func VerifyGet(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/verifysignature", nil)
	if err != nil {
		return err
	}

	// These are not actual keys, they only work with the Docker test server.
	// See https://github.com/eBay/digital-signature-verification-ebay-api
	privateKey := "MC4CAQAwBQYDK2VwBCIEIJ+DYvh6SEqVTm50DFtMDoQikTmiCqirVv9mWG9qfSnF"
	jwe := "eyJ6aXAiOiJERUYiLCJlbmMiOiJBMjU2R0NNIiwidGFnIjoiSXh2dVRMb0FLS0hlS0Zoa3BxQ05CUSIsImFsZyI6IkEyNTZHQ01LVyIsIml2IjoiaFd3YjNoczk2QzEyOTNucCJ9.2o02pR9SoTF4g_5qRXZm6tF4H52TarilIAKxoVUqjd8.3qaF0KJN-rFHHm_P.AMUAe9PPduew09mANIZ-O_68CCuv6EIx096rm9WyLZnYz5N1WFDQ3jP0RBkbaOtQZHImMSPXIHVaB96RWshLuJsUgCKmTAwkPVCZv3zhLxZVxMXtPUuJ-ppVmPIv0NzznWCOU5Kvb9Xux7ZtnlvLXgwOFEix-BaWNomUAazbsrUCbrp514GIea3butbyxXLNi6R9TJUNh8V2uan-optT1MMyS7eMQnVGL5rYBULk.9K5ucUqAu0DqkkhgubsHHw"

	req.Header.Set("x-ebay-signature-key", jwe)

	timestamp := time.Now()

	signatureInput := GenerateSignatureInput(nil, timestamp)
	req.Header.Set("Signature-Input", signatureInput)

	signature, err := GenerateSignature(req.Header, privateKey, SignatureComponents{
		Method:    "GET",
		Path:      "/verifysignature",
		Authority: "localhost:8080",
	}, nil, timestamp)
	if err != nil {
		return err
	}

	req.Header.Set("Signature", signature)

	fmt.Println("Headers after signing:")
	for k, v := range req.Header {
		fmt.Printf("%s: %s\n", k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
