package services

import (
	"errors"
	"time"

	"pismo-assignment/models"
	"pismo-assignment/repository"
)

type InstallmentService struct {
	instRepo *repository.InstallmentRepository
}

func NewInstallmentService() *InstallmentService {
	return &InstallmentService{
		instRepo: repository.NewInstallmentRepository(),
	}
}

func (s *InstallmentService) GetByID(id uint) (*models.Installment, error) {
	return s.instRepo.GetByID(id)
}

func (s *InstallmentService) GetDueInstallments(now time.Time) ([]models.Installment, error) {
	return s.instRepo.GetDueInstallments(now)
}

func (s *InstallmentService) GetEMIAmount(inst *models.Installment) (int64, error) {
	if inst.RemainingEMIs <= 0 {
		return 0, errors.New("no remaining EMIs")
	}
	if inst.RemainingEMIs == 1 {
		return inst.LastEMIAmountInPaisa, nil
	}
	return inst.EMIAmountInPaisa, nil
}

func (s *InstallmentService) MarkEMIPaid(inst *models.Installment) error {
	if inst.RemainingEMIs <= 0 {
		return errors.New("no remaining EMIs")
	}
	inst.RemainingEMIs--
	inst.NextDueDate = inst.NextDueDate.AddDate(0, 1, 0)
	return s.instRepo.Update(inst)
}

func (s *InstallmentService) Create(
	accountID uint,
	transactionID uint,
	amount int64,
	tenure int,
	startDate time.Time,
) (*models.Installment, error) {

	emi := amount / int64(tenure)
	lastEMI := amount - (emi * int64(tenure-1))

	inst := &models.Installment{
		AccountId:            accountID,
		TransactionId:        transactionID,
		TotalAmountInPaisa:   amount,
		Tenure:               tenure,
		EMIAmountInPaisa:     emi,
		LastEMIAmountInPaisa: lastEMI,
		StartDate:            startDate,
		EndDate:              startDate.AddDate(0, tenure-1, 0),
		RemainingEMIs:        tenure,
		NextDueDate:          startDate,
	}

	err := s.instRepo.Create(inst)
	return inst, err
}

func (s *InstallmentService) DeductEMI(
	inst *models.Installment,
) (int64, error) {

	var amount int64

	if inst.RemainingEMIs == 1 {
		amount = inst.LastEMIAmountInPaisa
	} else {
		amount = inst.EMIAmountInPaisa
	}

	inst.RemainingEMIs--
	inst.NextDueDate = inst.NextDueDate.AddDate(0, 1, 0)

	err := s.instRepo.Update(inst)

	return amount, err
}
