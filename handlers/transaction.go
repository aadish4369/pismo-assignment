package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"pismo-assignment/models"
	"pismo-assignment/services"
)

type TransactionHandler struct {
	txSvc *services.TransactionService
}

func NewTransactionHandler(txSvc *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{txSvc: txSvc}
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
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	log.Println("POST /transactions")
	var req models.CreateTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		resp := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		log.Printf("Request: %+v", req)
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("Request: %+v", req)

	tx, err := h.txSvc.CreateFromRupees(
		req.AccountID,
		models.OperationType(req.OperationTypeID),
		req.Amount,
	)
	if err != nil {
		resp := gin.H{"error": err.Error()}
		c.JSON(http.StatusBadRequest, resp)
		log.Printf("Error: %v", err)
		return
	}

	resp := models.CreateTransactionResponse{
		TransactionID:   tx.ID,
		AccountID:       tx.AccountId,
		OperationTypeID: int(tx.OperationTypeId),
		Amount:          float64(tx.AmountInPaisa) / 100.0,
		EventDate:       tx.EventDate,
	}
	c.JSON(http.StatusCreated, resp)
	log.Printf("Response: %d %+v", http.StatusCreated, resp)
}
