package models

type Account struct {
	BaseModel

	DocumentNumber string `gorm:"not null;uniqueIndex"`
	BalanceInPaisa int64  `gorm:"not null;default:0"`
}

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

type ErrorResponse struct {
	Error string `json:"error" example:"record not found"`
}
