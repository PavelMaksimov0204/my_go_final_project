package api

import (
	"net/http"
	"time"

	"my_go_final_project/pkg/db"
)

// doneTaskHandler обрабатывает выполнение задачи.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	// получаем задачу, чтобы проверить, есть ли у нее правило повторения
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"})
		return
	}

	// если правило повторения не задано, удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// если правило есть, вычисляем следующую дату
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		// и обновляем дату у задачи
		if err := db.UpdateDate(id, next); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// отправляем успешный пустой ответ
	writeJSON(w, map[string]string{})
}
