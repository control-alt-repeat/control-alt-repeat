package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// GetSecrets returns the secrets stored in SSM
func GetSecrets(region string, keynames []*string) (*map[string]interface{}, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))

	withDecryption := true
	param, err := ssmsvc.GetParameters(&ssm.GetParametersInput{
		Names:          keynames,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return nil, err
	}

	secretsInfo := make(map[string]interface{})

	for _, item := range param.Parameters {
		secretsInfo[*item.Name] = *item.Value
	}
	return &secretsInfo, nil
}

func GetParameterValue(region string, keyname string) (string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return "", err
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))

	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(keyname),
		WithDecryption: aws.Bool(false),
	})

	return *param.Parameter.Value, err
}
