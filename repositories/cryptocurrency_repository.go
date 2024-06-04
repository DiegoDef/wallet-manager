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
			  VALUES (LOWER(:name), :balance, :fiat_balance, :created_date) 
			  RETURNING cryptocurrency_id`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(&crypto.ID, crypto)
}

func (r *cryptocurrencyRepository) GetAll() ([]models.Cryptocurrency, error) {
	var cryptos []models.Cryptocurrency
	query := `
		SELECT c.*,
		((c.balance * p.price_usd)/ c.fiat_balance)*100 AS profit_percentage,
		(c.balance * p.price_usd) - c.fiat_balance AS usd_profit
		FROM cryptocurrency c
		INNER JOIN crypto_price p ON LOWER(c.name) = LOWER(p.name)
		ORDER BY profit_percentage DESC
	`
	err := r.db.Select(&cryptos, query)
	return cryptos, err
}

func (r *cryptocurrencyRepository) GetByID(id uint32) (*models.Cryptocurrency, error) {
	var crypto models.Cryptocurrency
	query := `
		SELECT c.*,
		((c.balance * cp.price_usd)/ c.fiat_balance)*100 AS profit_percentage,
		(c.balance * cp.price_usd) - c.fiat_balance AS usd_profit
		FROM cryptocurrency c
		INNER JOIN crypto_price cp ON LOWER(c.name) = LOWER(cp.name)
		WHERE c.cryptocurrency_id=$1
	`
	err := r.db.Get(&crypto, query, id)
	return &crypto, err
}

func (r *cryptocurrencyRepository) Update(crypto *models.Cryptocurrency) error {
	query := `UPDATE cryptocurrency SET name=LOWER(:name), balance=:balance, fiat_balance=:fiat_balance, created_date=:created_date WHERE cryptocurrency_id=:cryptocurrency_id`
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
