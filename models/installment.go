package models

import "time"

type Installment struct {
	BaseModel

	AccountId uint    `gorm:"not null;index"`
	Account   Account `gorm:"constraint:OnDelete:RESTRICT;"`

	TransactionId uint

	TotalAmountInPaisa int64

	Tenure int

	EMIAmountInPaisa int64

	LastEMIAmountInPaisa int64

	StartDate time.Time

	EndDate time.Time

	RemainingEMIs int

	NextDueDate time.Time
}
