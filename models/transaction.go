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

type CreateTransactionRequest struct {
	AccountID       uint    `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"1"`
	Amount          float64 `json:"amount" example:"123.45"`
}

type CreateTransactionResponse struct {
	TransactionID   uint      `json:"transaction_id"`
	AccountID       uint      `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}
