package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type TransactionService struct {
	txRepo *repository.TransactionRepository
}

func NewTransactionService(txRep *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		txRepo: txRep,
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

func (s *TransactionService) BalanceInRupees(accountID uint) (float64, error) {
	sum, err := s.txRepo.SumAmountInPaisaByAccountID(accountID)
	if err != nil {
		return 0, err
	}
	return float64(sum) / 100.0, nil
}
