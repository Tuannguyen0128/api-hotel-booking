package models

import (
	"time"
)

type Account struct {
	ID          string    `json:"id"`
	StaffID     string    `json:"staff_id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	UserRoleID  string    `json:"user_role_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}
