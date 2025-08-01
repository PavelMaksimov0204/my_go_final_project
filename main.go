package main

import (
	"log"

	"my_go_final_project/pkg/db"
	"my_go_final_project/pkg/server"

)

func main() {
	// инициализируем базу данных
	_, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	// запускаем наш сервер
	server.Run()
}