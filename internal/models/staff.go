package models

import (
	"time"
)

type Staff struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Position    string    `json:"position"`
	Salary      float64   `json:"salary"`
	DateOfBirth string    `json:"date_of_birth"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	StartDate   time.Time `json:"start_date"`
	DeletedAt   time.Time `json:"deleted_at"`
}
