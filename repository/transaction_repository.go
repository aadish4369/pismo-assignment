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
