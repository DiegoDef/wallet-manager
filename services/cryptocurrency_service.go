package services

import (
	"fmt"
	"wallet-manager/models"
	"wallet-manager/repositories"
	"wallet-manager/utils"

	"github.com/shopspring/decimal"
)

type CryptocurrencyService interface {
	Create(crypto *models.Cryptocurrency) error
	GetAll() ([]models.Cryptocurrency, error)
	GetByID(id uint32) (*models.Cryptocurrency, error)
	Update(crypto *models.Cryptocurrency) error
	UpdateBalance(crypto *models.Cryptocurrency) error
	Delete(id uint32) error
}

type cryptocurrencyService struct {
	repo repositories.CryptocurrencyRepository
}

func NewCryptocurrencyService(repo repositories.CryptocurrencyRepository) CryptocurrencyService {
	return &cryptocurrencyService{repo: repo}
}

func (s *cryptocurrencyService) Create(crypto *models.Cryptocurrency) error {
	return s.repo.Create(crypto)
}

func (s *cryptocurrencyService) GetAll() ([]models.Cryptocurrency, error) {
	return s.repo.GetAll()
}

func (s *cryptocurrencyService) GetByID(id uint32) (*models.Cryptocurrency, error) {
	crypto, err := s.repo.GetByID(id)
	if err != nil {
		return crypto, err
	}

	return calculateProfitPercentage(crypto)
}

func calculateProfitPercentage(crypto *models.Cryptocurrency) (*models.Cryptocurrency, error) {
	price, err := utils.GetCryptoPrice(crypto.Name)
	if err != nil {
		return crypto, fmt.Errorf("failed to calculate profit: %s", err)
	}

	cost := crypto.CostInFiat
	fiatBalance := crypto.Balance.Mul(price)
	profit := fiatBalance.Sub(cost)

	profitDecimal := (profit.Div(cost)).Mul(decimal.NewFromInt(100))
	profitPercentage, _ := profitDecimal.Float64()

	crypto.ProfitPercentage = float32(profitPercentage)

	return crypto, nil
}

func (s *cryptocurrencyService) Update(crypto *models.Cryptocurrency) error {
	return s.repo.Update(crypto)
}

func (s *cryptocurrencyService) UpdateBalance(crypto *models.Cryptocurrency) error {
	return s.repo.UpdateBalance(crypto)
}

func (s *cryptocurrencyService) Delete(id uint32) error {
	return s.repo.Delete(id)
}
