package models

import "github.com/shopspring/decimal"

type CryptoTransaction struct {
	ID                   uint32          `json:"id" db:"transaction_id"`
	CryptocurrencyId     uint32          `json:"cryptocurrency_id" db:"cryptocurrency_id"`
	CryptocurrencyAmount decimal.Decimal `json:"cryptocurrencyAmount" db:"cryptocurrency_amount"`
	FiatAmount           decimal.Decimal `json:"fiatAmount" db:"fiat_amount"`
	PurchaseDate         string          `json:"purchaseDate" db:"purchase_date"`
	CreatedDate          string          `json:"createdDate" db:"created_date"`
}
