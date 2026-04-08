package handlers

import "time"

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

type CreateAccountResponse struct {
	AccountID      uint   `json:"account_id" example:"1"`
	DocumentNumber string `json:"document_number" example:"12345678900"`
}

type InstallmentPlanItem struct {
	PlanID        uint      `json:"plan_id"`
	TotalAmount   float64   `json:"total_amount"`
	Tenure        int       `json:"tenure"`
	PaidEMIs      int       `json:"paid_emis"`
	RemainingEMIs int       `json:"remaining_emis"`
	NextDueDate   time.Time `json:"next_due_date"`
}

type GetAccountResponse struct {
	AccountID              uint                  `json:"account_id" example:"1"`
	DocumentNumber         string                `json:"document_number" example:"12345678900"`
	Balance                float64               `json:"balance" example:"0"`
	ActiveInstallmentPlans []InstallmentPlanItem `json:"active_installment_plans"`
}

type CreateTransactionRequest struct {
	AccountID       uint    `json:"account_id" example:"1"`
	OperationTypeID int     `json:"operation_type_id" example:"1"`
	Amount          float64 `json:"amount" example:"123.45"`
	Tenure          *int    `json:"tenure,omitempty" example:"3"`
}

type CreateTransactionResponse struct {
	TransactionID   uint                 `json:"transaction_id"`
	AccountID       uint                 `json:"account_id"`
	OperationTypeID int                  `json:"operation_type_id"`
	Amount          float64              `json:"amount"`
	EventDate       time.Time            `json:"event_date"`
	InstallmentPlan *InstallmentPlanItem `json:"installment_plan,omitempty"`
}

type NextInstallmentResponse struct {
	TransactionID uint `json:"transaction_id"`
	PaidEMIs      int  `json:"paid_emis"`
	RemainingEMIs int  `json:"remaining_emis"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"record not found"`
}
