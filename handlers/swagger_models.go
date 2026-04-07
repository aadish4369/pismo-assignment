package handlers

import "time"

// CreateAccountRequest is the body for POST /accounts.
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

// CreateAccountResponse is returned on successful account creation.
type CreateAccountResponse struct {
	AccountID      uint   `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

// GetAccountResponse is returned by GET /accounts/:accountId.
type GetAccountResponse struct {
	AccountID      uint   `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

// CreateTransactionRequest is the body for POST /transactions.
// For operation_type_id=2 (installment purchase), include tenure and start_date.
type CreateTransactionRequest struct {
	AccountID       uint    `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"1"`
	Amount          float64 `json:"amount" example:"123.45"`
	Tenure          int     `json:"tenure,omitempty" example:"3"`
	StartDate       string  `json:"start_date,omitempty" example:"2026-04-01"`
}

// CreateTransactionResponse is returned on successful transaction creation.
type CreateTransactionResponse struct {
	TransactionID   uint      `json:"transaction_id"`
	AccountID       uint      `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
	InstallmentID   *uint     `json:"installment_id,omitempty"`
	EMIAmount       *float64  `json:"emi_amount,omitempty"`
	RemainingEMIs   *int      `json:"remaining_emis,omitempty"`
}

// ErrorResponse is a common error payload.
type ErrorResponse struct {
	Error string `json:"error" example:"record not found"`
}

// PayEMIResponse is returned by POST /installments/:id/pay.
type PayEMIResponse struct {
	Status string `json:"status" example:"emi_paid"`
}
