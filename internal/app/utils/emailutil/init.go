package emailutil

import (
	"net/smtp"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/thirdparty/logger"
	"github.com/jordan-wright/email"
)

var log = logger.WithModule("email")

type EmailUtil interface {
	SendForgetPasswordEmail(emailTo, name, surname, token string) error
	SendCreateUserEmail(emailTo, name, surname, tempPassword string) error
}

func NewEmailService(config app.Email) EmailUtil {
	mailPool, err := email.NewPool(config.Address,
		config.ConnectionPoolSize,
		smtp.PlainAuth("",
			config.Auth.Username,
			config.Auth.Password,
			config.Auth.Host))
	if err != nil {
		panic(err)
	}

	return &emailServiceImpl{
		email: mailPool,
	}
}
