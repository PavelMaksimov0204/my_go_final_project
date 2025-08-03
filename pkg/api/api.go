package api

import (
	"encoding/json"
	"net/http"
)

// taskHandler — главный обработчик для всех запросов /api/task.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// если метод POST, вызываем обработчик добавления задачи
		addTaskHandler(w, r)
	default:
		// для всех других методов возвращаем ошибку
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Init регистрирует все обработчики API.
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	// регистрируем новый маршрут
	http.HandleFunc("/api/task", taskHandler)
}

// writeJSON отправляет JSON-ответ.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
