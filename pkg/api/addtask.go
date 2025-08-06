package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"my_go_final_project/pkg/db"
)

// addTaskHandler обрабатывает добавление новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// если дата в прошлом
	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		// и есть правило повторения, вычисляем следующую дату
		if task.Repeat != "" {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		} else {
			// иначе просто ставим сегодняшнюю
			task.Date = now.Format(dateLayout)
		}
	}

	// добавляем задачу в БД
	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// отправляем успешный ответ с ID новой задачи
	writeJSON(w, map[string]any{"id": strconv.FormatInt(id, 10)})
}
