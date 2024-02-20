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
	account := models.Account{}
	err = db.Debug().Model(models.Account{}).Where("email=?", email).Take(&account).Error
	if err != nil {
		return "", err
	}
	err = security.VerifyPassword(account.Password, password)
	if err != nil {
		return "", err
	}
	return CreateToken(email)
}
