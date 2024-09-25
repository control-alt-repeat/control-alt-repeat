package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/ebay"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func hello() error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: "eu-west-2",
	}))

	svc := s3.New(sess)

	rawObject, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("EBAY_TOKEN_S3_BUCKET")),
			Key:    aws.String(os.Getenv("EBAY_TOKEN_S3_KEY")),
		})

	buf := new(bytes.Buffer)
	buf.ReadFrom(rawObject.Body)
	myFileContentAsString := buf.String()

	err = ebay.ImportListing(myFileContentAsString)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(hello)
}
