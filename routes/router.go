package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pismo-assignment/handlers"
	"pismo-assignment/middleware"
	"pismo-assignment/repository"
	"pismo-assignment/services"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.APILogging())
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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
