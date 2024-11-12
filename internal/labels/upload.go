package labels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/control-alt-repeat/control-alt-repeat/internal/aws"
	"github.com/control-alt-repeat/control-alt-repeat/internal/models"
)

type PrinterResponse struct {
	Model  string `json:"model"`
	Active bool   `json:"active"`
	Label  string `json:"label"`
}

func CheckPrinterOnline(ctx context.Context, labelSize string) (models.LabelPrinter, error) {
	labelPrinterDomain, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/label_printer/host_domain")
	if err != nil {
		return models.LabelPrinter{}, err
	}

	printerURL := fmt.Sprintf("%s/printer?label=%s", labelPrinterDomain, labelSize)

	req, err := http.NewRequestWithContext(ctx, "GET", printerURL, nil)
	if err != nil {
		return models.LabelPrinter{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.LabelPrinter{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.LabelPrinter{}, err
	}

	var printerResponse PrinterResponse

	if err := json.NewDecoder(resp.Body).Decode(&printerResponse); err != nil {
		return models.LabelPrinter{}, err
	}

	return models.LabelPrinter{
		Name:   printerResponse.Model,
		Active: printerResponse.Active,
		Format: printerResponse.Label,
	}, nil
}

func UploadFileFromBytes(imageBytes []byte, filename string) error {
	labelPrinterDomain, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/label_printer/host_domain")
	if err != nil {
		return err
	}

	printURL := labelPrinterDomain + "/print"

	// Create a buffer and multipart writer
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Create form file field
	formFile, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return err
	}

	// Write the byte array to the form file
	_, err = formFile.Write(imageBytes)
	if err != nil {
		return err
	}

	// Close the multipart writer
	if err = writer.Close(); err != nil {
		return err
	}

	// Create a request to upload the file
	req, err := http.NewRequest("POST", printURL, &body)
	if err != nil {
		return err
	}

	// Set the content type to multipart
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
