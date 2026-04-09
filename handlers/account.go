package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pismo-assignment/models"
	"pismo-assignment/services"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// @Summary      Create account
// @Description  POST body: document_number
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateAccountRequest  true  "Document number"
// @Success      201   {object}  models.CreateAccountResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /accounts [post]
func (accountHandler *AccountHandler) CreateAccount(c *gin.Context) {
	log.Println("POST /accounts")
	var createAccountRequest models.CreateAccountRequest
	if err := c.BindJSON(&createAccountRequest); err != nil {
		errorResponse := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, errorResponse)
		log.Printf("Response: %d %+v", http.StatusBadRequest, errorResponse)
		return
	}
	log.Printf("Request: %+v", createAccountRequest)

	createdAccount, err := accountHandler.accountService.Create(createAccountRequest.DocumentNumber)
	if err != nil {
		errorResponse := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	createAccountResponse := models.CreateAccountResponse{
		AccountID:      createdAccount.ID,
		DocumentNumber: createdAccount.DocumentNumber,
	}
	c.JSON(http.StatusCreated, createAccountResponse)
	log.Printf("Response: %d %+v", http.StatusCreated, createAccountResponse)
}

// @Summary      Get account
// @Description  Balance in rupees
// @Tags         accounts
// @Produce      json
// @Param        accountId  path      int  true  "Account ID"
// @Success      200        {object}  models.GetAccountResponse
// @Failure      400        {object}  models.ErrorResponse
// @Failure      404        {object}  models.ErrorResponse
// @Router       /accounts/{accountId} [get]
func (accountHandler *AccountHandler) GetAccount(c *gin.Context) {
	log.Printf("GET /accounts/%s", c.Param("accountId"))
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		errorResponse := gin.H{"error": "invalid accountId"}
		c.JSON(http.StatusBadRequest, errorResponse)
		log.Printf("Response: %d", http.StatusBadRequest)
		return
	}
	log.Printf("Request: %d", accountID)

	account, err := accountHandler.accountService.GetByID(uint(accountID))
	if err != nil {
		errorResponse := gin.H{"error": err.Error()}
		c.JSON(http.StatusNotFound, errorResponse)
		log.Printf("Response: %d", http.StatusNotFound)
		return
	}

	getAccountResponse := models.GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
		Balance:        float64(account.BalanceInPaisa) / 100.0,
	}
	c.JSON(http.StatusOK, getAccountResponse)
	log.Printf("Response: %d %+v", http.StatusOK, getAccountResponse)
}
