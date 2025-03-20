package main

import (
	"fmt"
	"log"

	"github.com/cbe_telebirr_verifier/pkg/loader"
	"github.com/cbe_telebirr_verifier/pkg/parser"
	"github.com/cbe_telebirr_verifier/pkg/receipt"
)

func main() {
	// Example receipt number (replace with actual receipt number)
	receiptNo := "YOUR_RECEIPT_NUMBER"

	// Load receipt HTML
	html, err := loader.LoadReceipt(receiptNo, "")
	if err != nil {
		log.Fatalf("Failed to load receipt: %v", err)
	}

	// Parse receipt HTML
	parsedFields, err := parser.ParseHTML(html)
	if err != nil {
		log.Fatalf("Failed to parse receipt: %v", err)
	}

	// Create predefined fields for verification
	preDefinedFields := map[string]interface{}{
		"to":     "John Doe",    // Expected recipient
		"amount": float64(1000), // Expected amount
	}

	// Create receipt verifier
	r := receipt.New(parsedFields, preDefinedFields)

	// Verify specific fields
	verified := r.Verify(func(parsed, predefined map[string]interface{}) bool {
		return r.Equals(parsed["to"], predefined["to"]) &&
			r.Equals(parsed["amount"], predefined["amount"])
	})

	if verified {
		fmt.Println("Receipt verification successful!")
		fmt.Printf("Recipient: %v\n", parsedFields["to"])
		fmt.Printf("Amount: %v\n", parsedFields["amount"])
	} else {
		fmt.Println("Receipt verification failed!")
	}

	// You can also use VerifyAll to check all fields
	// ignoring certain fields like payer_name
	allVerified := r.VerifyAll([]string{"payer_name"})
	fmt.Printf("All fields verified: %v\n", allVerified)
}
