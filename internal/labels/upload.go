package labels

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

func UploadFileFromBytes(imageBytes []byte, filename string) error {
	sizeMB := float64(len(imageBytes)) / (1024 * 1024) // Convert bytes to MB
	fmt.Printf("Data size: %.2f MB\n", sizeMB)

	labelPrinterDomain, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/label_printer/host_domain")
	if err != nil {
		return fmt.Errorf("error creating POST request: %v", err)
	}

	printURL := labelPrinterDomain + "/print"

	fmt.Println("Sending print job to: ", printURL)
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
		return fmt.Errorf("upload failed: %s", resp.Status)
	}

	fmt.Println("File uploaded successfully")
	return nil
}
