package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// dateLayout определяет формат для наших дат.
const dateLayout = "20060102"

// nextDateHandler обрабатывает эндпоинт /api/nextdate.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса.
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	// Если 'now' не передан, используем текущее время.
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "неверный формат даты 'now'", http.StatusBadRequest)
			return
		}
	}

	// Вычисляем следующую дату.
	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем ответ.
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}

// NextDate вычисляет следующую дату на основе правила повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Проверяем на пустое правило.
	if repeat == "" {
		return "", errors.New("правило повторения не может быть пустым")
	}

	// Парсим начальную дату.
	d, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	// Убираем время, чтобы сравнивать только даты.
	now = now.Truncate(24 * time.Hour)

	switch rule {
	case "y":
		// Ежегодное повторение.
		for {
			d = d.AddDate(1, 0, 0)
			if d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	case "d":
		// Повторение по дням.
		if len(parts) < 2 {
			return "", errors.New("отсутствует количество дней для правила 'd'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное количество дней: %s", parts[1])
		}
		if days > 400 {
			return "", errors.New("количество дней не может превышать 400")
		}
		for {
			d = d.AddDate(0, 0, days)
			if d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	}

	return "", errors.New("неподдерживаемое правило повторения")
}
