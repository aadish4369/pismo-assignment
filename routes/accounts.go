package routes

import (
	"github.com/gin-gonic/gin"

	"pismo-assignment/handlers"
)

func registerAccountRoutes(r *gin.Engine, h *handlers.AccountHandler) {
	r.POST("/accounts", h.CreateAccount)
	r.GET("/accounts/:accountId", h.GetAccount)
}
