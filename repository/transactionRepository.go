package repository

import (
	"pismo-assignment/db"
	"pismo-assignment/models"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	result := db.DB.Create(transaction)
	return result.Error
}

func (r *TransactionRepository) GetByAccountID(accountID uint) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := db.DB.
		Where("account_id = ?", accountID).
		Order("id desc").
		Find(&txs).Error

	return txs, err
}
