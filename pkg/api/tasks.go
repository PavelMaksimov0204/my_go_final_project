package api

import (
	"net/http"

	"my_go_final_project/pkg/db"
)

// TasksResponse определяет структуру ответа для списка задач.
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает запрос на получение списка задач.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр search из URL
	search := r.URL.Query().Get("search")

	// получаем задачи из БД с учетом поиска
	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	response := TasksResponse{Tasks: tasks}
	writeJSON(w, response)
}
