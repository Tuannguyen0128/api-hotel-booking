package user

import (
	"context"
	"fmt"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/thirdparty/hash"
	"api-hotel-booking/internal/app/utils/header/pheader"
)

type ChangePassRequest struct {
	Header      pheader.Request `json:"header"`
	Token       string          `json:"token"`
	OldPassword string          `json:"oldPassword"`
	NewPassword string          `json:"newPassword"`
}

type ChangePassResponse struct {
	Header pheader.Response `json:"header"`
}

func (r *ChangePassRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	if r.OldPassword == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("old password cannot be empty")
	}

	if r.NewPassword == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("new password cannot be empty")
	}

	return nil
}

func (s *service) ChangePassword(ctx context.Context, req ChangePassRequest) (ChangePassResponse, error) {
	if err := req.Validate(); err != nil {
		return ChangePassResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return ChangePassResponse{}, err
	}

	userId := token.GetUserId()
	profile, err := s.UserProfileDb.GetById(ctx, userId)
	if err != nil {
		return ChangePassResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this user profile")
	}

	if isMatch, err := hash.ComparePassword(profile.GetActivePassword(), req.OldPassword); err != nil {
		return ChangePassResponse{}, app.ErrorMap.LoginFailed.AddDebug("error when validate password")
	} else if isMatch == false {
		if i, err := s.UserProfileDb.WrongPasswordCounterIncrease(profile.Id, s.WrongPasswordLockLimit, s.WrongPasswordLockToStatus); err != nil {
			return ChangePassResponse{}, app.ErrorMap.PasswordNotMatch.From(err).AddDebug(fmt.Sprintf("the password is wrong %d times", i))
		} else {
			return ChangePassResponse{}, app.ErrorMap.PasswordNotMatch.AddDebug(fmt.Sprintf("the password is wrong %d times", i))
		}
	} else if profile.WrongPassword != 0 {
		if err := s.UserProfileDb.WrongPasswordCounterReset(profile.Id); err != nil {
			return ChangePassResponse{}, app.ErrorMap.InternalServerError.From(err).AddDebug("cannot reset wrong password count")
		}
	}

	if err := s.doChangePassword(ctx, req.NewPassword, profile, err, userId); err != nil {
		return ChangePassResponse{}, err
	}

	return ChangePassResponse{}, nil
}

func (s *service) doChangePassword(ctx context.Context, newPassword string, profile persistence.UserProfile, err error, userId string) error {
	for _, oldPassword := range profile.OldPasswords {
		if isMatch, err := hash.ComparePassword(oldPassword, newPassword); err != nil {
			return app.ErrorMap.InternalServerError.AddDebug("error when validate password")
		} else if isMatch == true {
			return app.ErrorMap.PasswordValidateError.AddDebug(
				"new password cannot be the same as last 5 passwords")
		}
	}

	update := persistence.EditUserProfile{}

	if profile.TempPassword != "" {
		if isMatch, err := hash.ComparePassword(profile.TempPassword, newPassword); err != nil {
			return app.ErrorMap.InternalServerError.AddDebug("error when validate password")
		} else if isMatch == true {
			return app.ErrorMap.InvalidRequest.AddDebug(
				"new password cannot be the same as temporary password")
		}
		removePassword := ""
		update.TempPassword = &removePassword
	}

	newHashedPassword, err := hash.NewHashPassword(newPassword)
	if err != nil {
		return app.ErrorMap.InternalServerError.AddDebug("cannot hash new password")
	}

	if len(profile.OldPasswords) == 5 {
		profile.OldPasswords = profile.OldPasswords[1:]
	}
	update.OldPasswords = append(profile.OldPasswords, newHashedPassword)
	update.Password = newHashedPassword

	err = s.UserProfileDb.Update(ctx, userId, update)
	if err != nil {
		return app.ErrorMap.InternalServerError.AddDebug(
			"error when update user profile")
	}

	return nil
}
