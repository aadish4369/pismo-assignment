package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type RoutineService struct {
	accountSvc *AccountService
	txSvc      *TransactionService
	planRepo   *repository.InstallmentPlanRepository
}

func NewRoutineService() *RoutineService {
	accountRepo := repository.NewAccountRepository()
	txRepo := repository.NewTransactionRepository()
	planRepo := repository.NewInstallmentPlanRepository()
	return &RoutineService{
		accountSvc: NewAccountService(accountRepo),
		txSvc:      NewTransactionService(txRepo),
		planRepo:   planRepo,
	}
}

func (s *RoutineService) CreateAccount(documentNumber string) (*models.Account, error) {
	if documentNumber == "" {
		return nil, errors.New("document_number is required")
	}
	return s.accountSvc.CreateAccount(documentNumber)
}

func (s *RoutineService) GetAccount(accountID uint) (*models.Account, float64, []models.InstallmentPlan, error) {
	acc, err := s.accountSvc.GetAccount(accountID)
	if err != nil {
		return nil, 0, nil, err
	}
	bal, err := s.txSvc.BalanceInRupees(accountID)
	if err != nil {
		return nil, 0, nil, err
	}
	plans, err := s.planRepo.ListActiveByAccount(accountID)
	if err != nil {
		return nil, 0, nil, err
	}
	return acc, bal, plans, nil
}

func (s *RoutineService) CreateTransaction(
	accountID uint,
	opID models.OperationType,
	amountInRupees float64,
	tenure *int,
) (*models.Transaction, *models.InstallmentPlan, error) {
	if _, err := s.accountSvc.GetAccount(accountID); err != nil {
		return nil, nil, err
	}
	if amountInRupees <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	amountInPaisa := int64(amountInRupees * 100)

	if opID == models.InstallmentPurchase && tenure != nil {
		if *tenure <= 1 {
			return nil, nil, errors.New("tenure must be greater than 1 when provided for installment purchase")
		}
		tx, err := s.txSvc.Create(accountID, opID, amountInPaisa)
		if err != nil {
			return nil, nil, err
		}
		t := *tenure
		emi := amountInPaisa / int64(t)
		last := amountInPaisa - emi*int64(t-1)
		plan := &models.InstallmentPlan{
			TransactionId: tx.ID,
			AccountId:     accountID,
			TotalPaisa:    amountInPaisa,
			Tenure:        t,
			EMIPaisa:      emi,
			LastEMIPaisa:  last,
			PaidEMIs:      0,
			NextDueDate:   time.Now().UTC().Truncate(24 * time.Hour),
		}
		if err := s.planRepo.Create(plan); err != nil {
			return nil, nil, err
		}
		return tx, plan, nil
	}

	tx, err := s.txSvc.Create(accountID, opID, amountInPaisa)
	if err != nil {
		return nil, nil, err
	}
	return tx, nil, nil
}

func (s *RoutineService) PayInstallmentEMI(accountID, planID uint) (*models.InstallmentPlan, error) {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return nil, err
	}
	if plan.AccountId != accountID {
		return nil, errors.New("installment plan not found for this account")
	}
	if plan.PaidEMIs >= plan.Tenure {
		return nil, errors.New("installment plan already completed")
	}

	var paisa int64
	if plan.PaidEMIs == plan.Tenure-1 {
		paisa = plan.LastEMIPaisa
	} else {
		paisa = plan.EMIPaisa
	}

	if _, err := s.txSvc.Create(accountID, models.CreditVoucher, paisa); err != nil {
		return nil, err
	}

	plan.PaidEMIs++
	plan.NextDueDate = plan.NextDueDate.AddDate(0, 1, 0)
	if err := s.planRepo.Update(plan); err != nil {
		return nil, err
	}
	return plan, nil
}
