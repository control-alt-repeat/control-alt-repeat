package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// createOrUpdateSSMParameter creates or updates an SSM parameter with a secure string value
func CreateOrUpdateSSMParameter(parameters map[string]string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	if err != nil {
		return err
	}

	// SSM client
	ssmClient := ssm.New(sess)

	for name, value := range parameters {
		if value == "" {
			fmt.Printf("Skipping parameter %s, no value found in environment variable", name)
			continue
		}

		// Define the input for PutParameter
		input := &ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String(value),
			Type:      aws.String("SecureString"),
			Overwrite: aws.Bool(true), // This allows the value to be overwritten if the parameter already exists
		}

		// Attempt to put the parameter
		_, err := ssmClient.PutParameter(input)
		if err != nil {
			return err
		}
	}

	return nil
}
