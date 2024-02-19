package jwtutil

import (
	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type ForgetPWJwtToken interface {
	GenerateToken() (string, error)
	ParseToken(tokenString string) error
	GetUserId() string
	GetEmail() string
	GetExpire() string
	IsExpired() (bool, error)
}

func (u *util) NewForgetPWClaims(userId, email, expire string) ForgetPWJwtToken {
	return &forgetPWClaims{
		UserId:           userId,
		Email:            email,
		Expire:           expire,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
}

func (f *forgetPWClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, f)
	ss, err := token.SignedString([]byte(app.CFG.Service.ForgetPW.SigningKey))
	return ss, err
}

func (f *forgetPWClaims) ParseToken(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, f, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.CFG.Service.ForgetPW.SigningKey), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *forgetPWClaims) GetUserId() string {
	return f.UserId
}

func (f *forgetPWClaims) GetEmail() string {
	return f.Email
}

func (f *forgetPWClaims) GetExpire() string {
	return f.Expire
}

func (f *forgetPWClaims) IsExpired() (bool, error) {
	expire := f.GetExpire()
	expireTime, err := time.Parse(utils.RFC3339, expire)
	if err != nil {
		return true, ParseTimeError
	}

	timeNow := time.Now()
	if timeNow.Equal(expireTime) || timeNow.After(expireTime) {
		return true, nil
	}
	return false, nil
}
