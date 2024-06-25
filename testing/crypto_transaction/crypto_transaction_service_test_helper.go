package testing

import (
	"bytes"
	"encoding/json"
	"wallet-manager/models"
	"wallet-manager/repositories"
	"wallet-manager/utils"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

func createCryptoTransactionJson(transaction models.CryptoTransaction) *bytes.Reader {
	jsonBody, _ := json.Marshal(transaction)
	return bytes.NewReader(jsonBody)
}

func createCryptocurrency(testDbInstance *sqlx.DB) models.Cryptocurrency {
	cryptocurrencyRepository := repositories.NewCryptocurrencyRepository(testDbInstance)
	crypto := models.Cryptocurrency{
		Name:        "bitcoin",
		Balance:     decimal.NewFromInt(10),
		CostInFiat:  decimal.NewFromInt(10),
		CreatedDate: utils.NowFormatted(),
	}
	cryptocurrencyRepository.Create(&crypto)
	return crypto
}

func createTransactionWithoutCryptocurrencyId() models.CryptoTransaction {
	return models.CryptoTransaction{
		CryptocurrencyAmount: decimal.NewFromInt(10),
		FiatAmount:           decimal.NewFromInt(10),
		PurchaseDate:         utils.NowFormatted(),
		CreatedDate:          utils.NowFormatted(),
	}
}

func createTransaction(testDbInstance *sqlx.DB) models.CryptoTransaction {
	return models.CryptoTransaction{
		CryptocurrencyId:     createCryptocurrency(testDbInstance).ID,
		CryptocurrencyAmount: decimal.NewFromInt(10),
		FiatAmount:           decimal.NewFromInt(10),
		PurchaseDate:         utils.NowFormatted(),
		CreatedDate:          utils.NowFormatted(),
	}
}
