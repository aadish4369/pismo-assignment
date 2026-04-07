package models

import "time"

type Transaction struct {
	BaseModel

	AccountId uint    `gorm:"not null;index"`
	Account   Account `gorm:"constraint:OnDelete:RESTRICT;"`

	OperationTypeId OperationType `gorm:"not null"`

	AmountInPaisa int64 `gorm:"not null"`
	EventDate     time.Time
}
