package main

import (
	"audioPhile/database"
	"audioPhile/server"
	"fmt"
	"os"
)

func main() {
	host := os.Getenv("host")
	port := os.Getenv("port")
	databaseName := os.Getenv("database")
	user := os.Getenv("user")
	password := os.Getenv("password")
	err := database.ConnectAndMigrate(host, port, databaseName, user, password, database.SSLModeDisable)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
	srv := server.SetupRoutes()
	err = srv.Run(":8080")
	if err != nil {
		panic(err)
	}
}
