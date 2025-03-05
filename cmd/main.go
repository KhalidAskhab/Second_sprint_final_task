package main

import (
	"log"

	"SSFT/internal/application"
)

func main() {
	app := application.New()
	log.Println("Запуск сервера...")
	if err := app.RunServer(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
