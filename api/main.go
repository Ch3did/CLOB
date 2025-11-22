package main

import (
	"log"

	"CLOB/internal/config"
	"CLOB/internal/database"
	"CLOB/internal/http"
)

func main() {

	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	router := http.NewRouter(db)

	log.Println("Server running on port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
