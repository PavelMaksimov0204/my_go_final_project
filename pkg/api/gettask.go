package api

import (
	"net/http"

	"my_go_final_project/pkg/db"
)

// getTaskHandler обрабатывает запрос на получение одной задачи по ID.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем ID из параметров URL
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	// получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"})
		return
	}

	// отправляем задачу в виде JSON
	writeJSON(w, task)
}
