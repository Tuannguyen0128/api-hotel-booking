package app

import (
	"fmt"
	"strings"
)

type (
	ResponseError struct {
		Code         string `mapstructure:"Code" json:"code" bson:"code"`
		ErrorMessage string `mapstructure:"Message" json:"message" bson:"message"`
		HttpCode     int    `mapstructure:"HttpCode" json:"-" bson:"-"`
		DebugMessage string `json:"-" bson:"debugMessage"`
	}

	ErrorMapTemplate struct {
		InternalServerError   ResponseError `mapstructure:"InternalServerError"`
		InvalidRequest        ResponseError `mapstructure:"InvalidRequest"`
		InvalidStateError     ResponseError `mapstructure:"InvalidStateError"`
		ProgrammaticError     ResponseError `mapstructure:"ProgrammaticError"`
		FeatureDisabled       ResponseError `mapstructure:"FeatureDisabled"`
		SessionTimeout        ResponseError `mapstructure:"SessionTimeout"`
		PermissionDenied      ResponseError `mapstructure:"PermissionDenied"`
		LoginFailed           ResponseError `mapstructure:"LoginFailed"`
		PasswordNotMatch      ResponseError `mapstructure:"PasswordNotMatch"`
		PasswordValidateError ResponseError `mapstructure:"PasswordValidateError"`
	}
)

var (
	ErrorMap = &ErrorMapTemplate{}
)

func (re ResponseError) From(e error) ResponseError {
	if err, ok := e.(ResponseError); ok {
		re.DebugMessage = "extended from err:[" + err.Desc(true) + "]"
	} else {
		re.DebugMessage = "extended from err:[" + e.Error() + "]"
	}

	return re
}

func (re ResponseError) Error() string {
	return fmt.Sprintf("%s: %s", re.Code, re.ErrorMessage)
}

func (re ResponseError) SLog() string {
	return fmt.Sprintf("%s: %s [%03d][%s]", re.Code, re.ErrorMessage, re.HttpCode, re.DebugMessage)
}

func (re ResponseError) IsMatch(err error) bool {
	if err == nil {
		return false
	}

	if ree, ok := err.(ResponseError); !ok {
		return false
	} else {
		return re.Code == ree.Code
	}
}

func GetResponseError(i interface{}, isDebug bool) ResponseError {
	if err, ok := i.(ResponseError); ok {
		return err
	}

	if isDebug {
		e := ErrorMap.InternalServerError
		e.DebugMessage = fmt.Sprintf("[%#v]", i)
		return e
	}
	return ErrorMap.InternalServerError
}

func (re ResponseError) AddDebug(infos ...string) ResponseError {
	re.DebugMessage = strings.Join(append(infos, re.DebugMessage), "; ")
	return re
}

func (re ResponseError) CodeName() string {
	return re.Code
}

func (re ResponseError) Desc(isDebug bool) string {
	if isDebug {
		return re.ErrorMessage + "[" + re.DebugMessage + "]"
	}
	return re.ErrorMessage
}
