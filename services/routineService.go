package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type RoutineService struct {
	accountRepo *repository.AccountRepository
	txSvc       *TransactionService
	planRepo    *repository.InstallmentPlanRepository
}

func NewRoutineService() *RoutineService {
	accountRepo := repository.NewAccountRepository()
	txRepo := repository.NewTransactionRepository()
	planRepo := repository.NewInstallmentPlanRepository()
	return &RoutineService{
		accountRepo: accountRepo,
		txSvc:       NewTransactionService(txRepo),
		planRepo:    planRepo,
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

func (s *RoutineService) GetAccount(accountID uint) (*models.Account, float64, []models.InstallmentPlan, error) {
	acc, err := s.accountRepo.GetById(accountID)
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
	if _, err := s.accountRepo.GetById(accountID); err != nil {
		return nil, nil, err
	}
	if amountInRupees <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	totalPaisa := int64(amountInRupees * 100)

	if opID == models.InstallmentPurchase {
		if tenure == nil || *tenure <= 1 {
			return nil, nil, errors.New("tenure is required and must be greater than 1 for installment purchase")
		}
		t := *tenure
		emi := totalPaisa / int64(t)
		last := totalPaisa - emi*int64(t-1)
		firstPaisa := emi
		tx, err := s.txSvc.Create(accountID, opID, firstPaisa, nil)
		if err != nil {
			return nil, nil, err
		}
		plan := &models.InstallmentPlan{
			TransactionId: tx.ID,
			AccountId:     accountID,
			TotalPaisa:    totalPaisa,
			Tenure:        t,
			EMIPaisa:      emi,
			LastEMIPaisa:  last,
			PaidEMIs:      1,
			NextDueDate:   time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 1, 0),
		}
		if err := s.planRepo.Create(plan); err != nil {
			return nil, nil, err
		}
		tx.InstallmentPlanId = &plan.ID
		if err := s.txSvc.Save(tx); err != nil {
			return nil, nil, err
		}
		return tx, plan, nil
	}

	tx, err := s.txSvc.Create(accountID, opID, totalPaisa, nil)
	if err != nil {
		return nil, nil, err
	}
	return tx, nil, nil
}

func (s *RoutineService) RecordNextInstallmentEMI(accountID, planID uint) (*models.Transaction, *models.InstallmentPlan, error) {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return nil, nil, err
	}
	if plan.AccountId != accountID {
		return nil, nil, errors.New("installment plan not found for this account")
	}
	if plan.PaidEMIs >= plan.Tenure {
		return nil, nil, errors.New("installment plan already completed")
	}

	var paisa int64
	if plan.PaidEMIs+1 == plan.Tenure {
		paisa = plan.LastEMIPaisa
	} else {
		paisa = plan.EMIPaisa
	}

	tx, err := s.txSvc.Create(accountID, models.InstallmentPurchase, paisa, &planID)
	if err != nil {
		return nil, nil, err
	}

	plan.PaidEMIs++
	plan.NextDueDate = plan.NextDueDate.AddDate(0, 1, 0)
	if err := s.planRepo.Update(plan); err != nil {
		return nil, nil, err
	}
	return tx, plan, nil
}
