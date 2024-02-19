package logstats

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	timeFormatRFC3339Milli      = "2006-01-02T15:04:05.000Z07:00"
	fieldRequestId              = "reqUID"
	fieldCallServiceID          = "callServiceID"
	fieldPartnerId              = "partnerId"
	fieldAPIName                = "apiName"
	fieldAdditionInfo           = "additionInfo"
	fieldRemark                 = "remark"
	fieldValueAutoRemarkSuccess = "success"
	fieldValueAutoRemarkFail    = "fail"
	fieldValueUndefined         = "undefined"
)

func getRequestId(ctx *gin.Context, log logger) string {
	if rid := ctx.GetHeader(fieldRequestId); rid != "" {
		return rid
	}
	if rid := log.GetRequestID(ctx); rid != "" {
		return rid
	}

	return fmt.Sprintf("%s_ServerGen", uuid.New().String())
}

func getCallerServiceId(ctx *gin.Context) string {
	cid := ctx.GetString(fieldCallServiceID)

	if cid != "" {
		return cid
	}

	cid = ctx.GetHeader(fieldCallServiceID)

	if cid != "" {
		return cid
	}
	return fieldValueUndefined
}

func SetCallerServiceId(ctx *gin.Context, callerServiceId string) {
	ctx.Set(fieldCallServiceID, callerServiceId)
}

func getPartnerId(ctx *gin.Context) string {
	pid := ctx.GetString(fieldPartnerId)
	if pid != "" {
		return pid
	}

	return fieldValueUndefined
}

func SetPartnerId(ctx *gin.Context, partnerIdFromBody string) {
	ctx.Set(fieldPartnerId, partnerIdFromBody)
}

func getAPIName(ctx *gin.Context) string {
	api := ctx.GetString(fieldAPIName)
	if api != "" {
		return api
	}

	return fieldValueUndefined
}

func SetAPIName(ctx *gin.Context, apiName string) {
	ctx.Set(fieldAPIName, apiName)
}

func getRemark(ctx *gin.Context) string {
	remark := ctx.GetString(fieldRemark)
	if remark != "" {
		return remark
	}

	if ctx.Writer.Status()/100 == 2 {
		return fieldValueAutoRemarkSuccess
	}

	return fieldValueAutoRemarkFail
}

func SetRemark(ctx *gin.Context, remark string) {
	ctx.Set(fieldRemark, remark)
}

func getAdditionInfo(ctx *gin.Context) interface{} {
	info, _ := ctx.Get(fieldAdditionInfo)
	return info
}

func SetAdditionInfo(ctx *gin.Context, info interface{}) {
	ctx.Set(fieldAdditionInfo, info)
}
