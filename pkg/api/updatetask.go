package api

import (
	"encoding/json"
	"net/http"
	"time"

	"my_go_final_project/pkg/db"
)

// updateTaskHandler обрабатывает обновление задачи.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// десериализуем JSON из тела запроса
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// проверяем обязательные поля
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	now := time.Now()

	// если дата не указана, берем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	// проверяем и обрабатываем дату
	taskDate, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат даты"})
		return
	}

	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		} else {
			task.Date = now.Format(dateLayout)
		}
	}

	// обновляем задачу в БД
	err = db.UpdateTask(task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// отправляем успешный пустой ответ
	writeJSON(w, map[string]string{})
}
