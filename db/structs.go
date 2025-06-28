package db

import "database/sql"

type Pavilion struct {
	ID     int
	Number string
	Area   float64
}

type ActivityType struct {
	ID   int
	Name string
}

type User struct {
	TelegramID sql.NullInt64
	Username   sql.NullString
	Role       string
}

type Tenant struct {
	ID               int
	UserID           int
	FullName         string
	RegistrationType string
	HasCashRegister  bool
}
