package session

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/pheader"
)

type ExtendRequest struct {
	Header             pheader.Request `json:"header"`
	Token              string          `json:"token"`
	EncryptedSessionId string          `json:"encryptedSessionId"`
}
type ExtendResponse struct {
	Header pheader.Response `json:"header"`
	Token  string           `json:"token"`
}

func (r *ExtendRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}
	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}
	if r.EncryptedSessionId == "" && app.CFG.Service.Session.IsRequireEncryption {
		return app.ErrorMap.InvalidRequest.AddDebug("encrypted session id cannot be empty")
	}

	return nil
}

func (s *service) ExtendSession(ctx context.Context, req ExtendRequest) (ExtendResponse, error) {
	if err := req.Validate(); err != nil {
		return ExtendResponse{}, err
	}

	blankClaims := s.JwtUtil.NewBlankClaims()
	err := blankClaims.ParseToken(req.Token)
	if err != nil {
		return ExtendResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot parse the jwt token")
	}
	userId := blankClaims.GetUserId()
	companyId := blankClaims.GetCompanyId()
	sessionId := blankClaims.GetSessionId()
	email := blankClaims.GetEmail()
	name := blankClaims.GetName()
	surname := blankClaims.GetSurname()
	role := blankClaims.GetRole()

	session, err := s.SessionDb.Get(ctx, sessionId)
	if err != nil {
		return ExtendResponse{}, app.ErrorMap.InvalidRequest.AddDebug(err.Error(), "cannot find session info")
	}

	if app.CFG.Service.Session.IsRequireEncryption {
		decodedKey, err := base64.StdEncoding.DecodeString(session.Key)
		if err != nil {
			return ExtendResponse{}, app.ErrorMap.InternalServerError.AddDebug(err.Error(), "cannot decode session key")
		}

		encrypt := s.EncryptUtil.NewEncryptor(string(decodedKey))
		decryptedSessionId, err := encrypt.Decrypt(req.EncryptedSessionId)

		if err != nil {
			return ExtendResponse{}, app.ErrorMap.InvalidRequest.AddDebug(err.Error(), "cannot decrypt session id")
		}
		if strings.Compare(sessionId, decryptedSessionId) != 0 {
			return ExtendResponse{}, app.ErrorMap.InvalidRequest.AddDebug("session id is not matched")
		}
	}

	timeNow := time.Now()
	if timeNow.Equal(session.Timeout) || timeNow.After(session.Timeout) {
		return ExtendResponse{}, app.ErrorMap.SessionTimeout.AddDebug("this session is expired")
	}

	timeout := timeNow.Add(time.Duration(s.Timeout) * time.Second)
	err = s.SessionDb.UpdateTimeout(ctx, sessionId, timeout)
	if err != nil {
		return ExtendResponse{}, app.ErrorMap.InternalServerError.AddDebug(err.Error(), "error when update session timeout")
	}

	claims := s.JwtUtil.NewClaims(sessionId, userId, email, name, surname, companyId, role, timeout.Format(utils.RFC3339))
	token, err := claims.GenerateToken()
	if err != nil {
		return ExtendResponse{}, app.ErrorMap.InternalServerError.AddDebug(err.Error(), "error when generating token")
	}

	return ExtendResponse{
		Header: pheader.Response{},
		Token:  token,
	}, nil
}
