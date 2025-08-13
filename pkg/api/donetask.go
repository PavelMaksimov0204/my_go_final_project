package api

import (
	"log"
	"net/http"
	"time"

	"my_go_final_project/pkg/db"
)

// doneTaskHandler обрабатывает выполнение задачи.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("ошибка при удалении задачи: %v", err)
			writeJSON(w, map[string]string{"error": "внутренняя ошибка сервера"}, http.StatusInternalServerError)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}
		if err := db.UpdateDate(id, next); err != nil {
			log.Printf("ошибка при обновлении даты задачи: %v", err)
			writeJSON(w, map[string]string{"error": "внутренняя ошибка сервера"}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
