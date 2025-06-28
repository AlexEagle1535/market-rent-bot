package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./market.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Не удалось включить foreign_keys:", err)
	}

	queries := []string{
		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER UNIQUE,
			username TEXT UNIQUE,
			role TEXT NOT NULL CHECK (role IN ('admin', 'tenant'))
		);`,

		// Таблица арендаторов
		`CREATE TABLE IF NOT EXISTS tenants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
			full_name TEXT NOT NULL,
			registration_type TEXT NOT NULL,
			has_cash_register BOOLEAN NOT NULL DEFAULT 0
		);`,

		// Таблица сотрудников арендатора
		`CREATE TABLE IF NOT EXISTS tenant_employees (
			tenant_id INTEGER PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
			employee_count INTEGER NOT NULL,
			avg_salary DECIMAL NOT NULL
		);`,

		// Таблица павильонов
		`CREATE TABLE IF NOT EXISTS pavilions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pavilion_number TEXT NOT NULL UNIQUE,
			area REAL NOT NULL
		);`,

		// Договоры аренды
		`CREATE TABLE IF NOT EXISTS contracts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			pavilion_id INTEGER NOT NULL REFERENCES pavilions(id) ON DELETE CASCADE,
			contract_number TEXT NOT NULL,
			settlement_day DATE NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			rent_amount DECIMAL NOT NULL,
			amount_per_m2 DECIMAL NOT NULL
		);`,

		// Таблица кассовых аппаратов
		`CREATE TABLE IF NOT EXISTS cash_registers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			model TEXT NOT NULL,
			reg_number TEXT NOT NULL
		);`,

		// Таблица видов деятельности
		`CREATE TABLE IF NOT EXISTS activity_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);`,

		// Связь арендаторов и видов деятельности
		`CREATE TABLE IF NOT EXISTS tenant_activity_types (
			tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			activity_type_id INTEGER NOT NULL REFERENCES activity_types(id) ON DELETE CASCADE,
			PRIMARY KEY (tenant_id, activity_type_id)
		);`,

		// Договоры на коммунальные услуги
		`CREATE TABLE IF NOT EXISTS utilities_contracts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			contract_number TEXT NOT NULL,
			settlement_day DATE NOT NULL,
			start_date DATE,
			end_date DATE
		);`,
	}

	// Выполняем все запросы
	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatalf("Ошибка при создании таблиц: %v\nSQL: %s", err, query)
		}
	}
}
