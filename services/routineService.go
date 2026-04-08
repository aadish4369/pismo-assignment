package services

import (
	"errors"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type RoutineService struct {
	accountRepo *repository.AccountRepository
	txSvc       *TransactionService
}

func NewRoutineService() *RoutineService {
	accountRepo := repository.NewAccountRepository()
	txRepo := repository.NewTransactionRepository()
	return &RoutineService{
		accountRepo: accountRepo,
		txSvc:       NewTransactionService(txRepo),
	}
}

func (s *RoutineService) CreateAccount(documentNumber string) (*models.Account, error) {
	if documentNumber == "" {
		return nil, errors.New("document_number is required")
	}
	account := &models.Account{DocumentNumber: documentNumber}
	err := s.accountRepo.Create(account)
	return account, err
}

func (s *RoutineService) GetAccount(accountID uint) (*models.Account, float64, error) {
	acc, err := s.accountRepo.GetById(accountID)
	if err != nil {
		return nil, 0, err
	}
	bal, err := s.txSvc.BalanceInRupees(accountID)
	if err != nil {
		return nil, 0, err
	}
	return acc, bal, nil
}

func (s *RoutineService) CreateTransaction(
	accountID uint,
	opID models.OperationType,
	amountInRupees float64,
) (*models.Transaction, error) {
	if _, err := s.accountRepo.GetById(accountID); err != nil {
		return nil, err
	}
	if amountInRupees <= 0 {
		return nil, errors.New("amount must be positive")
	}

	totalPaisa := int64(amountInRupees * 100)
	return s.txSvc.Create(accountID, opID, totalPaisa)
}
