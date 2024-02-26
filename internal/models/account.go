package models

import (
	"time"
)

type Account struct {
	ID          string    `json:"id"json:",omitempty"`
	StaffID     string    `json:"staff_id,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	UserRoleID  string    `json:"user_role_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DeletedAt   time.Time `json:"deleted_at,omitempty"`
	LastLoginAt time.Time `json:"last_login_at,omitempty"`
}
