package main

import (
	"audioPhile/database"
	"audioPhile/server"
	"fmt"
)

func main() {
	err := database.ConnectAndMigrate("localhost", "5435", "audioPhile", "local", "local", database.SSLModeDisable)
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
