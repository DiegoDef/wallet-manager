package services

import (
	"wallet-manager/models"
	"wallet-manager/repositories"
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
	return cryptos, nil
}

func (s *cryptocurrencyService) GetByID(id uint32) (*models.Cryptocurrency, error) {
	return s.repo.GetByID(id)
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
