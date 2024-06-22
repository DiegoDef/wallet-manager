package testing

import (
	"bytes"
	"encoding/json"
	"wallet-manager/models"
	"wallet-manager/utils"

	"github.com/shopspring/decimal"
)

func createCryptocurrencyJson(crypto models.Cryptocurrency) *bytes.Reader {
	jsonBody, _ := json.Marshal(crypto)
	return bytes.NewReader(jsonBody)
}

func createCryptocurrency() models.Cryptocurrency {
	return createCryptoWithParamaters("Bitcoin", decimal.NewFromInt(1), decimal.NewFromInt(60000))
}

func createCryptoWithParamaters(name string, balance decimal.Decimal, costInFiat decimal.Decimal) models.Cryptocurrency {
	return models.Cryptocurrency{
		Name:        name,
		Balance:     balance,
		CostInFiat:  costInFiat,
		CreatedDate: utils.NowFormatted(),
	}
}
