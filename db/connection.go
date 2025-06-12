package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// Подключение к базе и создание таблицы
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./market.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)

	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER UNIQUE,
		username TEXT UNIQUE,
		role TEXT NOT NULL CHECK (role IN ('admin', 'tenant'))
	);`

	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Ошибка создания таблицы users:", err)
	}

}
