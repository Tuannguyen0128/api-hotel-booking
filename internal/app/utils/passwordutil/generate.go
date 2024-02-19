package passwordutil

import (
	"math/rand"
	"time"
)

const PasswordLength = 32

type PasswordUtil interface {
	GeneratePassword() string
}

func (u *util) GeneratePassword() string {
	randomPass := ""
	var chars []rune
	for i := 'A'; i <= 'H'; i++ {
		chars = append(chars, i)
	}
	chars = append(chars, []rune{'@', '*', '#', '$'}...)
	for i := '0'; i <= '9'; i++ {
		chars = append(chars, i)
	}
	for i := 'a'; i <= 'h'; i++ {
		chars = append(chars, i)
	}
	for i := 'I'; i <= 'Z'; i++ {
		chars = append(chars, i)
	}
	for i := 'i'; i <= 'z'; i++ {
		chars = append(chars, i)
	}

	for i := 0; i < PasswordLength; i++ {
		rand.Seed(time.Now().UnixNano())
		randInt := rand.Intn(PasswordLength)
		randomPass += string(chars[randInt])
	}
	return randomPass
}
