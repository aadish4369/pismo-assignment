package main

import (
	"os"

	"pismo-assignment/db"
	"pismo-assignment/routes"

	_ "pismo-assignment/docs"
)

// @title           Pismo transactions API
// @version         1.0
// @description     Accounts, transactions (signed amounts per operation type), and optional installment EMI extension.

// @contact.name   API Support

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	db.ConnectDatabase()
	if err := db.Migrate(); err != nil {
		panic(err)
	}

	r := routes.SetupRouter()
	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	r.Run(":" + addr)
}
