package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pismo-assignment/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := handlers.NewAPIHandler()
	router.POST("/accounts", api.CreateAccount)
	router.POST("/accounts/:accountId/installments/:planId/pay", api.PayInstallment)
	router.GET("/accounts/:accountId", api.GetAccount)
	router.POST("/transactions", api.CreateTransaction)

	return router
}
