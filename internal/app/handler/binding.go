package handler

import (
	"github.com/gin-gonic/gin"
)

// do not wrap anything that do not help shorten or simplify anythings that has been already optimized

func BindRequestHeader(c *gin.Context, request interface{}) error {
	if err := c.BindHeader(request); err != nil {
		log.KError(c, "error binding request header: %v", err)
		return err
	}
	log.DDebug(c, request, "request header")
	return nil
}

func BindRequestBody(c *gin.Context, request interface{}) error {
	if err := c.BindJSON(request); err != nil {
		log.KError(c, "error binding request body: %v", err)
		return err
	}
	log.DDebug(c, request, "request body")
	return nil
}

func BindRequestParams(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindUri(request); err != nil {
		log.KError(c, "error binding request params: %v", err)
		return err
	}
	log.DDebug(c, request, "request params")
	return nil
}

func BindRequestForm(c *gin.Context, request interface{}) error {
	if err := c.ShouldBind(request); err != nil {
		log.KError(c, "error binding request form: %v", err)
		return err
	}
	log.DDebug(c, request, "request form")
	return nil
}
