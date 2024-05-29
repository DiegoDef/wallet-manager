package services

import (
	"fmt"
	"wallet-manager/models"
	"wallet-manager/repositories"
)

type CryptoTransactionService interface {
	Create(crypto *models.CryptoTransaction) error
	GetAll(cryptoId uint32) ([]models.CryptoTransaction, error)
	GetByID(id uint32) (*models.CryptoTransaction, error)
	Update(crypto *models.CryptoTransaction) error
	Delete(id uint32) error
}

type cryptoTransactionService struct {
	repo          repositories.CryptoTransactionRepository
	cryptoService CryptocurrencyService
}

func NewCryptoTransactionService(repo repositories.CryptoTransactionRepository, cryptoService CryptocurrencyService) CryptoTransactionService {
	return &cryptoTransactionService{repo: repo, cryptoService: cryptoService}
}

func (s *cryptoTransactionService) Create(crypto *models.CryptoTransaction) error {
	err := s.repo.Create(crypto)
	if err == nil {
		cryptocurrency := models.Cryptocurrency{ID: crypto.CryptocurrencyId, Balance: crypto.CryptocurrencyAmount, CostInFiat: crypto.FiatAmount}
		err = s.cryptoService.UpdateBalance(&cryptocurrency)
		if err != nil {
			fmt.Println("Erro ao atualizar balance")
		}
	}
	return err
}

func (s *cryptoTransactionService) GetAll(cryptoId uint32) ([]models.CryptoTransaction, error) {
	return s.repo.GetAll(cryptoId)
}

func (s *cryptoTransactionService) GetByID(id uint32) (*models.CryptoTransaction, error) {
	return s.repo.GetByID(id)
}

func (s *cryptoTransactionService) Update(crypto *models.CryptoTransaction) error {
	return s.repo.Update(crypto)
}

func (s *cryptoTransactionService) Delete(id uint32) error {
	return s.repo.Delete(id)
}
