package pheader

import (
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils"
)

type PartnerRequest interface {
	GetPartnerId() string
	GetPartnerSecret() string
}

type Request struct {
	PartnerId     string `json:"partnerId"`
	PartnerSecret string `json:"partnerSecret"`
	RequestId     string `json:"requestId"`
	RequestDt     string `json:"requestDt"`
}

func (r *Request) GetPartnerId() string {
	return r.PartnerId
}

func (r *Request) GetPartnerSecret() string {
	return r.PartnerSecret
}

func (r *Request) Validate(omitSecret bool) error {
	if r.PartnerId == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("partnerId cannot be empty")
	}
	if !omitSecret && r.PartnerSecret == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("partnerSecret cannot be empty")
	}
	if r.RequestId == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("requestId cannot be empty")
	}
	if r.RequestDt == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("requestDt cannot be empty")
	} else if _, e := time.Parse(utils.RFC3339, r.RequestDt); e != nil {
		return app.ErrorMap.InvalidRequest.AddDebug("requestDt wrong format", e.Error())
	}

	return nil
}
