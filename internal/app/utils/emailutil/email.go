package emailutil

import (
	"api-hotel-booking/internal/app"
	"fmt"
	"github.com/jordan-wright/email"
	"time"
)

func (e *emailServiceImpl) SendForgetPasswordEmail(emailTo, name, surname, token string) error {
	em := email.NewEmail()
	em.To = []string{emailTo}
	em.Subject = "Car AI Core Cap : Labeling Tool - Forgot Password"
	line1 := "Dear, " + name + " " + surname
	line2 := "We received a request to reset the password for your account"
	line3 := "To reset your password, Click on the button below"
	line4 := app.CFG.Email.Url + "/reset-password?token=" + token
	line5 := "Or copy and paste the URL into your browser :"
	html := fmt.Sprintf("<div>%s</div><div>%s</div><div>%s</div><a href=\"%s\"><button style=\"background-color:#04AA6D;color: white;border-radius: 5px;border: none;\">RESET YOUR PASSWORD NOW</button></a><div>%s</div><div>%s</div>", line1, line2, line3, line4, line5, line4)
	em.HTML = []byte(html)
	go e.send(em)
	return nil
}

func (e *emailServiceImpl) SendCreateUserEmail(emailTo, name, surname, tempPassword string) error {
	em := email.NewEmail()
	em.To = []string{emailTo}
	em.Subject = "Car AI Core Cap : Labeling Tool - Registration Confirmed"
	line1 := "Dear, " + name + " " + surname
	line2 := "Welcome to Car AI Core Cap : Labeling Tool, you will be able to log in using :"
	line3 := "Password : " + tempPassword
	html := fmt.Sprintf("<div>%s</div><div>%s</div><div>%s</div>", line1, line2, line3)
	em.HTML = []byte(html)
	go e.send(em)
	return nil
}

func (e *emailServiceImpl) send(em *email.Email) {
	em.From = app.CFG.Email.Sender
	retryDeadline := time.Now().Add(time.Duration(app.CFG.Email.RetryTimeout) * time.Second)
	sendTimeout := time.Duration(app.CFG.Email.SendTimeout) * time.Second
	attemptsRemain := app.CFG.Email.RetryAttempts

	for err := e.email.Send(em, sendTimeout); attemptsRemain > 0 && time.Now().Before(retryDeadline) && err != nil; err = e.email.Send(em, sendTimeout) {
		attemptsRemain--
		log.Error("cannot send email attemptsRemain: %d, err:%s", attemptsRemain, err.Error())
	}
}
