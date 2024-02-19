package pheader

import (
	"fmt"
	"math/rand"
	"time"

	"api-hotel-booking/internal/app/utils"
	"github.com/google/uuid"
)

const (
	SuccessCode = "00"
	FailCode    = "10"
)

type PartnerResponse interface {
	FormRequestHeader(reqH Request, err Err)
}

type Response struct {
	PartnerId  string `json:"partnerId"`
	RequestId  string `json:"requestId"`
	RequestDt  string `json:"requestDt"`
	ResponseId string `json:"responseId"`
	ResponseDt string `json:"responseDt"`
	StatusCode string `json:"statusCode"`
	Error      *err   `json:"error"`
}

type ResponseError struct {
	Header Response `json:"header"`
}

func GetResponseHeader(reqH Request, err Err) Response {
	return Response{
		PartnerId:  reqH.PartnerId,
		RequestId:  reqH.RequestId,
		RequestDt:  reqH.RequestDt,
		ResponseId: getResponseId(),
		ResponseDt: time.Now().Format(utils.RFC3339),
		StatusCode: getSuccessCode(err == nil),
		Error:      getErr(err),
	}
}

func GetResponseError(reqH Request, err Err) ResponseError {
	return ResponseError{Header: GetResponseHeader(reqH, err)}
}

func getResponseId() string {
	if resId, e := uuid.NewUUID(); e != nil {
		log.Warn(fmt.Sprintf("this error on create new uuid should be impossible but:%s", e.Error()))
		return fmt.Sprintf("%020d", rand.Int63())
	} else {
		return resId.String()
	}
}

func getSuccessCode(isSuccess bool) string {
	if isSuccess {
		return SuccessCode
	}

	return FailCode
}
