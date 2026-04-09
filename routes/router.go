package routes

import (
	"github.com/gin-gonic/gin"

	"pismo-assignment/handlers"
	"pismo-assignment/repository"
	"pismo-assignment/services"
)

func SetupRouter() *gin.Engine {
	engine := gin.New()

	accountRepository := repository.NewAccountRepository()
	transactionRepository := repository.NewTransactionRepository()

	accountService := services.NewAccountService(accountRepository)
	transactionService := services.NewTransactionService(transactionRepository, accountRepository)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	registerAccountRoutes(engine, accountHandler)
	registerTransactionRoutes(engine, transactionHandler)

	return engine
}
