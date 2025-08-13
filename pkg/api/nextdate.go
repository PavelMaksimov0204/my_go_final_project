package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// dateLayout определяет стандартный формат даты для приложения.
const dateLayout = "20060102"

// nextDateHandler обрабатывает GET-запросы к /api/nextdate.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "неверный формат даты 'now'", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}

// NextDate вычисляет следующую дату выполнения задачи на основе текущей даты и правила повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не может быть пустым")
	}

	d, err := time.Parse(dateLayout, date)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	now = now.Truncate(24 * time.Hour)

	switch rule {
	case "y":
		for {
			d = d.AddDate(1, 0, 0)
			if d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	case "d":
		if len(parts) < 2 {
			return "", errors.New("отсутствует количество дней для правила 'd'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("неверное количество дней: %s", parts[1])
		}
		if days <= 0 || days > 400 {
			return "", errors.New("количество дней должно быть от 1 до 400")
		}
		for {
			d = d.AddDate(0, 0, days)
			if d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	case "w":
		if len(parts) < 2 {
			return "", errors.New("не указаны дни недели для правила 'w'")
		}
		daysOfWeek := make(map[int]bool)
		for _, dayStr := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("неверный формат дней недели")
			}
			daysOfWeek[day] = true
		}

		for {
			d = d.AddDate(0, 0, 1)
			weekday := int(d.Weekday())
			// time.Weekday() считает воскресенье за 0, а в нашем формате это 7.
			if weekday == 0 {
				weekday = 7
			}
			if daysOfWeek[weekday] && d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	case "m":
		if len(parts) < 2 {
			return "", errors.New("не указаны дни месяца для правила 'm'")
		}
		daysOfMonth := make(map[int]bool)
		negativeDays := make(map[int]bool)
		for _, dayStr := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("неверный формат дней месяца")
			}
			if day > 0 && day <= 31 {
				daysOfMonth[day] = true
			} else if day == -1 || day == -2 {
				negativeDays[day] = true
			} else {
				return "", errors.New("неверный день месяца")
			}
		}

		validMonths := make(map[int]bool)
		allMonths := true
		if len(parts) > 2 {
			allMonths = false
			for _, monthStr := range strings.Split(parts[2], ",") {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("неверный формат месяцев")
				}
				validMonths[month] = true
			}
		}

		for {
			d = d.AddDate(0, 0, 1)
			if !allMonths && !validMonths[int(d.Month())] {
				continue
			}

			// Проверяем, подходит ли текущий день по правилу.
			dayIsValid := daysOfMonth[d.Day()]
			if !dayIsValid && len(negativeDays) > 0 {
				lastDayOfMonth := time.Date(d.Year(), d.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				if negativeDays[-1] && d.Day() == lastDayOfMonth {
					dayIsValid = true
				}
				if negativeDays[-2] && d.Day() == lastDayOfMonth-1 {
					dayIsValid = true
				}
			}

			if dayIsValid && d.After(now) {
				return d.Format(dateLayout), nil
			}
		}
	}

	return "", errors.New("неподдерживаемое правило повторения")
}
