package services

import (
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
	repo repositories.CryptoTransactionRepository
}

func NewCryptoTransactionService(repo repositories.CryptoTransactionRepository) CryptoTransactionService {
	return &cryptoTransactionService{repo: repo}
}

func (s *cryptoTransactionService) Create(crypto *models.CryptoTransaction) error {
	return s.repo.Create(crypto)
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
