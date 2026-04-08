package handlers

import "time"

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

type CreateAccountResponse struct {
	AccountID      uint   `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

type GetAccountResponse struct {
	AccountID      uint    `json:"account_id" example:"1"`
	DocumentNumber string  `json:"document_number" example:"12345678900"`
	Balance        float64 `json:"balance" example:"0"`
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

type ErrorResponse struct {
	Error string `json:"error" example:"record not found"`
}
