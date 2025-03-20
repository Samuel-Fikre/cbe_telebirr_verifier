package parser

import (
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

// fieldMapping maps Amharic/English field names to struct field names
var fieldMapping = map[string]string{
	"የከፋይስም/payername":                       "payer_name",
	"የከፋይቴሌብርቁ./payertelebirrno.":            "payer_phone",
	"የከፋይአካውንትአይነት/payeraccounttype":         "payer_acc_type",
	"የገንዘብተቀባይስም/creditedpartyname":          "credited_party_name",
	"የገንዘብተቀባይቴሌብርቁ./creditedpartyaccountno": "credited_party_acc_no",
	"የክፍያውሁኔታ/transactionstatus":             "transaction_status",
	"የባንክአካውንትቁጥር/bankaccountnumber":         "bank_acc_no",
	"የክፍያቁጥር/receiptno.":                     "receiptNo",
	"የክፍያቀን/paymentdate":                     "date",
	"የተከፈለውመጠን/settledamount":                "settled_amount",
	"ቅናሽ/discountamount":                     "discount_amount",
	"15%ቫት/vat":                              "vat_amount",
	"ጠቅላላየተክፈለ/totalamountpaid":              "total_amount",
	"የገንዘቡልክበፊደል/totalamountinword":          "amount_in_word",
	"የክፍያዘዴ/paymentmode":                     "payment_mode",
	"የክፍያምክንያት/paymentreason":                "payment_reason",
	"የክፍያመንገድ/paymentchannel":                "payment_channel",
}

// ParseHTML parses the HTML content of a telebirr receipt
func ParseHTML(htmlContent string) (map[string]interface{}, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	fields := make(map[string]interface{})
	var tdNodes []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "td" {
			tdNodes = append(tdNodes, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for i := 0; i < len(tdNodes); i++ {
		text := getTextContent(tdNodes[i])
		text = normalizeText(text)

		if fieldName, exists := fieldMapping[text]; exists {
			if i+1 < len(tdNodes) {
				value := getTextContent(tdNodes[i+1])
				value = strings.TrimSpace(value)

				if fieldName == "bank_acc_no" {
					// Handle special case for bank account and recipient name
					accNo := strings.TrimSpace(extractNumbers(value))
					name := strings.TrimSpace(extractLetters(value))
					fields["bank_acc_no"] = accNo
					fields["to"] = name
				} else if strings.HasSuffix(fieldName, "amount") {
					// Handle amount fields
					amount := extractAmount(value)
					fields[fieldName] = amount
				} else {
					fields[fieldName] = value
				}
			}
		}
	}

	return fields, nil
}

// Helper functions
func getTextContent(n *html.Node) string {
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		}
		text += getTextContent(c)
	}
	return text
}

func normalizeText(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func extractNumbers(s string) string {
	var numbers strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			numbers.WriteRune(r)
		}
	}
	return numbers.String()
}

func extractLetters(s string) string {
	var letters strings.Builder
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == ' ' {
			letters.WriteRune(r)
		}
	}
	return letters.String()
}

func extractAmount(s string) float64 {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "birr", "")
	s = strings.TrimSpace(s)
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return amount
}
