package db

import (
	"database/sql"
	"fmt"
)

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
		_, err := DB.Exec(`
			INSERT INTO users (telegram_id, role)
			VALUES (?, ?)
			ON CONFLICT(telegram_id) DO UPDATE SET role = excluded.role
		`, telegramID, role)
		return err

	case username != "":
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

func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT telegram_id, username, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.TelegramID, &u.Username, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func DeleteUser(telegramID int64, username string) error {
	if telegramID == 0 && username == "" {
		return fmt.Errorf("не указан telegramID и username")
	}
	if telegramID != 0 {
		_, err := DB.Exec("DELETE FROM users WHERE telegram_id = ?", telegramID)
		return err
	}
	if username != "" {
		_, err := DB.Exec("DELETE FROM users WHERE username = ?", username)
		return err
	}
	return nil
}

type User struct {
	TelegramID sql.NullInt64
	Username   sql.NullString
	Role       string
}
