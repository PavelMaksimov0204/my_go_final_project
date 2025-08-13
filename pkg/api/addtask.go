package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"my_go_final_project/pkg/db"
)

// addTaskHandler обрабатывает добавление новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	taskDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат даты"}, http.StatusBadRequest)
		return
	}

	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}
			task.Date = next
		} else {
			task.Date = now.Format(dateLayout)
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		// Логируем настоящую ошибку
		log.Printf("ошибка при добавлении задачи: %v", err)
		// Пользователю отдаем общее сообщение
		writeJSON(w, map[string]string{"error": "внутренняя ошибка сервера"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}
