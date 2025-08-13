package api

import (
	"net/http"

	"my_go_final_project/pkg/db"

)

// getTaskHandler обрабатывает запрос на получение одной задачи по ID.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"}, http.StatusNotFound)
		return
	}

	// успешное получение - статус 200 OK
	writeJSON(w, task, http.StatusOK)
}
