# telebirr_verifier

Go library for verifying Telebirr payment receipts.

## Install

```bash
# Latest version
go get github.com/Samuel-Fikre/telebirr_verifier@latest

# Or specific version
go get github.com/Samuel-Fikre/telebirr_verifier@v0.1.0
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/Samuel-Fikre/telebirr_verifier/pkg/loader"
    "github.com/Samuel-Fikre/telebirr_verifier/pkg/parser"
    "github.com/Samuel-Fikre/telebirr_verifier/pkg/receipt"
)

func main() {
    receiptNo := "CATExample"
    expectedAmount := 85.00
    expectedPhone := "251912345678"

    html, _ := loader.LoadReceipt(receiptNo, "")
    fields, _ := parser.ParseHTML(html)

    verifier := receipt.New(fields, map[string]interface{}{
        "credited_party_acc_no": expectedPhone,
        "settled_amount":        expectedAmount,
        "transaction_status":    "Completed",
    })

    if verifier.VerifyAll([]string{"date", "payer_name"}) {
        fmt.Println("✅ Payment verified!")
        fmt.Printf("Amount: %.2f Birr\n", fields["settled_amount"])
        fmt.Printf("From: %s (%s)\n", fields["payer_name"], fields["payer_phone"])
        fmt.Printf("To: %s (%s)\n", fields["credited_party_name"], fields["credited_party_acc_no"])
        fmt.Printf("Date: %s\n", fields["date"])
        fmt.Printf("Receipt number: %s\n", fields["receiptNo"])
        fmt.Printf("Payment reason: %s\n", fields["payment_reason"])
    } else {
        fmt.Println("❌ Payment verification failed!")
    }
}
```

Example output:
```
✅ Payment verified!
Amount: 85.00 Birr
From: John Doe (251912345678)
To: Your Business (251987654321)
Date: 29-01-2025 15:13:43
Receipt number: CATExample
Payment reason: Buy Package Mini APP
```

## Available Fields

| Field | Example |
|-------|---------|
| `settled_amount` | 85.00 |
| `credited_party_acc_no` | "251912345678" |
| `transaction_status` | "Completed" |
| `date` | "29-01-2025 15:13:43" |
| `payer_name` | "John Doe" |
| `payer_phone` | "251912345678" |

## License

MIT

## Disclaimer

Unofficial library, not affiliated with Telebirr or Ethio Telecom.