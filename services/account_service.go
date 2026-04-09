package services

import (
	"errors"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Create(documentNumber string) (*models.Account, error) {
	if documentNumber == "" {
		return nil, errors.New("document_number is required")
	}
	account := &models.Account{DocumentNumber: documentNumber}
	err := s.repo.Create(account)
	return account, err
}

func (s *AccountService) GetByID(accountID uint) (*models.Account, error) {
	return s.repo.GetById(accountID)
}
