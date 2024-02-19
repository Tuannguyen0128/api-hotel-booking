package user

import (
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/utils"
	"context"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils/header/pheader"
)

type ForgetPassRequest struct {
	Header pheader.Request `json:"header"`
	Email  string          `json:"email"`
}

func (r *ForgetPassRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Email == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("email cannot be empty")
	}

	return nil
}

func (s *service) ForgetPassword(ctx context.Context, req ForgetPassRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	profile, err := s.UserProfileDb.GetByEmail(ctx, req.Email)
	if err != nil {
		return app.ErrorMap.InvalidRequest.AddDebug("cannot find this user profile")
	}

	if profile.ForgotPasswordEmailTimeOut.Before(time.Now()) {
		timeout := time.Now().Add(time.Duration(s.ForgetPWConfig.Timeout) * time.Second)
		//update forgot password timeout
		update := persistence.EditUserProfile{
			ForgotPasswordEmailTimeOut: timeout,
		}
		s.UserProfileDb.Update(ctx, profile.Id, update)

		token := s.JwtUtil.NewForgetPWClaims(profile.Id, profile.Email, timeout.Format(utils.RFC3339))
		tokenString, err := token.GenerateToken()
		if err != nil {
			return app.ErrorMap.InternalServerError.AddDebug("error when generating token")
		}
		err = s.EmailUtil.SendForgetPasswordEmail(profile.Email, profile.Name, profile.Surname, tokenString)
		if err != nil {
			return app.ErrorMap.InternalServerError.AddDebug("cannot send email")
		}
	}

	return nil
}
