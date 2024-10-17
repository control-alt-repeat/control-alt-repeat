package access_tokens

import (
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

func GetAccessToken() (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
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

	expiryTime := access_token_last_modified.Add(time.Duration(access_token_expiry_seconds * int(time.Second)))

	if time.Now().UTC().After(expiryTime) {
		return refreshOAuthToken()
	}

	return access_token, nil
}
