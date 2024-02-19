package user

import (
	"context"
	"fmt"
	"time"

	"api-hotel-booking/internal/app/persistence"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"api-hotel-booking/internal/app/utils/jwtutil"
	"api-hotel-booking/internal/app/utils/rbac"
)

type DeleteRequest struct {
	Header pheader.Request `json:"header"`
	Token  string          `json:"token"`
	UserId string          `json:"userId,omitempty"`
}

type DeleteResponse struct {
	Header pheader.Response `json:"header"`
}

func (r *DeleteRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	return nil
}

func (s *service) DeleteUser(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	if err := req.Validate(); err != nil {
		return DeleteResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return DeleteResponse{}, err
	}

	var targetId string
	if req.UserId != "" {
		targetId = req.UserId
	} else {
		targetId = token.GetUserId()
	}

	targetProfile, err := s.UserProfileMemorizer.Get(ctx, targetId)
	if err != nil {
		return DeleteResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this user profile")
	}

	if err := checkCanDeleteUser(targetProfile.Role, targetProfile.CompanyId, token); err != nil {
		return DeleteResponse{}, err
	}
	update := persistence.EditUserProfile{
		Status:    StatusDelete,
		DeletedDt: time.Now(),
	}

	err = s.UserProfileDb.Update(ctx, targetProfile.Id, update)
	if err != nil {
		return DeleteResponse{}, app.ErrorMap.InternalServerError.AddDebug("error when update profile")
	}
	return DeleteResponse{
		Header: pheader.Response{},
	}, nil
}

func checkCanDeleteUser(targetUserRoleName, targetCompanyId string, token jwtutil.JwtToken) error {
	if rbac.IsCompanyAdmin(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.AdminUserDelete, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to delete profile", err.Error())
		}
	} else if rbac.IsViewerCheckerLabeler(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.UserDelete, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to delete profile", err.Error())
		}
	} else {
		return app.ErrorMap.InvalidRequest.AddDebug(fmt.Sprintf("delete profile role no-supported: %s", targetUserRoleName))
	}

	return nil
}
