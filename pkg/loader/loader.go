package loader

import (
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://transactioninfo.ethiotelecom.et/receipt/"

// LoadReceipt loads a receipt from the telebirr website using either a receipt number or full URL
func LoadReceipt(receiptNo string, fullURL string) (string, error) {
	var url string
	if receiptNo != "" {
		url = baseURL + receiptNo
	} else if fullURL != "" {
		url = fullURL
	} else {
		return "", fmt.Errorf("either receipt number or full URL must be provided")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch receipt: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch receipt: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
