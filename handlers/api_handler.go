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
// @Description  Balance in rupees
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

	account, balance, err := h.svc.GetAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
		Balance:        balance,
	})
}

// @Summary      Create transaction
// @Description  Types 1–3 debit full amount, 4 credit. Type 2 is stored as installment purchase and debited like type 1 (full amount).
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

	tx, err := h.svc.CreateTransaction(
		req.AccountID,
		models.OperationType(req.OperationTypeID),
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateTransactionResponse{
		TransactionID:   tx.ID,
		AccountID:       tx.AccountId,
		OperationTypeID: int(tx.OperationTypeId),
		Amount:          float64(tx.AmountInPaisa) / 100.0,
		EventDate:       tx.EventDate,
	})
}
