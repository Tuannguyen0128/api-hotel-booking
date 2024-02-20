package auth

import (
	"ProjectPractice/src/api/database"
	"ProjectPractice/src/api/models"
	"ProjectPractice/src/api/security"
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
