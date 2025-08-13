package api

import (
	"encoding/json"
	"log"
	"net/http"

	"my_go_final_project/pkg/db"
)

// taskHandler — главный обработчик для всех запросов /api/task.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
			return
		}
		if err := db.DeleteTask(id); err != nil {
			log.Printf("ошибка при удалении задачи: %v", err)
			writeJSON(w, map[string]string{"error": "внутренняя ошибка сервера"}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}, http.StatusOK)
	default:

		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
	}
}

// Init регистрирует все обработчики API.
func Init() {
	// защищенные маршруты
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))

	// открытые маршруты
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/signin", signinHandler)
}

// writeJSON отправляет JSON-ответ c указанным кодом статуса.
func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
