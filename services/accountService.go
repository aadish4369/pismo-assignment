package services

import (
	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: repo,
	}
}

func (s *AccountService) CreateAccount(documentNumber string) (*models.Account, error) {
	account := &models.Account{
		DocumentNumber: documentNumber,
	}

	err := s.accountRepo.Create(account)
	return account, err
}

func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	return s.accountRepo.GetById(id)
}
