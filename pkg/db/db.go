package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // регистрируем драйвер
)

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
func Init() (*sql.DB, error) {
	// получаем путь к файлу БД из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		// если переменная не задана, используем значение по умолчанию
		dbFile = "scheduler.db"
	}

	// открываем или создаем файл БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	// выполняем SQL для создания таблиц и индексов
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
