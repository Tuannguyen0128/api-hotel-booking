package user

import (
	"context"
	"fmt"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"api-hotel-booking/internal/app/utils/jwtutil"
	"api-hotel-booking/internal/app/utils/rbac"
)

type GetRequest struct {
	Header pheader.Request `json:"header"`
	Token  string          `json:"token"`
	Id     string          `json:"id,omitempty"`
}

type GetResponse struct {
	Header pheader.Response `json:"header"`
	userInfo
}

func (r *GetRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	return nil
}

func (s *service) GetUser(ctx context.Context, req GetRequest) (GetResponse, error) {
	if err := req.Validate(); err != nil {
		return GetResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return GetResponse{}, err
	}

	var profile persistence.UserProfile
	if req.Id == "" {
		profile, err = s.UserProfileMemorizer.Get(ctx, token.GetUserId())
		if err != nil {
			return GetResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this profile")
		}

	} else {
		profile, err = s.UserProfileMemorizer.Get(ctx, req.Id)
		if err != nil {
			return GetResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this profile")
		}

		if err := rbac.IsHavePermission(token, rbac.UserView, profile.CompanyId); err != nil {
			return GetResponse{}, app.ErrorMap.PermissionDenied.From(err).AddDebug("do not have permission to view this profile")
		}
	}

	if err := checkCanViewUser(profile.Id, profile.Role, profile.CompanyId, token); err != nil {
		return GetResponse{}, err
	}

	var companyName string
	//if profile.CompanyId != "" {
	//	if companyProfile, err := s.CompanyProfileMemorizer.Get(ctx, profile.CompanyId); err != nil {
	//		if err.Error() != persistence.NotFoundError.Error() {
	//			return GetResponse{}, app.ErrorMap.InternalServerError.From(err).AddDebug("cannot get user profile dur to cannot get company info")
	//		}
	//	} else {
	//		companyName = companyProfile.CompanyName
	//	}
	//}

	var createTime string
	var editTime string
	if !profile.CreateDt.IsZero() {
		createTime = profile.CreateDt.Format(utils.RFC3339)
	}
	if !profile.LastEditDt.IsZero() {
		editTime = profile.LastEditDt.Format(utils.RFC3339)
	}

	return GetResponse{
		Header: pheader.Response{},
		userInfo: userInfo{
			Id:          profile.Id,
			Email:       profile.Email,
			Name:        profile.Name,
			Surname:     profile.Surname,
			CompanyId:   profile.CompanyId,
			Role:        profile.Role,
			Status:      profile.Status,
			CreateBy:    profile.CreateBy,
			CreateDt:    createTime,
			LastEditDt:  editTime,
			CompanyName: companyName,
		},
	}, nil
}

func checkCanViewUser(targetUserId, targetUserRoleName, targetCompanyId string, token jwtutil.JwtToken) error {
	if token.GetUserId() == targetUserId {
		return nil
	}

	if rbac.IsCompanyAdmin(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.AdminUserView, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to edit profile", err.Error())
		}
	} else if rbac.IsViewerCheckerLabeler(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.UserView, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to edit profile", err.Error())
		}
	} else {
		return app.ErrorMap.InvalidRequest.AddDebug(fmt.Sprintf("edit profile role no-supported: %s", targetUserRoleName))
	}

	return nil
}
