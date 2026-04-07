package repository

import (
	"pismo-assignment/db"
	"pismo-assignment/models"
	"time"
)

type InstallmentRepository struct{}

func NewInstallmentRepository() *InstallmentRepository {
	return &InstallmentRepository{}
}

func (r *InstallmentRepository) Create(inst *models.Installment) error {
	return db.DB.Create(inst).Error
}

func (r *InstallmentRepository) GetByID(id uint) (*models.Installment, error) {
	var inst models.Installment
	err := db.DB.First(&inst, id).Error
	return &inst, err
}

func (r *InstallmentRepository) Update(inst *models.Installment) error {
	return db.DB.Save(inst).Error
}

func (r *InstallmentRepository) GetDueInstallments(now time.Time) ([]models.Installment, error) {
	var inst []models.Installment

	err := db.DB.
		Where("next_due_date <= ?", now).
		Find(&inst).Error

	return inst, err
}
