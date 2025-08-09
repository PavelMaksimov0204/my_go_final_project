package db

import (
	"database/sql"
	"errors"
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

	// Проверяем, является ли поисковый запрос датой.
	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		// Если это дата, ищем точное совпадение.
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
		rows, err = DB.Query(query, t.Format("20060102"), limit)
	} else if search != "" {
		// Если это текст, ищем вхождение подстроки.
		likeSearch := "%" + search + "%"
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, likeSearch, likeSearch, limit)
	} else {
		// Если поиска нет, возвращаем все задачи.
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
	return t, err
}

// UpdateTask обновляет задачу в базе данных.
func UpdateTask(task Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	// проверяем, была ли обновлена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return err
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
		return err
	}

	// проверяем, была ли удалена хотя бы одна запись
	count, err := res.RowsAffected()
	if err != nil {
		return err
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
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
