# telebirr_verifier

Go library for verifying Telebirr payment receipts.

## Install

```bash
go get github.com/Samuel-Fikre/telebirr_verifier
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
    }
}
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