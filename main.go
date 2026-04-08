package main

import (
	"os"

	"pismo-assignment/db"
	"pismo-assignment/routes"

	_ "pismo-assignment/docs"
)

// @title       Pismo API
// @version     1.0
// @description Accounts and transactions (SQLite, Gin, GORM).
// @host        localhost:8080
// @BasePath    /

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
