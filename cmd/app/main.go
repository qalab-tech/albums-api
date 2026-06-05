package main

import (
	"fmt"
	"log"

	"albums-api/internal/config"
	"albums-api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаём сервер
	srv := server.New(cfg)

	fmt.Println("Starting server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
