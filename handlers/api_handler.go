package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pismo-assignment/models"
	"pismo-assignment/services"
)

type APIHandler struct {
	svc *services.RoutineService
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{svc: services.NewRoutineService()}
}

// @Summary      Create account
// @Description  POST body: document_number
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      CreateAccountRequest  true  "Document number"
// @Success      201   {object}  CreateAccountResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /accounts [post]
func (h *APIHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.svc.CreateAccount(req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

// @Summary      Get account
// @Description  Balance and active installment plans
// @Tags         accounts
// @Produce      json
// @Param        accountId  path      int  true  "Account ID"
// @Success      200        {object}  GetAccountResponse
// @Failure      400        {object}  ErrorResponse
// @Failure      404        {object}  ErrorResponse
// @Router       /accounts/{accountId} [get]
func (h *APIHandler) GetAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid accountId"})
		return
	}

	account, balance, plans, err := h.svc.GetAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	items := make([]InstallmentPlanItem, 0, len(plans))
	for i := range plans {
		items = append(items, installmentPlanItemFromModel(&plans[i]))
	}

	c.JSON(http.StatusOK, GetAccountResponse{
		AccountID:              account.ID,
		DocumentNumber:         account.DocumentNumber,
		Balance:                balance,
		ActiveInstallmentPlans: items,
	})
}

// @Summary      Create transaction
// @Description  Types 1–3 debit, 4 credit. Type 2 needs tenure; first EMI debited. Use /installments/.../next for further EMIs.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body  body      CreateTransactionRequest  true  "Transaction payload"
// @Success      201   {object}  CreateTransactionResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /transactions [post]
func (h *APIHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, plan, err := h.svc.CreateTransaction(
		req.AccountID,
		models.OperationType(req.OperationTypeID),
		req.Amount,
		req.Tenure,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := CreateTransactionResponse{
		TransactionID:   tx.ID,
		AccountID:       tx.AccountId,
		OperationTypeID: int(tx.OperationTypeId),
		Amount:          float64(tx.AmountInPaisa) / 100.0,
		EventDate:       tx.EventDate,
	}
	if plan != nil {
		item := installmentPlanItemFromModel(plan)
		resp.InstallmentPlan = &item
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Next EMI
// @Description  Debit next EMI (type 2), update plan
// @Tags         accounts
// @Produce      json
// @Param        accountId  path      int  true  "Account ID"
// @Param        planId     path      int  true  "Installment plan ID"
// @Success      200        {object}  NextInstallmentResponse
// @Failure      400        {object}  ErrorResponse
// @Router       /accounts/{accountId}/installments/{planId}/next [post]
func (h *APIHandler) PostNextInstallment(c *gin.Context) {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid accountId"})
		return
	}
	planID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planId"})
		return
	}

	tx, plan, err := h.svc.RecordNextInstallmentEMI(uint(accountID), uint(planID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, NextInstallmentResponse{
		TransactionID: tx.ID,
		PaidEMIs:      plan.PaidEMIs,
		RemainingEMIs: plan.Tenure - plan.PaidEMIs,
	})
}

func installmentPlanItemFromModel(p *models.InstallmentPlan) InstallmentPlanItem {
	return InstallmentPlanItem{
		PlanID:        p.ID,
		TotalAmount:   float64(p.TotalPaisa) / 100.0,
		Tenure:        p.Tenure,
		PaidEMIs:      p.PaidEMIs,
		RemainingEMIs: p.Tenure - p.PaidEMIs,
		NextDueDate:   p.NextDueDate,
	}
}
