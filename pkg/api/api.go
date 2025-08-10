package api

import (
	"encoding/json"
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
		// обрабатываем удаление задачи
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "не указан идентификатор"})
			return
		}
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]string{})
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
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

// writeJSON отправляет JSON-ответ.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
