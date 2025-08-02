package server

import (
	"log"
	"net/http"
	"os"

	"my_go_final_project/pkg/api" // импортируем наш пакет api

)

// Run запускает веб-сервер
func Run() {
	// определяем порт, который будет слушать сервер
	port := os.Getenv("TODO_PORT")
	if port == "" {
		// если порт не задан, используем по умолчанию
		port = "7540"
	}

	// регистрируем обработчики API
	api.Init()

	// отдаём статику из папки web
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Println("Запуск веб-сервера на порту", port)
	// запускаем сервер
	log.Fatal(http.ListenAndServe(":"+port, nil))
}