package models

import "time"

type InstallmentPlan struct {
	BaseModel

	TransactionId uint        `gorm:"not null;uniqueIndex"`
	Transaction   Transaction `gorm:"constraint:OnDelete:RESTRICT;"`

	AccountId uint    `gorm:"not null;index"`
	Account   Account `gorm:"constraint:OnDelete:RESTRICT;"`

	TotalPaisa   int64 `gorm:"not null"`
	Tenure       int   `gorm:"not null"`
	EMIPaisa     int64 `gorm:"not null"`
	LastEMIPaisa int64 `gorm:"not null"`
	PaidEMIs     int   `gorm:"not null;default:0;column:paid_emis"`
	NextDueDate  time.Time
}
