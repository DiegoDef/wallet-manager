package services

import (
	"fmt"
	"strings"
	"sync"
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
	cryptos, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	prices, err := utils.GetCryptoPrices(utils.GetCryptoNames(cryptos))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate profit: %s", err)
	}

	var wg sync.WaitGroup
	// Canal para capturar erros
	errCh := make(chan error, len(cryptos))

	for i := range cryptos {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Processar cada criptomoeda e capturar erro, se houver
			err := calculateProfitPercentage(&cryptos[i], prices)
			if err != nil {
				errCh <- err
			}
		}(i)
	}

	// Aguardar a conclusão de todas as goroutines
	wg.Wait()
	// Fechar o canal de erros
	close(errCh)

	// Verificar se houve algum erro
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return cryptos, nil
}

func (s *cryptocurrencyService) GetByID(id uint32) (*models.Cryptocurrency, error) {
	crypto, err := s.repo.GetByID(id)
	if err != nil {
		return crypto, err
	}

	cryptoName := strings.ToLower(crypto.Name)
	prices, err := utils.GetCryptoPrices([]string{cryptoName})
	if err != nil {
		return nil, fmt.Errorf("failed to get crypto price: %s", err)
	}

	price, exists := prices[cryptoName]
	if !exists {
		return nil, fmt.Errorf("price for %s not found", cryptoName)
	}

	calculateProfitPercentage(crypto, map[string]decimal.Decimal{cryptoName: price})
	return crypto, err
}

func calculateProfitPercentage(crypto *models.Cryptocurrency, prices map[string]decimal.Decimal) error {
	if crypto.CostInFiat.IsZero() || crypto.Balance.IsZero() {
		return fmt.Errorf("failed to calculate profit. CostInFiat or balance is zero")
	}

	cryptoName := strings.ToLower(crypto.Name)
	price, exists := prices[cryptoName]
	if !exists {
		return fmt.Errorf("price for %s not found", cryptoName)
	}

	cost := crypto.CostInFiat
	fiatBalance := crypto.Balance.Mul(price)
	profit := fiatBalance.Sub(cost)

	profitDecimal := (profit.Div(cost)).Mul(decimal.NewFromInt(100))
	profitPercentage, _ := profitDecimal.Float64()

	crypto.ProfitPercentage = float32(profitPercentage)

	return nil
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
