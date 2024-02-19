package jwtutil

import (
	"api-hotel-booking/internal/app"
	"github.com/golang-jwt/jwt/v5"
	"time"

	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/rbac"
)

type JwtToken interface {
	GenerateToken() (string, error)
	ParseToken(tokenString string) error
	GetSessionId() string
	GetUserId() string
	GetEmail() string
	GetName() string
	GetSurname() string
	GetFullName() string
	GetCompanyId() string
	GetPermission() int
	GetRole() string
	GetExpire() string
	IsExpired() (bool, error)
}

func (c *customClaims) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(app.CFG.Service.Jwt.SigningKey))
	return ss, err
}

func (c *customClaims) ParseToken(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.CFG.Service.Jwt.SigningKey), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *customClaims) GetSessionId() string {
	return c.SessionId
}

func (c *customClaims) GetUserId() string {
	return c.UserId
}

func (c *customClaims) GetEmail() string {
	return c.Email
}

func (c *customClaims) GetName() string {
	return c.Name
}

func (c *customClaims) GetSurname() string {
	return c.Surname
}

func (c *customClaims) GetFullName() string {
	return c.Name + " " + c.Surname
}

func (c *customClaims) GetCompanyId() string {
	return c.CompanyId
}

func (c *customClaims) GetRole() string {
	return c.Role
}

func (c *customClaims) GetPermission() int {
	role, isIn := rbac.Roles[c.Role]
	if isIn == false {
		return 0
	}
	return role.Permission()
}

func (c *customClaims) GetExpire() string {
	return c.Expire
}

func (c *customClaims) IsExpired() (bool, error) {
	expire := c.GetExpire()
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
