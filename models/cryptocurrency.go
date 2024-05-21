package models

import "github.com/shopspring/decimal"

type Cryptocurrency struct {
	ID             uint32          `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	Balance        decimal.Decimal `json:"balance" db:"balance"`
	PurchaseDate   string          `json:"purchaseDate" db:"purchase_date"`
	PurchaseAmount decimal.Decimal `json:"purchaseAmount" db:"purchase_amount"`
	CreatedDate    string          `json:"createdDate" db:"created_date"`
}
