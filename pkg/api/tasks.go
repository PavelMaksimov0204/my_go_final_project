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
	// получаем задачи из БД, ограничиваемся 50-ю записями
	tasks, err := db.Tasks(50)
	if err != nil {
		// в случае ошибки отправляем JSON с ошибкой
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// отправляем успешный ответ со списком задач
	response := TasksResponse{Tasks: tasks}
	writeJSON(w, response)
}
