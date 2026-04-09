package repository

import (
	"pismo-assignment/db"
	"pismo-assignment/models"
)

type AccountRepository struct{}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{}
}

func (r *AccountRepository) Create(account *models.Account) error {
	result := db.DB.Create(account)
	return result.Error
}

func (r *AccountRepository) GetById(id uint) (*models.Account, error) {
	var account models.Account
	result := db.DB.First(&account, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (r *AccountRepository) UpdateBalanceByID(id uint, balanceInPaisa int64) error {
	return db.DB.Model(&models.Account{}).
		Where("id = ?", id).
		Update("balance_in_paisa", balanceInPaisa).Error
}
