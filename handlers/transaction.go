package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"pismo-assignment/models"
	"pismo-assignment/services"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// @Summary      Create transaction
// @Description  Types 1–3 debit full amount, 4 credit. Type 2 is stored as installment purchase and debited like type 1 (full amount).
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateTransactionRequest  true  "Transaction payload"
// @Success      201   {object}  models.CreateTransactionResponse
// @Failure      400   {object}  models.ErrorResponse
// @Router       /transactions [post]
func (transactionHandler *TransactionHandler) CreateTransaction(c *gin.Context) {
	log.Println("POST /transactions")
	var createTransactionRequest models.CreateTransactionRequest
	if err := c.BindJSON(&createTransactionRequest); err != nil {
		errorResponse := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, errorResponse)
		log.Printf("Request: %+v", createTransactionRequest)
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("Request: %+v", createTransactionRequest)

	if createTransactionRequest.Amount <= 0 {
		errorResponse := gin.H{"error": "amount must be positive"}
		c.JSON(http.StatusBadRequest, errorResponse)
		log.Printf("Error: %v", errorResponse["error"])
		return
	}

	createdTransaction, err := transactionHandler.transactionService.Create(
		createTransactionRequest.AccountID,
		models.OperationType(createTransactionRequest.OperationTypeID),
		int64(createTransactionRequest.Amount*100),
	)
	if err != nil {
		errorResponse := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, errorResponse)
		log.Printf("Error: %v", err)
		return
	}

	createTransactionResponse := models.CreateTransactionResponse{
		TransactionID:   createdTransaction.ID,
		AccountID:       createdTransaction.AccountId,
		OperationTypeID: int(createdTransaction.OperationTypeId),
		Amount:          float64(createdTransaction.AmountInPaisa) / 100.0,
		EventDate:       createdTransaction.EventDate,
	}
	c.JSON(http.StatusCreated, createTransactionResponse)
	log.Printf("Response: %d %+v", http.StatusCreated, createTransactionResponse)
}
