package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// createOrUpdateSSMParameter creates or updates an SSM parameter with a secure string value
func createOrUpdateSSMParameter(ssmClient *ssm.SSM, name, value string) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(true),
	}

	_, err := ssmClient.PutParameter(input)
	if err != nil {
		return fmt.Errorf("failed to put parameter %s: %v", name, err)
	}
	return nil
}

func main() {
	// Load environment variables
	envVars := map[string]string{
		"/control_alt_repeat/freeagent/live/token_url":     "https://api.freeagent.com/v2/token_endpoint",
		"/control_alt_repeat/freeagent/live/client_id":     os.Getenv("FREEAGENT_LIVE_CLIENT_ID"),
		"/control_alt_repeat/freeagent/live/client_secret": os.Getenv("FREEAGENT_LIVE_CLIENT_SECRET"),
		"/control_alt_repeat/freeagent/live/refresh_token": os.Getenv("FREEAGENT_LIVE_REFRESH_TOKEN"),
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// Create SSM client
	ssmClient := ssm.New(sess)

	// Iterate over environment variables and create/replace SSM parameters
	for name, value := range envVars {
		if value == "" {
			log.Printf("Skipping parameter %s, no value found in environment variable", name)
			continue
		}

		log.Printf("Saving parameter %s with value from environment variable", name)
		err := createOrUpdateSSMParameter(ssmClient, name, value)
		if err != nil {
			log.Fatalf("failed to save parameter %s: %v", name, err)
		}
	}

	log.Println("All parameters saved successfully")
}
