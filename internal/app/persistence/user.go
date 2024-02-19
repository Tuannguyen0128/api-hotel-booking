package persistence

import (
	"context"
	"time"
)

type UserDB interface {
	Insert(ctx context.Context, profile UserProfile) (string, error)
	GetByEmail(ctx context.Context, email string) (UserProfile, error)
	GetById(ctx context.Context, userId string) (UserProfile, error)
	GetByCompanyId(ctx context.Context, companyId string, filter UserFilter) ([]UserProfile, error)
	GetAll(ctx context.Context, filter UserFilter) ([]UserProfile, error)
	Update(ctx context.Context, userId string, document EditUserProfile) error
	UpdateInfo(ctx context.Context, updateReq EditUserProfileInfo) error
	WrongPasswordCounterReset(userId string) error
	WrongPasswordCounterIncrease(userId string, lockIfOver int, changeStatusTo string) (int, error)
}

type UserFilter struct {
	Status string
}

type UserProfile struct {
	Id                         string    `json:"id" bson:"_id,omitempty"`
	Email                      string    `json:"email" bson:"email"`
	Password                   string    `json:"password" bson:"password"`
	WrongPassword              int       `json:"-" bson:"wrong_password"`
	Name                       string    `json:"name" bson:"name"`
	Surname                    string    `json:"surname" bson:"surname"`
	CompanyId                  string    `json:"companyId" bson:"companyId"`
	Role                       string    `json:"role" bson:"role"`
	Status                     string    `json:"status" bson:"status"`
	OldPasswords               []string  `json:"oldPasswords" bson:"oldPasswords"`
	TempPassword               string    `json:"tempPassword" bson:"tempPassword"`
	CreateBy                   string    `json:"createBy" bson:"createBy"`
	CreateDt                   time.Time `json:"createDt" bson:"createDt"`
	LastEditDt                 time.Time `json:"lastEditDt" bson:"lastEditDt"`
	DeletedDt                  time.Time `json:"-" bson:"deletedDt"`
	ForgotPasswordEmailTimeOut time.Time `json:"-" bson:"forgotPasswordEmailTimeOut"`
}

func (u UserProfile) GetId() string {
	return u.Id
}

type EditUserProfile struct {
	Email                      string    `json:"email" bson:"email,omitempty"`
	Password                   string    `json:"password" bson:"password,omitempty"`
	Name                       string    `json:"name" bson:"name,omitempty"`
	Surname                    string    `json:"surname" bson:"surname,omitempty"`
	Role                       string    `json:"role" bson:"role,omitempty"`
	Status                     string    `json:"status" bson:"status,omitempty"`
	OldPasswords               []string  `json:"oldPasswords" bson:"oldPasswords,omitempty"`
	TempPassword               *string   `json:"tempPassword" bson:"tempPassword,omitempty"`
	CreateDt                   time.Time `json:"createDt" bson:"createDt,omitempty"`
	LastEditDt                 time.Time `json:"lastEditDt" bson:"lastEditDt,omitempty"`
	DeletedDt                  time.Time `json:"DeletedDt" bson:"DeletedDt,omitempty"`
	ForgotPasswordEmailTimeOut time.Time `json:"forgotPasswordEmailTimeOut" bson:"forgotPasswordEmailTimeOut,omitempty"`
}

type EditUserProfileInfo struct {
	UserId    string `bson:"-"`
	CompanyId string `bson:"-"`
	Role      string `bson:"-"`

	NewEmail         *string    `bson:"email,omitempty"`
	NewName          *string    `bson:"name,omitempty"`
	NewSurname       *string    `bson:"surname,omitempty"`
	NewRole          *string    `bson:"role,omitempty"`
	NewStatus        *string    `bson:"status,omitempty"`
	NewWrongPassword *int       `bson:"wrong_password,omitempty"`
	LastEditDt       *time.Time `bson:"lastEditDt,omitempty"`
	DeletedDt        *time.Time `bson:"deletedDt,omitempty"`
}
