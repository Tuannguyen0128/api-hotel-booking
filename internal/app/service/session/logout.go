package session

import (
	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils/header/pheader"
)

type LogoutRequest struct {
	Header pheader.Request `json:"header"`
	Token  string          `json:"token"`
}

func (r *LogoutRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}
	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	return nil
}

func (s *service) Logout(req LogoutRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return err
	}

	if err := s.SessionDb.Delete(token.GetSessionId()); err != nil {
		return app.ErrorMap.InternalServerError.From(err).AddDebug("cannot logout")
	}

	return nil
}
