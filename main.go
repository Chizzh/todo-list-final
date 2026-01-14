package main

import (
	"log"

	"todo-list-final/pkg/db"
	"todo-list-final/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	// Init db
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}

	// Start server
	if err := server.Run("7540"); err != nil {
		panic(err)
	}
}
