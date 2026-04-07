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

func (s *TransactionService) Create(
	accountID uint,
	opID models.OperationType,
	amountInPaisa int64,
	installmentID *uint,
) (*models.Transaction, error) {
	if opID < models.NormalPurchase || opID > models.CreditVoucher {
		return nil, errors.New("invalid operation_type_id")
	}

	amount := amountInPaisa
	if opID == models.CreditVoucher {
		if amount < 0 {
			amount = -amount
		}
	} else {
		if amount > 0 {
			amount = -amount
		}
	}

	tx := &models.Transaction{
		AccountId:       accountID,
		OperationTypeId: opID,
		AmountInPaisa:   amount,
		EventDate:       time.Now().UTC(),
		InstallmentId:   installmentID,
	}

	err := s.txRepo.Create(tx)
	return tx, err
}
