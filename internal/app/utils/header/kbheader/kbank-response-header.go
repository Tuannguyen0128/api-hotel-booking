package kbheader

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

type Response struct {
	FuncNm     string `json:"funcNm"`
	RqUID      string `json:"rqUID"`
	RsAppId    string `json:"rsAppId"`
	RsUID      string `json:"rsUID"`
	RsDt       string `json:"rsDt"`
	CorrID     string `json:"corrID"`
	StatusCode string `json:"statusCode"`
	Errors     []KErr `json:"errors"`
}

type ResponseError struct {
	Header Response `json:"header"`
}

func GetResponseHeader(reqH Request, err Err) Response {
	var errors []KErr
	if err != nil {
		errors = []KErr{getKErr(err)}
	}
	return Response{
		FuncNm:     reqH.FuncNm,
		RqUID:      reqH.RqUID,
		RsAppId:    appId,
		RsUID:      getKResponseId(),
		RsDt:       time.Now().Format(utils.RFC3339),
		CorrID:     getCorrId(reqH.CorrID),
		StatusCode: getSuccessCode(err == nil),
		Errors:     errors,
	}
}

func GetResponseError(reqH Request, err Err) ResponseError {
	return ResponseError{Header: GetResponseHeader(reqH, err)}
}

func getCorrId(reqCorrId string) string {
	if reqCorrId != "" {
		return reqCorrId
	}

	if corrId, e := uuid.NewUUID(); e != nil {
		log.Warn(fmt.Sprintf("this error on create new uuid should be impossible but:%s", e.Error()))
		return fmt.Sprintf("%020d", rand.Int63())
	} else {
		return corrId.String()
	}
}
func getKResponseId() string {
	if resId, e := uuid.NewUUID(); e != nil {
		log.Warn(fmt.Sprintf("this error on create new uuid should be impossible but:%s", e.Error()))
		return fmt.Sprintf("%s_%s_%020d", appIdDigit, time.Now().Format(requestResponseIdDateFormat), rand.Int63())
	} else {
		return fmt.Sprintf("%s_%s_%s", appIdDigit, time.Now().Format(requestResponseIdDateFormat), resId.String())
	}
}

func getSuccessCode(isSuccess bool) string {
	if isSuccess {
		return SuccessCode
	}

	return FailCode
}
