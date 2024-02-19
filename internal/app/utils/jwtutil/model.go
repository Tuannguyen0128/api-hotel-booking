package jwtutil

import "github.com/golang-jwt/jwt/v5"

type customClaims struct {
	SessionId string `json:"sessionId"`
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	CompanyId string `json:"companyId"`
	Role      string `json:"role"`
	Expire    string `json:"expire"`
	jwt.RegisteredClaims
}

type forgetPWClaims struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Expire string `json:"expire"`
	jwt.RegisteredClaims
}
