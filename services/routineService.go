package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type RoutineService struct {
	accountSvc     *AccountService
	txSvc          *TransactionService
	installmentSvc *InstallmentService
}

func NewRoutineService() *RoutineService {
	accountRepo := repository.NewAccountRepository()
	txRepo := repository.NewTransactionRepository()
	return &RoutineService{
		accountSvc:     NewAccountService(accountRepo),
		txSvc:          NewTransactionService(txRepo),
		installmentSvc: NewInstallmentService(),
	}
}

func (s *RoutineService) CreateAccount(documentNumber string) (*models.Account, error) {
	if documentNumber == "" {
		return nil, errors.New("document_number is required")
	}
	return s.accountSvc.CreateAccount(documentNumber)
}

func (s *RoutineService) GetAccount(accountID uint) (*models.Account, error) {
	return s.accountSvc.GetAccount(accountID)
}

func (s *RoutineService) CreateTransaction(accountID uint, opID models.OperationType, amountInRupees float64, tenure int, startDate *time.Time) (*models.Transaction, *models.Installment, error) {
	if _, err := s.accountSvc.GetAccount(accountID); err != nil {
		return nil, nil, err
	}
	if amountInRupees <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	amountInPaisa := int64(amountInRupees * 100)
	tx, err := s.txSvc.Create(accountID, opID, amountInPaisa, nil)
	if err != nil {
		return nil, nil, err
	}

	// Optional EMI extension: create plan only for installment purchases.
	if opID != models.InstallmentPurchase {
		return tx, nil, nil
	}
	if tenure <= 1 {
		return nil, nil, errors.New("tenure must be greater than 1 for installment purchase")
	}
	if startDate == nil {
		return nil, nil, errors.New("start_date is required for installment purchase")
	}

	inst, err := s.installmentSvc.Create(accountID, tx.ID, amountInPaisa, tenure, *startDate)
	if err != nil {
		return nil, nil, err
	}
	return tx, inst, nil
}

func (s *RoutineService) DeductInstallmentEMI(installmentID uint) error {
	inst, err := s.installmentSvc.GetByID(installmentID)
	if err != nil {
		return err
	}

	amount, err := s.installmentSvc.GetEMIAmount(inst)
	if err != nil {
		return err
	}
	if err := s.installmentSvc.MarkEMIPaid(inst); err != nil {
		return err
	}
	_, err = s.txSvc.Create(inst.AccountId, models.Withdraw, amount, &inst.ID)
	return err
}
