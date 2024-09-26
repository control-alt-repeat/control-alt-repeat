package ebay

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type EbayAccessToken struct {
	AccessToken string `ssm:"access_token" required:"false"`
	ExpiresIn   int    `ssm:"expires_in" default:"0"`
}

func getAccessToken() (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"), // Specify your region
	})
	if err != nil {
		return "", err
	}

	ssmClient := ssm.New(sess)

	result, err := ssmClient.GetParameters(&ssm.GetParametersInput{
		Names: aws.StringSlice([]string{
			"/control_alt_repeat/ebay/live/access_token",
			"/control_alt_repeat/ebay/live/expires_in",
		}),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	var access_token string
	var access_token_last_modified time.Time
	var access_token_expiry_seconds int

	for _, param := range result.Parameters {
		if *param.Name == "/control_alt_repeat/ebay/live/access_token" {
			access_token = *param.Value
			access_token_last_modified = *param.LastModifiedDate
		}
		if *param.Name == "/control_alt_repeat/ebay/live/expires_in" {
			access_token_expiry_seconds, err = strconv.Atoi(*param.Value)
			if err != nil {
				return "", err
			}
		}
	}

	fmt.Printf("Access token is: '%s...'\n", access_token[0:30])
	fmt.Printf("Access token last modified: %s\n", access_token_last_modified)
	fmt.Printf("Access token expiry seconds: %d\n", access_token_expiry_seconds)

	expiryTime := access_token_last_modified.Add(time.Duration(access_token_expiry_seconds * int(time.Second)))

	fmt.Printf("Access token expires at: %s\n", expiryTime)
	fmt.Printf("Access now is: %s\n", time.Now())

	if expiryTime.After(time.Now()) {
		fmt.Printf("Access token expired, getting a new one\n")

		return refreshOAuthToken()
	}

	fmt.Printf("Access token is valid, returning it for use.\n")

	return access_token, nil
}
