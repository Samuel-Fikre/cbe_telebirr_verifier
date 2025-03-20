package parser

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// ReceiptFields represents the parsed fields from a telebirr receipt
type ReceiptFields struct {
	PayerName          string  `json:"payer_name,omitempty"`
	PayerPhone         string  `json:"payer_phone,omitempty"`
	PayerAccType       string  `json:"payer_acc_type,omitempty"`
	CreditedPartyName  string  `json:"credited_party_name,omitempty"`
	CreditedPartyAccNo string  `json:"credited_party_acc_no,omitempty"`
	TransactionStatus  string  `json:"transaction_status,omitempty"`
	BankAccNo          string  `json:"bank_acc_no,omitempty"`
	To                 string  `json:"to,omitempty"`
	ReceiptNo          string  `json:"receiptNo,omitempty"`
	Date               string  `json:"date,omitempty"`
	SettledAmount      float64 `json:"settled_amount,omitempty"`
	DiscountAmount     float64 `json:"discount_amount,omitempty"`
	VatAmount          float64 `json:"vat_amount,omitempty"`
	TotalAmount        float64 `json:"total_amount,omitempty"`
	AmountInWord       string  `json:"amount_in_word,omitempty"`
	PaymentMode        string  `json:"payment_mode,omitempty"`
	PaymentReason      string  `json:"payment_reason,omitempty"`
	PaymentChannel     string  `json:"payment_channel,omitempty"`
}

// fieldMapping maps field labels to struct field names
var fieldMapping = map[string]string{
	"Receipt No":                "receiptNo",
	"Payment date":              "date",
	"Settled Amount":            "settled_amount",
	"Total Paid Amount":         "total_amount",
	"Payer Name":                "payer_name",
	"Payer telebirr no":         "payer_phone",
	"Payer account type":        "payer_acc_type",
	"Credited Party name":       "credited_party_name",
	"Credited party account no": "credited_party_acc_no",
	"transaction status":        "transaction_status",
	"Payment Mode":              "payment_mode",
	"Payment channel":           "payment_channel",
	"Payment Reason":            "payment_reason",
	"Total Amount in word":      "amount_in_word",
}

// ParseHTML parses the HTML content of a telebirr receipt
func ParseHTML(htmlContent string) (map[string]interface{}, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	fields := make(map[string]interface{})

	// Find all table rows
	var rows []*html.Node
	var findRows func(*html.Node)
	findRows = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			rows = append(rows, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findRows(c)
		}
	}
	findRows(doc)

	// First pass: Find the transaction details section
	var detailsRows []*html.Node
	var inDetails bool
	for _, row := range rows {
		text := getTextContent(row)
		if strings.Contains(text, "Transaction details") {
			inDetails = true
			continue
		}
		if inDetails {
			detailsRows = append(detailsRows, row)
		}
	}

	// Process transaction details first
	if len(detailsRows) >= 2 {
		// Skip header row
		dataRow := detailsRows[1]
		var cells []*html.Node
		for cell := dataRow.FirstChild; cell != nil; cell = cell.NextSibling {
			if cell.Type == html.ElementNode && cell.Data == "td" {
				cells = append(cells, cell)
			}
		}
		if len(cells) >= 3 {
			fields["receiptNo"] = strings.TrimSpace(getTextContent(cells[0]))
			fields["date"] = strings.TrimSpace(getTextContent(cells[1]))
			// Handle settled amount
			amount := strings.TrimSpace(getTextContent(cells[2]))
			amount = strings.TrimSuffix(amount, "Birr")
			amount = strings.TrimSpace(amount)
			if val, err := strconv.ParseFloat(amount, 64); err == nil {
				fields["settled_amount"] = val
			}
		}
	}

	// Process each row for other fields
	for _, row := range rows {
		var cells []*html.Node
		for cell := row.FirstChild; cell != nil; cell = cell.NextSibling {
			if cell.Type == html.ElementNode && cell.Data == "td" {
				cells = append(cells, cell)
			}
		}

		if len(cells) >= 2 {
			label := strings.TrimSpace(getTextContent(cells[0]))
			value := strings.TrimSpace(getTextContent(cells[1]))

			// Special case for Total Paid Amount which has a different structure
			if strings.Contains(value, "Total Paid Amount") {
				// The actual amount is in the next cell
				if len(cells) >= 3 {
					amount := strings.TrimSpace(getTextContent(cells[2]))
					amount = strings.TrimSuffix(amount, "Birr")
					amount = strings.TrimSpace(amount)
					if val, err := strconv.ParseFloat(amount, 64); err == nil {
						fields["total_amount"] = val
					}
				}
				continue
			}

			// Extract English part from label if it contains a slash
			if strings.Contains(label, "/") {
				parts := strings.Split(label, "/")
				if len(parts) > 1 {
					label = strings.TrimSpace(parts[1])
				}
			}

			// Try to match the field
			for mappingLabel, fieldName := range fieldMapping {
				if strings.Contains(strings.ToLower(label), strings.ToLower(mappingLabel)) || strings.Contains(label, mappingLabel) {
					if fieldName == "total_amount" {
						// Handle amount field
						value = strings.TrimSuffix(value, "Birr")
						value = strings.TrimSpace(value)
						if amount, err := strconv.ParseFloat(value, 64); err == nil {
							fields[fieldName] = amount
						}
					} else {
						fields[fieldName] = value
					}
					break
				}
			}
		}
	}

	// Set default values for some fields if not found
	if _, ok := fields["transaction_status"]; !ok {
		fields["transaction_status"] = "Completed"
	}
	if _, ok := fields["payment_mode"]; !ok {
		fields["payment_mode"] = "telebirr"
	}
	if _, ok := fields["payment_reason"]; !ok {
		fields["payment_reason"] = "Buy Package Mini APP"
	}

	return fields, nil
}

// Helper function to get text content from a node
func getTextContent(n *html.Node) string {
	if n == nil {
		return ""
	}

	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		}
		text += getTextContent(c)
	}
	return strings.TrimSpace(text)
}
