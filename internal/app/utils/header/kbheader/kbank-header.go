package kbheader

import (
	"regexp"
	"strings"

	"api-hotel-booking/internal/app/thirdparty/logger"
)

var (
	log = logger.WithModule("kbheader")

	appId      = "undefined"
	appIdDigit = "0"
	appAbbr    = "undefined"
	isDebug    = false
)

const requestResponseIdDateFormat = "20060102"

func InitKBankHeader(AppId string, AppAbbr string, IsDebug bool) {
	appId = AppId
	appAbbr = AppAbbr
	appIdDigit = strings.Join(regexp.MustCompile("[0-9]+").FindAllString(AppId, -1), "")
	isDebug = IsDebug
}
