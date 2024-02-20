package models

import (
	"api-hotel-booking/internal/security"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/go-passwd/validator"
	"gorm.io/gorm"
)

type TeamMember struct {
	ID           uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Fullname     string    `gorm:"size:20;not null" json:"fullname"`
	Email        string    `gorm:"size:50;not null;unique" json:"email"`
	Password     string    `gorm:"size:60;not null" json:"password"`
	CreatedAt    time.Time `gorm:"autoCreatedTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdatedTime" json:"updated_at"`
	MerchantCode string    `json:"merchantcode"`
}

func (u *TeamMember) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	//log.Fatal(string(hashedPassword))
	return nil
}
func (t *TeamMember) Prepare() {
	t.ID = 0
	t.Fullname = html.EscapeString(strings.TrimSpace(t.Fullname))
	t.Email = html.EscapeString(strings.TrimSpace(t.Email))
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
}
func (t *TeamMember) Validtate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if t.Fullname == "" {
			return errors.New("Full name is require")

		}
		if t.Email == "" {
			return errors.New("Email is require")

		}
		if err := checkmail.ValidateFormat(t.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	case "login":

		if t.Password == "" {
			return errors.New("Password is require")
		}
		var err error
		passwordValidator := validator.New(validator.MinLength(5, err), validator.MaxLength(10, err))
		err = passwordValidator.Validate(t.Password)
		if err != nil {
			return errors.New("Password must be 5 -10 length  ")
		}
		if t.Email == "" {
			return errors.New("Email is require")

		}
		if err := checkmail.ValidateFormat(t.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	default:
		if t.Fullname == "" {
			return errors.New("Full name is require")

		}
		if t.Password == "" {
			return errors.New("Password is require")
		}
		var err error
		passwordValidator := validator.New(validator.MinLength(5, err), validator.MaxLength(10, err))
		err = passwordValidator.Validate(t.Password)
		if err != nil {
			return errors.New("Password must be 5 -10 length  ")
		}
		if t.Email == "" {
			return errors.New("Email is require")

		}
		if err := checkmail.ValidateFormat(t.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}
