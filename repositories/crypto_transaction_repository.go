package repositories

import (
	"wallet-manager/models"

	"github.com/jmoiron/sqlx"
)

type CryptoTransactionRepository interface {
	Create(transaction *models.CryptoTransaction) error
	GetAll(cryptoId uint32) ([]models.CryptoTransaction, error)
	GetByID(id uint32) (*models.CryptoTransaction, error)
	Update(crypto *models.CryptoTransaction) error
	Delete(id uint32) error
}

type cryptoTransactionRepository struct {
	db *sqlx.DB
}

func NewCryptoTransactionRepository(db *sqlx.DB) CryptoTransactionRepository {
	return &cryptoTransactionRepository{db: db}
}

func (r *cryptoTransactionRepository) Create(transaction *models.CryptoTransaction) error {
	query := `INSERT INTO crypto_transaction (cryptocurrency_id, cryptocurrency_amount, fiat_amount, purchase_date, created_date) 
			  VALUES (:cryptocurrency_id, :cryptocurrency_amount, :fiat_amount, :purchase_date, :created_date) RETURNING transaction_id`
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.Get(&transaction.ID, transaction)
}

func (r *cryptoTransactionRepository) GetAll(cryptoId uint32) ([]models.CryptoTransaction, error) {
	var cryptos []models.CryptoTransaction
	err := r.db.Select(&cryptos, "SELECT * FROM crypto_transaction WHERE cryptocurrency_id=$1", cryptoId)
	return cryptos, err
}

func (r *cryptoTransactionRepository) GetByID(id uint32) (*models.CryptoTransaction, error) {
	var crypto models.CryptoTransaction
	err := r.db.Get(&crypto, "SELECT * FROM crypto_transaction WHERE transaction_id=$1", id)
	return &crypto, err
}

func (r *cryptoTransactionRepository) Update(crypto *models.CryptoTransaction) error {
	query := `UPDATE crypto_transaction SET cryptocurrency_amount=:cryptocurrency_amount, fiat_amount=:fiat_amount, purchase_date=:purchase_date WHERE transaction_id=:transaction_id`
	_, err := r.db.NamedExec(query, crypto)
	return err
}

func (r *cryptoTransactionRepository) Delete(id uint32) error {
	_, err := r.db.Exec("DELETE FROM crypto_transaction WHERE transaction_id=$1", id)
	return err
}
