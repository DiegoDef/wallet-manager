package testing

import (
	"bytes"
)

func createCryptocurrencyJson() *bytes.Reader {
	jsonBody := []byte(`{"name": "Bitcoin", "balance": 1, "CostInFiat": 60000}`)
	return bytes.NewReader(jsonBody)
}
