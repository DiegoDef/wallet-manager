package models

import "github.com/shopspring/decimal"

type Cryptocurrency struct {
	ID               uint32          `json:"id" db:"cryptocurrency_id"`
	Name             string          `json:"name" db:"name"`
	Balance          decimal.Decimal `json:"balance" db:"balance"`
	CostInFiat       decimal.Decimal `json:"fiatBalance" db:"fiat_balance"`
	CreatedDate      string          `json:"createdDate" db:"created_date"`
	ProfitPercentage float32
}
