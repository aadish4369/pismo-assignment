package routes

import (
	"github.com/gin-gonic/gin"

	"pismo-assignment/handlers"
)

func registerTransactionRoutes(r *gin.Engine, h *handlers.TransactionHandler) {
	r.POST("/transactions", h.CreateTransaction)
}
