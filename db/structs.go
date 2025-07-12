package db

import (
	"database/sql"
	"time"
)

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

type Contract struct {
	ID             int
	PavilionNumber string
	ContractNumber string
	SigningDate    time.Time
	StartDate      time.Time
	EndDate        time.Time
	RentAmount     float64
}
