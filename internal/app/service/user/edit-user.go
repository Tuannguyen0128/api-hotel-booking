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

type EditRequest struct {
	Header    pheader.Request `json:"header"`
	Token     string          `json:"token"`
	UserId    string          `json:"userId"`
	CompanyId string          `json:"companyId"`
	Name      *string         `json:"name,omitempty"`
	Surname   *string         `json:"surname,omitempty"`
	Role      *string         `json:"role,omitempty"`
	Status    *string         `json:"status,omitempty"`
}

type EditResponse struct {
	Header pheader.Response `json:"header"`
}

func (r *EditRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	if r.UserId == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("userId cannot be empty")
	}

	if r.Name == nil && r.Surname == nil && r.Role == nil && r.Status == nil {
		return app.ErrorMap.InvalidRequest.AddDebug("must update at least one parameter")
	}

	if r.Name != nil && r.Surname != nil {
		if *r.Name == "" && *r.Surname == "" {
			return app.ErrorMap.InvalidRequest.AddDebug("name and surname cannot be both empty")
		}
	}

	if r.Status != nil {
		status := *r.Status
		if status != StatusActive && status != StatusInactive && status != StatusDelete {
			return app.ErrorMap.InvalidRequest.AddDebug("invalid status")
		}
	}

	if r.Role != nil {
		_, isIn := rbac.Roles[*r.Role]
		if isIn == false {
			return app.ErrorMap.InvalidRequest.AddDebug("user role is wrong")
		}
	}

	return nil
}

func (s *service) EditUser(ctx context.Context, req EditRequest) (EditResponse, error) {
	if err := req.Validate(); err != nil {
		return EditResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return EditResponse{}, err
	}

	targetProfile, err := s.UserProfileMemorizer.Get(ctx, req.UserId)
	if err != nil {
		return EditResponse{}, app.ErrorMap.InvalidRequest.AddDebug("cannot find this user targetProfile")
	}

	if targetProfile.Status == StatusDelete {
		return EditResponse{}, app.ErrorMap.InternalServerError.AddDebug("cannot edit targetProfile because this profile has been deleted")
	}

	if err := checkCanEditUser(targetProfile.Role, targetProfile.CompanyId, token); err != nil {
		return EditResponse{}, err
	}

	update, err := req.getUpdateProfile(token)
	if err != nil {
		return EditResponse{}, err
	}

	err = s.UserProfileDb.UpdateInfo(ctx, update)
	if err != nil {
		return EditResponse{}, app.ErrorMap.InternalServerError.AddDebug("cannot edit the targetProfile")
	}
	return EditResponse{
		Header: pheader.Response{},
	}, nil
}

func (r *EditRequest) getUpdateProfile(token jwtutil.JwtToken) (persistence.EditUserProfileInfo, error) {
	now := time.Now()
	update := persistence.EditUserProfileInfo{
		UserId:     r.UserId,
		LastEditDt: &now,
	}

	// only KBank admin can edit other company account
	if rbac.RoleKBankAdmin.IsRole(token) {
		if r.Role != nil && *r.Role != rbac.RoleAdmin.Name() {
			return persistence.EditUserProfileInfo{}, app.ErrorMap.InvalidRequest.AddDebug("KBankAdmin cannot change user role")
		}
		update.CompanyId = r.CompanyId
		update.Role = rbac.RoleAdmin.Name()
	} else {
		update.CompanyId = token.GetCompanyId()
	}

	if r.Role != nil {
		if r.UserId != token.GetUserId() { // not update themselves
			if err := canAssignToRole(token.GetRole(), *r.Role); err != nil {
				return persistence.EditUserProfileInfo{}, err
			}
			update.NewRole = r.Role
		}
	}

	if r.Name != nil {
		update.NewName = r.Name
	}

	if r.Surname != nil {
		update.NewSurname = r.Surname
	}

	if r.Status != nil {
		update.NewStatus = r.Status
	}

	if *update.NewStatus == StatusDelete {
		update.DeletedDt = &now
	}

	return update, nil
}

func checkCanEditUser(targetUserRoleName, targetCompanyId string, token jwtutil.JwtToken) error {
	if rbac.IsCompanyAdmin(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.AdminUserEdit, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to edit profile", err.Error())
		}
	} else if rbac.IsViewerCheckerLabeler(targetUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.UserEdit, targetCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to edit profile", err.Error())
		}
	} else {
		return app.ErrorMap.InvalidRequest.AddDebug(fmt.Sprintf("edit profile role no-supported: %s", targetUserRoleName))
	}

	return nil
}

func canAssignToRole(assigner, assignee string) error {
	if rbac.IsKBankAdmin(assigner) {
		if rbac.IsCompanyAdmin(assignee) {
			return nil
		}
	} else if rbac.IsCompanyAdmin(assigner) {
		if rbac.IsViewerCheckerLabeler(assignee) {
			return nil
		}
	}

	return app.ErrorMap.InvalidRequest.AddDebug(fmt.Sprintf("%s cannot assign user to %s", assigner, assignee))
}
