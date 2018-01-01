# go-cexio-rest-api
Cex.io REST API written in Golang.

### Example

```go
package main

import "github/b1gb1t/go-cexio-rest-api"

func main() {
    api := cexioapi.NewCexioAPI("username","key","secret", false)

    // BTC last price
    res := api.LastPrice("BTC", "EUR")

    // Convert amount
    res := api.Converter("ETH", "EUR", "1337.0")
}
```
