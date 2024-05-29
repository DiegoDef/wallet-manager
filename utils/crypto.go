package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type MultiCoinGeckoResponse map[string]struct {
	Usd float64 `json:"usd"`
}

func GetCryptoPrice(cryptoName string) (decimal.Decimal, error) {
	cryptoName = strings.ToLower(cryptoName)
	prices, err := GetMultipleCryptoPrices([]string{cryptoName})
	if err != nil {
		return decimal.Zero, err
	}

	price, exists := prices[cryptoName]
	if !exists {
		return decimal.Zero, fmt.Errorf("price for %s not found", cryptoName)
	}

	return decimal.NewFromFloat(price), nil
}

func GetMultipleCryptoPrices(cryptoNames []string) (map[string]float64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	ids := strings.Join(cryptoNames, ",")
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", ids)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}

	var result MultiCoinGeckoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	prices := make(map[string]float64, len(result))
	for name, data := range result {
		prices[name] = data.Usd
	}

	return prices, nil
}
