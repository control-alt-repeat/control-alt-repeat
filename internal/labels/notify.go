package labels

import (
	"fmt"
	"net/http"

	"github.com/Control-Alt-Repeat/control-alt-repeat/internal/aws"
)

// SendPostRequest sends a POST request to the /webhook endpoint with no body
func NotifyLabelPrintServer() error {
	labelPrinterDomain, err := aws.GetParameterValue("eu-west-2", "/control_alt_repeat/ebay/live/label_printer/host_domain")
	if err != nil {
		return fmt.Errorf("error creating POST request: %v", err)
	}

	webhookURL := labelPrinterDomain + "/webhook"

	fmt.Println("Letting the printer know import complete at: ", webhookURL)

	// Create a new POST request with no body
	req, err := http.NewRequest("POST", webhookURL, nil)
	if err != nil {
		return fmt.Errorf("error creating POST request: %v", err)
	}
	go func() {
		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error sending POST request: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("received non-200 response: %d", resp.StatusCode)
			return
		}
	}()

	return nil
}
