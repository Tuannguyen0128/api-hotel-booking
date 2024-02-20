package auth

import (
	"api-hotel-booking/internal/database"
	"api-hotel-booking/internal/models"
	"api-hotel-booking/internal/security"
)

func SignIn(email, password string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	teammember := models.TeamMember{}
	err = db.Debug().Model(models.TeamMember{}).Where("email=?", email).Take(&teammember).Error
	if err != nil {
		return "", err
	}
	err = security.VerifyPassword(teammember.Password, password)
	if err != nil {
		return "", err
	}
	return CreateToken(email)
}
