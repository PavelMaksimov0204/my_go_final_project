package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
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

// Init инициализирует соединение с БД.
func Init() error {
	var dbFile string
	if os.Getenv("RUNNING_IN_DOCKER") == "true" {
		dbFile = "file::memory:?cache=shared"
	} else {
		dbFile = os.Getenv("TODO_DBFILE")
		if dbFile == "" {
			dbFile = "scheduler.db"
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		// оборачиваем ошибку
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	_, err = DB.Exec(schema)
	if err != nil {
		// оборачиваем ошибку
		return fmt.Errorf("ошибка выполнения схемы: %w", err)
	}
	return nil
}
