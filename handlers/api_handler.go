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

// CreateAccount godoc
// @Summary      Create account
// @Description  Creates an account for a document number.
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

// GetAccount godoc
// @Summary      Get account
// @Description  Returns account id, document number, balance, and active installment plans (incomplete).
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

// CreateTransaction godoc
// @Summary      Create transaction
// @Description  Types 1–3 debit (stored negative), type 4 credit (stored positive). For type 2, optional tenure (>1) creates an installment plan; omit tenure for a lump debit only. EMI repayment: POST .../installments/.../pay (credit voucher). Fails if balance would go negative on debits.
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

// PayInstallment godoc
// @Summary      Pay one installment
// @Description  Records a credit voucher for this EMI and advances the plan (does not change the original purchase debit).
// @Tags         accounts
// @Produce      json
// @Param        accountId  path      int  true  "Account ID"
// @Param        planId     path      int  true  "Installment plan ID"
// @Success      200        {object}  PayInstallmentResponse
// @Failure      400        {object}  ErrorResponse
// @Router       /accounts/{accountId}/installments/{planId}/pay [post]
func (h *APIHandler) PayInstallment(c *gin.Context) {
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

	plan, err := h.svc.PayInstallmentEMI(uint(accountID), uint(planID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, PayInstallmentResponse{
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
