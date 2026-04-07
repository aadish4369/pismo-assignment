package handlers

import (
	"net/http"
	"strconv"
	"time"

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
// @Description  Creates an account for a document number (CPF-style identifier).
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
// @Description  Returns account id and document number.
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

	account, err := h.svc.GetAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

// CreateTransaction godoc
// @Summary      Create transaction
// @Description  Records a transaction. Amounts are normalized: types 1–3 stored negative, type 4 positive. For operation_type_id=2, pass tenure and start_date to create an installment plan.
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

	var startDate *time.Time
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
			return
		}
		startDate = &parsed
	}

	tx, inst, err := h.svc.CreateTransaction(
		req.AccountID,
		models.OperationType(req.OperationTypeID),
		req.Amount,
		req.Tenure,
		startDate,
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
	if inst != nil {
		emi := float64(inst.EMIAmountInPaisa) / 100.0
		rem := inst.RemainingEMIs
		resp.InstallmentID = &inst.ID
		resp.EMIAmount = &emi
		resp.RemainingEMIs = &rem
	}
	c.JSON(http.StatusCreated, resp)
}

// PayEMI godoc
// @Summary      Pay one EMI
// @Description  Marks one installment payment and records a withdrawal transaction for that EMI amount.
// @Tags         installments
// @Produce      json
// @Param        id   path      int  true  "Installment ID"
// @Success      200  {object}  PayEMIResponse
// @Failure      400  {object}  ErrorResponse
// @Router       /installments/{id}/pay [post]
func (h *APIHandler) PayEMI(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.DeductInstallmentEMI(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, PayEMIResponse{Status: "emi_paid"})
}
