package repositories

import (
	"wallet-manager/models"

	"github.com/jmoiron/sqlx"
)

const (
	insertCryptocurrencyQuery = `
		INSERT INTO cryptocurrency (name, balance, fiat_balance, created_date) 
		VALUES (LOWER(:name), :balance, :fiat_balance, :created_date) 
		RETURNING cryptocurrency_id;
	`
	getByIDQuery = `
		SELECT c.*,
		((c.balance * cp.price_usd)/ c.fiat_balance)*100 AS profit_percentage,
		(c.balance * cp.price_usd) - c.fiat_balance AS usd_profit
		FROM cryptocurrency c
		INNER JOIN crypto_price cp ON LOWER(c.name) = LOWER(cp.name)
		WHERE c.cryptocurrency_id=$1;
	`
	updateCryptocurrencyQuery = `
		UPDATE cryptocurrency 
		SET name=LOWER(:name), balance=:balance, fiat_balance=:fiat_balance, created_date=:created_date 
		WHERE cryptocurrency_id=:cryptocurrency_id;
	`
	updateCryptocurrencyBalanceQuery = `
		UPDATE cryptocurrency 
		SET balance = :balance + balance, fiat_balance = :fiat_balance + fiat_balance 
		WHERE cryptocurrency_id=:cryptocurrency_id;
	`
	deleteCryptocurrencyQuery = `DELETE FROM cryptocurrency WHERE cryptocurrency_id=$1;`
	getAllCryptocurrencyQuery = `
		SELECT c.*, ((c.balance * p.price_usd)/ c.fiat_balance)*100 AS profit_percentage,
		(c.balance * p.price_usd) - c.fiat_balance AS usd_profit
		FROM cryptocurrency c
		INNER JOIN crypto_price p ON LOWER(c.name) = LOWER(p.name)
		ORDER BY profit_percentage DESC;
	`
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
	stmt, err := r.db.PrepareNamed(insertCryptocurrencyQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(&crypto.ID, crypto)
}

func (r *cryptocurrencyRepository) GetAll() ([]models.Cryptocurrency, error) {
	var cryptos []models.Cryptocurrency
	err := r.db.Select(&cryptos, getAllCryptocurrencyQuery)
	return cryptos, err
}

func (r *cryptocurrencyRepository) GetByID(id uint32) (*models.Cryptocurrency, error) {
	var crypto models.Cryptocurrency
	err := r.db.Get(&crypto, getByIDQuery, id)
	return &crypto, err
}

func (r *cryptocurrencyRepository) Update(crypto *models.Cryptocurrency) error {
	_, err := r.db.NamedExec(updateCryptocurrencyQuery, crypto)
	return err
}

func (r *cryptocurrencyRepository) UpdateBalance(crypto *models.Cryptocurrency) error {
	_, err := r.db.NamedExec(updateCryptocurrencyBalanceQuery, crypto)
	return err
}

func (r *cryptocurrencyRepository) Delete(id uint32) error {
	_, err := r.db.Exec(deleteCryptocurrencyQuery, id)
	return err
}
