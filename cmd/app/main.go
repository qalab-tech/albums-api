package main

import (
	"log"

	"albums-api/internal/config"
	"albums-api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	app, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Не удалось инициализировать приложение: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}
}
