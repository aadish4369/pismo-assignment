package repository

import (
	"pismo-assignment/db"
	"pismo-assignment/models"
)

type InstallmentPlanRepository struct{}

func NewInstallmentPlanRepository() *InstallmentPlanRepository {
	return &InstallmentPlanRepository{}
}

func (r *InstallmentPlanRepository) Create(plan *models.InstallmentPlan) error {
	return db.DB.Create(plan).Error
}

func (r *InstallmentPlanRepository) GetByID(id uint) (*models.InstallmentPlan, error) {
	var plan models.InstallmentPlan
	err := db.DB.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *InstallmentPlanRepository) Update(plan *models.InstallmentPlan) error {
	return db.DB.Save(plan).Error
}

func (r *InstallmentPlanRepository) ListActiveByAccount(accountID uint) ([]models.InstallmentPlan, error) {
	var plans []models.InstallmentPlan
	err := db.DB.
		Where("account_id = ? AND paid_emis < tenure", accountID).
		Order("id asc").
		Find(&plans).Error
	return plans, err
}
