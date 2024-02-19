package pheader

import "api-hotel-booking/internal/app/thirdparty/logger"

var isDebug = false

var log = logger.WithModule("pheader")

func Init(setDebug bool) {
	isDebug = setDebug
}
