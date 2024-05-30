package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"wallet-manager/models"

	"github.com/shopspring/decimal"
)

type MultiCoinGeckoResponse map[string]struct {
	Usd float64 `json:"usd"`
}

func GetCryptoPrices(cryptoNames []string) (map[string]decimal.Decimal, error) {
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

	prices := make(map[string]decimal.Decimal, len(result))
	for name, data := range result {
		prices[name] = decimal.NewFromFloat(data.Usd)
	}

	return prices, nil
}

func GetCryptoNames(cryptos []models.Cryptocurrency) []string {
	var cryptoNames []string = make([]string, len(cryptos))
	for i, crypto := range cryptos {
		cryptoNames[i] = strings.ToLower(crypto.Name)
	}
	return cryptoNames
}
