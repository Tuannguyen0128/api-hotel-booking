package persistence

import "errors"

var (
	NotFoundError           = errors.New("mongo: no documents in result")
	DuplicateEmailError     = errors.New("already have a profile with the same email")
	DuplicateUniqueTagError = errors.New("duplicate unique tag")
	DatabaseConnectError    = errors.New("cannot connect to the database")
)
