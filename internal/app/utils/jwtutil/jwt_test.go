package jwtutil

import (
	"fmt"
	"testing"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/rbac"
)

func Test_util_NewClaims(t *testing.T) {
	u := &util{}
	app.CFG.Service.Jwt.SigningKey = "mysigningkey"

	got := u.NewClaims("sessionId", "userId", "email", "name", "surname", "companyId", rbac.RoleAdmin.Name(), time.Now().Add(time.Hour*time.Duration(12)).Format(utils.RFC3339))

	to, _ := got.GenerateToken()

	_, e := u.ParseAndVerifyToken(to)

	fmt.Println("jwt", to)
	fmt.Println("print", e)
}
