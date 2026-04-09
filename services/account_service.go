package services

import (
	"errors"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type AccountService struct {
	accountRepository *repository.AccountRepository
}

func NewAccountService(accountRepository *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepository: accountRepository}
}

func (accountService *AccountService) Create(documentNumber string) (*models.Account, error) {
	if documentNumber == "" {
		return nil, errors.New("document_number is required")
	}
	account := &models.Account{DocumentNumber: documentNumber}
	err := accountService.accountRepository.Create(account)
	return account, err
}

func (accountService *AccountService) GetByID(accountID uint) (*models.Account, error) {
	return accountService.accountRepository.GetById(accountID)
}
