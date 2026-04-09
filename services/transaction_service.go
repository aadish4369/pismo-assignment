package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type TransactionService struct {
	txRepo      *repository.TransactionRepository
	accountRepo *repository.AccountRepository
}

func NewTransactionService(
	txRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *TransactionService {
	return &TransactionService{
		txRepo:      txRepo,
		accountRepo: accountRepo,
	}
}

func normalizedAmountInPaisa(opID models.OperationType, amountInPaisa int64) int64 {
	if opID == models.CreditVoucher {
		if amountInPaisa < 0 {
			return -amountInPaisa
		}
		return amountInPaisa
	}
	if amountInPaisa > 0 {
		return -amountInPaisa
	}
	return amountInPaisa
}

func (s *TransactionService) Create(
	accountID uint,
	opID models.OperationType,
	amountInPaisa int64,
) (*models.Transaction, error) {
	if opID < models.NormalPurchase || opID > models.CreditVoucher {
		return nil, errors.New("invalid operation_type_id")
	}

	amount := normalizedAmountInPaisa(opID, amountInPaisa)
	current, err := s.txRepo.SumAmountInPaisaByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	resultingBalance := current + amount
	if resultingBalance < 0 {
		return nil, errors.New("insufficient balance")
	}

	tx := &models.Transaction{
		AccountId:       accountID,
		OperationTypeId: opID,
		AmountInPaisa:   amount,
		EventDate:       time.Now().UTC(),
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *TransactionService) CreateFromRupees(
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
	return s.Create(accountID, opID, totalPaisa)
}

func (s *TransactionService) BalanceInRupees(accountID uint) (float64, error) {
	sum, err := s.txRepo.SumAmountInPaisaByAccountID(accountID)
	if err != nil {
		return 0, err
	}
	return float64(sum) / 100.0, nil
}
