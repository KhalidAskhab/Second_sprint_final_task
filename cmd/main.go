package main

import (
	"log"

	"Second_sprint_final_task/internal/application"
)

func main() {
	app := application.New()
	log.Println("Запуск сервера...")
	if err := app.RunServer(); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
