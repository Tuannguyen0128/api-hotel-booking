package utils

import (
	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/thirdparty/logger"
	"github.com/gin-gonic/gin"
)

type (
	Closer struct {
		Name    string
		CloseFn func() error
	}
)

var (
	log     = logger.WithModule("utils")
	closers []Closer
)

const (
	TimeFormatRFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
	TimeFormatYYYYMMDD     = "20060102"
)

func RegisterCloser(name string, fn func() error) {
	closers = append(closers, Closer{name, fn})
}

func OnClose() {
	for _, c := range closers {
		log.Info("closing connection of %s", c.Name)
		if err := c.CloseFn(); err != nil {
			log.Error("error while closing %s: %v", c.Name, err)
		}
	}
}

func AbortWithError(ctx *gin.Context, err error) {
	if responseError, ok := err.(app.ResponseError); ok {
		ctx.AbortWithStatusJSON(responseError.HttpCode, gin.H{"error": responseError})
		return
	}

	log.KError(ctx, "got undefined error in response: %v", err)
	AbortWithError(ctx, app.ErrorMap.InternalServerError)
}
