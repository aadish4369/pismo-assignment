package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pismo-assignment/services"
)

type AccountHandler struct {
	accountSvc *services.AccountService
	txSvc      *services.TransactionService
}

func NewAccountHandler(accountSvc *services.AccountService, txSvc *services.TransactionService) *AccountHandler {
	return &AccountHandler{accountSvc: accountSvc, txSvc: txSvc}
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
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acct, err := h.accountSvc.Create(req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAccountResponse{
		AccountID:      acct.ID,
		DocumentNumber: acct.DocumentNumber,
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
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid accountId"})
		return
	}

	acct, err := h.accountSvc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.txSvc.BalanceInRupees(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetAccountResponse{
		AccountID:      acct.ID,
		DocumentNumber: acct.DocumentNumber,
		Balance:        balance,
	})
}
