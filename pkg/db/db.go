package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // регистрируем драйвер
)

// DB хранит активное соединение с базой данных.
var DB *sql.DB

// Схема для создания таблицы задач
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL,
	title VARCHAR(255) NOT NULL,
	comment TEXT,
	repeat VARCHAR(128) NOT NULL
);

CREATE INDEX IF NOT EXISTS date_idx ON scheduler (date);
`

// Init инициализирует соединение с БД и создает таблицы, если их нет
func Init() error {
	// получаем путь к файлу БД из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		// если переменная не задана, используем значение по умолчанию
		dbFile = "scheduler.db"
	}

	// открываем или создаем файл БД
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// выполняем SQL для создания таблиц и индексов
	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}

// AddTask добавляет новую задачу в базу данных.
func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	// получаем ID последней вставленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
