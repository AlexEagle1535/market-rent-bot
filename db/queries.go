package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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

func GetUsernameByID(userID int) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Нет такого пользователя
		}
		return "", err
	}
	return username, nil
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

func GetUsersByRole(role string) ([]User, error) {
	rows, err := DB.Query(`SELECT telegram_id, username, role FROM users WHERE role = $1`, role)
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

func SearchUsers(query string, roleFilter string) ([]User, error) {
	query = "%" + strings.ToLower(query) + "%"
	sql := `
		SELECT telegram_id, username, role
		FROM users
		WHERE (LOWER(username) LIKE ? OR CAST(telegram_id AS TEXT) LIKE ?)
	`

	args := []any{query, query}

	if roleFilter != "all" {
		sql += " AND role = ?"
		args = append(args, roleFilter)
	}

	rows, err := DB.Query(sql, args...)
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

// Добавление павильона
func AddPavilion(number string, area float64) error {
	_, err := DB.Exec(`INSERT INTO pavilions (pavilion_number, area) VALUES (?, ?)`, number, area)
	return err
}

// Получение всех павильонов
func GetAllPavilions() ([]Pavilion, error) {
	rows, err := DB.Query(`SELECT id, pavilion_number, area FROM pavilions ORDER BY pavilion_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pavilions []Pavilion
	for rows.Next() {
		var p Pavilion
		err := rows.Scan(&p.ID, &p.Number, &p.Area)
		if err != nil {
			return nil, err
		}
		pavilions = append(pavilions, p)
	}
	return pavilions, nil
}

func GetPavilionByID(id int) (*Pavilion, error) {
	var p Pavilion
	err := DB.QueryRow(`SELECT id, pavilion_number, area FROM pavilions WHERE id = ?`, id).Scan(&p.ID, &p.Number, &p.Area)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет такого павильона
		}
		return nil, err
	}
	return &p, nil
}

func GetPavilionByNumber(num string) (*Pavilion, error) {
	var p Pavilion
	err := DB.QueryRow(`SELECT id, pavilion_number, area FROM pavilions WHERE pavilion_number = ?`, num).Scan(&p.ID, &p.Number, &p.Area)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет такого павильона
		}
		return nil, err
	}
	return &p, nil
}

// Добавление вида деятельности
func AddActivityType(name string) error {
	_, err := DB.Exec(`INSERT INTO activity_types (name) VALUES (?)`, name)
	return err
}

// Получение всех видов деятельности
func GetAllActivityTypes() ([]ActivityType, error) {
	rows, err := DB.Query(`SELECT id, name FROM activity_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ActivityType
	for rows.Next() {
		var a ActivityType
		err := rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, a)
	}
	return types, nil
}

func GetAllTenants() ([]Tenant, error) {
	rows, err := DB.Query("SELECT * FROM tenants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []Tenant
	for rows.Next() {
		var tenant Tenant
		err := rows.Scan(&tenant.ID, &tenant.UserID, &tenant.FullName, &tenant.RegistrationType, &tenant.HasCashRegister)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func GetTenantByID(id int) (*Tenant, error) {
	var tenant Tenant
	err := DB.QueryRow("SELECT * FROM tenants WHERE id = ?", id).Scan(&tenant.ID, &tenant.UserID, &tenant.FullName, &tenant.RegistrationType, &tenant.HasCashRegister)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет такого арендатора
		}
		return nil, err
	}
	return &tenant, nil
}

func AddTenant(username, fullName, registrationType string, hasCashRegister bool) (int64, error) {
	// Проверка, есть ли пользователь с таким username
	var userID int64
	err := DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователя нет — создаём
			res, err := DB.Exec("INSERT INTO users (username, role) VALUES (?, 'tenant')", username)
			if err != nil {
				return 0, fmt.Errorf("не удалось создать пользователя: %w", err)
			}
			userID, err = res.LastInsertId()
			if err != nil {
				return 0, fmt.Errorf("не удалось получить ID нового пользователя: %w", err)
			}
		} else {
			return 0, fmt.Errorf("ошибка при поиске пользователя: %w", err)
		}
	}
	if err == nil {
		return 0, fmt.Errorf("пользователь с таким username уже существует")
	}
	// Добавляем арендатора
	res, err := DB.Exec(`
		INSERT INTO tenants (user_id, full_name, registration_type, has_cash_register)
		VALUES (?, ?, ?, ?)
	`, userID, fullName, registrationType, hasCashRegister)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить арендатора: %w", err)
	}

	tenantID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID арендатора: %w", err)
	}

	return tenantID, nil
}

func SaveTenantActivityTypes(tenantID int, activityTypeIDs []int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR IGNORE INTO tenant_activity_types (tenant_id, activity_type_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range activityTypeIDs {
		_, err := stmt.Exec(tenantID, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func AddTenantContract(tenantID int, contractNumber, pavilionNumber string, dateStartTime, dateEndTime time.Time, amount float64) error {
	// Получаем ID павильона по номеру
	var pavilionID int
	err := DB.QueryRow("SELECT id FROM pavilions WHERE pavilion_number = ?", pavilionNumber).Scan(&pavilionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("павильон с номером %s не найден", pavilionNumber)
		}
		return fmt.Errorf("ошибка при получении ID павильона: %w", err)
	}

	// Добавляем договор аренды
	_, err = DB.Exec(`
		INSERT INTO contracts (tenant_id, pavilion_id, contract_number, start_date, end_date, rent_amount)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenantID, pavilionID, contractNumber, dateStartTime, dateEndTime, amount)
	return err
}

func GetTenantActivityTypes(tenantID int) ([]ActivityType, error) {
	rows, err := DB.Query(`
		SELECT at.id, at.name
		FROM tenant_activity_types tat
		JOIN activity_types at ON tat.activity_type_id = at.id
		WHERE tat.tenant_id = ?
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ActivityType
	for rows.Next() {
		var a ActivityType
		err := rows.Scan(&a.ID, &a.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, a)
	}
	return types, nil
}

func AddCashRegister(tenantID int, model, cashRegisterNumber string) error {
	_, err := DB.Exec(`
		INSERT INTO cash_registers (tenant_id, model, reg_number)
		VALUES (?, ?, ?)
	`, tenantID, model, cashRegisterNumber)
	return err
}
