package routes

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"pismo-assignment/handlers"
)

func SetupRouter() *gin.Engine {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		if strings.HasPrefix(absolutePath, "/swagger") {
			return
		}
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] %-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Skip: func(c *gin.Context) bool {
			return strings.HasPrefix(c.Request.URL.Path, "/swagger")
		},
	}))
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := handlers.NewAPIHandler()
	router.POST("/accounts", api.CreateAccount)
	router.POST("/accounts/:accountId/installments/:planId/next", api.PostNextInstallment)
	router.GET("/accounts/:accountId", api.GetAccount)
	router.POST("/transactions", api.CreateTransaction)

	return router
}
