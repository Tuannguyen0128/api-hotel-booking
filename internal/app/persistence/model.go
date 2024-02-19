package persistence

import (
	"time"
)

type (
	Session struct {
		SessionId string    `bson:"sessionId"`
		UserId    string    `bson:"userId"`
		Key       string    `bson:"key"`
		Timeout   time.Time `bson:"timeout"`
	}
)

func (u UserProfile) GetActivePassword() string {
	if u.TempPassword != "" {
		return u.TempPassword
	} else {
		return u.Password
	}
}

func (u UserProfile) HaveAnyPassword() bool {
	if u.TempPassword == "" && u.Password == "" {
		return false
	}
	return true
}
