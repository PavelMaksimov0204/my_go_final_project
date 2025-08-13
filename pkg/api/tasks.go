package api

import (
	"log"
	"net/http"

	"my_go_final_project/pkg/db"

)

// TasksResponse определяет структуру ответа для списка задач.
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает запрос на получение списка задач.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(search, 50)
	if err != nil {
		log.Printf("ошибка при получении задач: %v", err)
		writeJSON(w, map[string]string{"error": "внутренняя ошибка сервера"}, http.StatusInternalServerError)
		return
	}

	response := TasksResponse{Tasks: tasks}
	writeJSON(w, response, http.StatusOK)
}