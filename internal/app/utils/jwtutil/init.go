package jwtutil

import "errors"

type util struct {
}

func NewUtil() *util {
	return &util{}
}

var (
	ParseTimeError = errors.New("error when parsing expire time")
)
