package kbheader

import (
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils"
)

type Request struct {
	FuncNm  string `json:"funcNm"`
	RqUID   string `json:"rqUID"`
	RqAppId string `json:"rqAppId"`
	RqDt    string `json:"rqDt"`
	CorrID  string `json:"corrID"`
}

func (r *Request) Validate(fn string) error {
	if r.FuncNm != fn {
		return app.ErrorMap.InvalidRequest.AddDebug("funcNm not match")
	}
	if r.RqUID == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("rqUID cannot be empty")
	}
	if r.RqAppId == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("rqAppId cannot be empty")
	}
	if r.RqDt == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("rqDt cannot be empty")
	} else if _, e := time.Parse(utils.RFC3339, r.RqDt); e != nil {
		return app.ErrorMap.InvalidRequest.AddDebug("requestDt wrong format", e.Error())
	}

	return nil
}
