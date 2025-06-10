package db

import (
	"database/sql"
	"fmt"
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

// Получение роли по ID Telegram
func GetUserRole(telegramID int64, username string) (string, error) {
	var role string
	if telegramID == 0 && username == "" {
		return "", fmt.Errorf("не указан telegramID и username")
	}
	if telegramID != 0 || username != "" {
		err := DB.QueryRow("SELECT role FROM users WHERE telegram_id = ? OR username = ?", telegramID, username).Scan(&role)
		if err != nil {
			if err == sql.ErrNoRows {
				return "", nil // Нет такого пользователя
			}
			return "", err
		}
	}
	return role, nil
}

func SetUserRole(telegramID int64, username, role string) error {
	switch {
	case telegramID != 0 && username != "":
		// Обновляем по telegram_id, username просто обновляется тоже
		_, err := DB.Exec(`
			INSERT INTO users (telegram_id, username, role)
			VALUES (?, ?, ?)
			ON CONFLICT(telegram_id) DO UPDATE SET 
				username = excluded.username,
				role = excluded.role
		`, telegramID, username, role)
		return err

	case telegramID != 0:
		// Только по telegram_id
		_, err := DB.Exec(`
			INSERT INTO users (telegram_id, role)
			VALUES (?, ?)
			ON CONFLICT(telegram_id) DO UPDATE SET role = excluded.role
		`, telegramID, role)
		return err

	case username != "":
		// Только по username
		_, err := DB.Exec(`
			INSERT INTO users (username, role)
			VALUES (?, ?)
			ON CONFLICT(username) DO UPDATE SET role = excluded.role
		`, username, role)
		return err

	default:
		return fmt.Errorf("не указан ни telegramID, ни username")
	}
}
