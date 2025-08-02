package api

import (
	"net/http"

)

// Init регистрирует все обработчики API.
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}