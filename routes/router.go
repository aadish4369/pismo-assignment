package routes

import (
	"github.com/gin-gonic/gin"

	"pismo-assignment/handlers"
	"pismo-assignment/repository"
	"pismo-assignment/services"
)

func SetupRouter() *gin.Engine {
	router := gin.New()

	accountRepo := repository.NewAccountRepository()
	txRepo := repository.NewTransactionRepository()

	accountSvc := services.NewAccountService(accountRepo)
	txSvc := services.NewTransactionService(txRepo, accountRepo)

	accountHandler := handlers.NewAccountHandler(accountSvc, txSvc)
	txHandler := handlers.NewTransactionHandler(txSvc)

	registerAccountRoutes(router, accountHandler)
	registerTransactionRoutes(router, txHandler)

	return router
}
