package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Task описывает задачу в планировщике.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Tasks возвращает список задач из базы данных с учетом поиска.
func Tasks(search string, limit int) ([]Task, error) {
	var rows *sql.Rows
	var err error

	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
		rows, err = DB.Query(query, t.Format("20060102"), limit)
	} else if search != "" {
		likeSearch := "%" + search + "%"
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, likeSearch, likeSearch, limit)
	} else {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса на получение задач: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования задачи: %w", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после итерации по задачам: %w", err)
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// GetTask возвращает одну задачу по её ID.
func GetTask(id string) (Task, error) {
	var t Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, fmt.Errorf("ошибка получения задачи по id %s: %w", id, err)
	}
	return t, nil
}

// UpdateTask обновляет задачу в базе данных.
func UpdateTask(task Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи id %s: %w", task.ID, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке обновленных строк для id %s: %w", task.ID, err)
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи id %s: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке удаленных строк для id %s: %w", id, err)
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет только дату задачи.
func UpdateDate(id, date string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, date, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты для id %s: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке обновленной даты для id %s: %w", id, err)
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

// AddTask добавляет новую задачу в базу данных.
func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения id новой задачи: %w", err)
	}

	return id, nil
}
