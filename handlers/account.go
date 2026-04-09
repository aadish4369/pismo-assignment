package handlers

import (
	"log"
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
	log.Println("POST /accounts")
	var req CreateAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Response: %d", http.StatusBadRequest)
		return
	}

	acct, err := h.accountSvc.Create(req.DocumentNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Response: %d %+v", http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := CreateAccountResponse{
		AccountID:      acct.ID,
		DocumentNumber: acct.DocumentNumber,
	}
	c.JSON(http.StatusCreated, answer)
	log.Printf("Response: %d %+v", http.StatusCreated, answer)
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
	log.Printf("GET /accounts/%s", c.Param("accountId"))
	id, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid accountId"})
		log.Printf("Response: %d", http.StatusBadRequest)
		return
	}

	acct, err := h.accountSvc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		log.Printf("Response: %d", http.StatusNotFound)
		return
	}

	balance, err := h.txSvc.BalanceInRupees(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		log.Printf("Response: %d", http.StatusNotFound)
		return
	}

	answer := GetAccountResponse{
		AccountID:      acct.ID,
		DocumentNumber: acct.DocumentNumber,
		Balance:        balance,
	}
	c.JSON(http.StatusOK, answer)
	log.Printf("Response: %d %+v", http.StatusOK, answer)
}
