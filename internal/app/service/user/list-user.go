package user

import (
	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"api-hotel-booking/internal/app/utils/jwtutil"
	"api-hotel-booking/internal/app/utils/rbac"
	"context"
)

type ListRequest struct {
	Header pheader.Request `json:"header"`
	Token  string          `json:"token"`
	Filter Filter          `json:"filter,omitempty"`
}

type Filter struct {
	Status string `json:"status"`
}

func (f Filter) getDBFilter() persistence.UserFilter {
	return persistence.UserFilter{
		Status: f.Status,
	}
}

type ListResponse struct {
	Header pheader.Response `json:"header"`
	Users  []userInfo       `json:"users"`
}

func (r *ListRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	return nil
}

func (s *service) ListUser(ctx context.Context, req ListRequest) (ListResponse, error) {
	if err := req.Validate(); err != nil {
		return ListResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return ListResponse{}, err
	}

	var profiles []persistence.UserProfile
	if rbac.Roles[token.GetRole()].Scope() == rbac.Global {
		profiles, err = s.listAllUser(ctx, req.Filter, token)
		if err != nil {
			return ListResponse{}, err
		}
	} else {
		profiles, err = s.listUserInCompany(ctx, token.GetCompanyId(), req.Filter, token)
		if err != nil {
			return ListResponse{}, err
		}
	}

	var users []userInfo
	for _, profile := range profiles {
		var companyName string
		//if profile.CompanyId != "" {
		//	if companyProfile, err := s.CompanyProfileMemorizer.Get(ctx, profile.CompanyId); err != nil {
		//		if err.Error() != persistence.NotFoundError.Error() {
		//			return ListResponse{}, app.ErrorMap.InternalServerError.From(err).AddDebug("cannot list user profile due to cannot get company info")
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

		users = append(users, userInfo{
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
		})
	}

	return ListResponse{
		Users: users,
	}, nil
}

func (s *service) listUserInCompany(ctx context.Context, companyId string, filter Filter, token jwtutil.JwtToken) ([]persistence.UserProfile, error) {
	if err := rbac.IsHavePermission(token, rbac.UserList, companyId); err != nil {
		return []persistence.UserProfile{}, app.ErrorMap.PermissionDenied.From(err).AddDebug("do not have permission to list profile in this company")
	}

	if profiles, err := s.UserProfileDb.GetByCompanyId(ctx, companyId, filter.getDBFilter()); err != nil {
		if err != persistence.NotFoundError {
			return []persistence.UserProfile{}, app.ErrorMap.InternalServerError.AddDebug("error when query data")
		}
		return []persistence.UserProfile{}, nil
	} else {
		return profiles, nil
	}
}

func (s *service) listAllUser(ctx context.Context, filter Filter, token jwtutil.JwtToken) ([]persistence.UserProfile, error) {
	if err := rbac.IsHavePermission(token, rbac.AdminUserList, "*"); err != nil {
		return []persistence.UserProfile{}, app.ErrorMap.PermissionDenied.From(err).AddDebug("do not have permission to list all profile")
	}

	if profiles, err := s.UserProfileDb.GetAll(ctx, filter.getDBFilter()); err != nil {
		if err != persistence.NotFoundError {
			return []persistence.UserProfile{}, app.ErrorMap.InternalServerError.AddDebug("error when query data")
		}
		return []persistence.UserProfile{}, nil
	} else {
		return profiles, nil
	}
}
