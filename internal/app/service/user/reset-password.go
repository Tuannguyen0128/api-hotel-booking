package user

import (
	"context"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils/header/pheader"
)

type ResetPassRequest struct {
	Header      pheader.Request `json:"header"`
	Token       string          `json:"token"`
	NewPassword string          `json:"newPassword"`
}

type ResetPassResponse struct {
	Header pheader.Response `json:"header"`
}

func (r *ResetPassRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	if r.NewPassword == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("new password cannot be empty")
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPassRequest) (ResetPassResponse, error) {
	if err := req.Validate(); err != nil {
		return ResetPassResponse{}, err
	}

	var userId string
	if token, err := s.JwtUtil.ParseAndVerifyToken(req.Token); err != nil {
		userId = token.GetUserId()
	} else if forgotToken, err := s.JwtUtil.ParseAndVerifyForgetPWToken(req.Token); err != nil {
		userId = forgotToken.GetUserId()
	} else {
		return ResetPassResponse{}, app.ErrorMap.InvalidRequest.AddDebug("Can't parse token.")
	}

	profile, err := s.UserProfileDb.GetById(ctx, userId)
	if err != nil {
		return ResetPassResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this user profile")
	}

	err = s.doChangePassword(ctx, req.NewPassword, profile, err, userId)
	if err != nil {
		return ResetPassResponse{}, err
	}

	return ResetPassResponse{
		Header: pheader.Response{},
	}, nil
}
