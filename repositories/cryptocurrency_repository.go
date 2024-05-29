package repositories

import (
	"wallet-manager/models"

	"github.com/jmoiron/sqlx"
)

type CryptocurrencyRepository interface {
	Create(crypto *models.Cryptocurrency) error
	GetAll() ([]models.Cryptocurrency, error)
	GetByID(id uint32) (*models.Cryptocurrency, error)
	Update(crypto *models.Cryptocurrency) error
	UpdateBalance(crypto *models.Cryptocurrency) error
	Delete(id uint32) error
}

type cryptocurrencyRepository struct {
	db *sqlx.DB
}

func NewCryptocurrencyRepository(db *sqlx.DB) CryptocurrencyRepository {
	return &cryptocurrencyRepository{db: db}
}

func (r *cryptocurrencyRepository) Create(crypto *models.Cryptocurrency) error {
	query := `INSERT INTO cryptocurrency (name, balance, fiat_balance, created_date) 
			  VALUES (:name, :balance, :fiat_balance, :created_date)`
	_, err := r.db.NamedExec(query, crypto)
	return err
}

func (r *cryptocurrencyRepository) GetAll() ([]models.Cryptocurrency, error) {
	var cryptos []models.Cryptocurrency
	err := r.db.Select(&cryptos, "SELECT * FROM cryptocurrency")
	return cryptos, err
}

func (r *cryptocurrencyRepository) GetByID(id uint32) (*models.Cryptocurrency, error) {
	var crypto models.Cryptocurrency
	err := r.db.Get(&crypto, "SELECT * FROM cryptocurrency WHERE cryptocurrency_id=$1", id)
	return &crypto, err
}

func (r *cryptocurrencyRepository) Update(crypto *models.Cryptocurrency) error {
	query := `UPDATE cryptocurrency SET name=:name, balance=:balance, fiat_balance=:fiat_balance, created_date=:created_date WHERE cryptocurrency_id=:cryptocurrency_id`
	_, err := r.db.NamedExec(query, crypto)
	return err
}

func (r *cryptocurrencyRepository) UpdateBalance(crypto *models.Cryptocurrency) error {
	query := `UPDATE cryptocurrency SET balance = :balance + balance, fiat_balance = :fiat_balance + fiat_balance WHERE cryptocurrency_id=:cryptocurrency_id`
	_, err := r.db.NamedExec(query, crypto)
	return err
}

func (r *cryptocurrencyRepository) Delete(id uint32) error {
	_, err := r.db.Exec("DELETE FROM cryptocurrency WHERE cryptocurrency_id=$1", id)
	return err
}
