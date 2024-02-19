package jwtutil

import (
	"api-hotel-booking/internal/app"
	"github.com/golang-jwt/jwt/v5"
)

func (u *util) NewClaims(sessionId, userId, email, name, surname, companyId, role, expire string) JwtToken {
	return &customClaims{
		SessionId:        sessionId,
		UserId:           userId,
		Email:            email,
		Name:             name,
		Surname:          surname,
		CompanyId:        companyId,
		Role:             role,
		Expire:           expire,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
}

func (u *util) NewBlankClaims() JwtToken {
	return &customClaims{}
}

func (u *util) ParseAndVerifyToken(token string) (JwtToken, error) {
	c := &customClaims{}
	err := c.ParseToken(token)
	if err != nil {
		return c, app.ErrorMap.InvalidRequest.AddDebug("cannot parse the jwt token")
	}

	isExpired, err := c.IsExpired()
	if err != nil {
		return c, app.ErrorMap.InvalidRequest.AddDebug("error when parsing expire time")
	}
	if isExpired == true {
		return c, app.ErrorMap.SessionTimeout.AddDebug("this session is expired")
	}

	roleInt := c.GetPermission()
	if roleInt == 0 {
		return c, app.ErrorMap.InvalidRequest.AddDebug("cannot find this role")
	}

	return c, nil
}

func (u *util) ParseAndVerifyForgetPWToken(token string) (ForgetPWJwtToken, error) {
	f := &forgetPWClaims{}
	err := f.ParseToken(token)
	if err != nil {
		return f, app.ErrorMap.InvalidRequest.AddDebug("cannot parse the forget password jwt token")
	}

	isExpired, err := f.IsExpired()
	if err != nil {
		return f, app.ErrorMap.InvalidRequest.AddDebug("error when parsing expire time")
	}
	if isExpired == true {
		return f, app.ErrorMap.SessionTimeout.AddDebug("this session is expired")
	}

	return f, nil
}
