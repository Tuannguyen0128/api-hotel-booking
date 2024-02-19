package user

import (
	"context"
	"fmt"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/thirdparty/hash"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"api-hotel-booking/internal/app/utils/jwtutil"
	"api-hotel-booking/internal/app/utils/rbac"
)

type CreateRequest struct {
	Header    pheader.Request `json:"header"`
	Token     string          `json:"token"`
	Email     string          `json:"email"`
	CompanyId string          `json:"companyId,omitempty"`
	Name      string          `json:"name,omitempty"`
	Surname   string          `json:"surname,omitempty"`
	Role      string          `json:"role"`
	Status    string          `json:"status"`
}

type CreateResponse struct {
	Header   pheader.Response `json:"header"`
	UserId   string           `json:"userId"`
	Password string           `json:"password"`
}

func (r *CreateRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}

	if r.Token == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("token cannot be empty")
	}

	if r.Email == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("email cannot be empty")
	}

	if r.Role == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("user role cannot be empty")
	}

	_, isIn := rbac.Roles[r.Role]
	if isIn == false {
		return app.ErrorMap.InvalidRequest.AddDebug("user role is wrong")
	}

	if r.Status == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("status cannot be empty")
	} else {
		if r.Status != StatusActive && r.Status != StatusInactive && r.Status != StatusDelete {
			return app.ErrorMap.InvalidRequest.AddDebug("invalid status")
		}
	}

	return nil
}

func (s *service) CreateUser(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return CreateResponse{}, err
	}

	token, err := s.JwtUtil.ParseAndVerifyToken(req.Token)
	if err != nil {
		return CreateResponse{}, err
	}

	var companyId string
	if rbac.RoleKBankAdmin.IsRole(token) {
		companyId = req.CompanyId
	} else {
		companyId = token.GetCompanyId()
	}
	if err := checkCanCreateUser(req.Role, companyId, token); err != nil {
		return CreateResponse{}, err
	}

	tempPassword := s.PasswordUtil.GeneratePassword()
	hashedPassword, err := hash.NewHashPassword(tempPassword)
	if err != nil {
		return CreateResponse{}, app.ErrorMap.InternalServerError.AddDebug("cannot hash the password")
	}

	// check company exist
	//if _, err := s.CompanyProfileMemorizer.Get(ctx, companyId); err != nil {
	//	if err == persistence.NotFoundError {
	//		return CreateResponse{}, app.ErrorMap.InvalidRequest.AddDebug("new user company not found")
	//	}
	//	return CreateResponse{}, app.ErrorMap.InternalServerError.AddDebug("new user error while getting company")
	//}

	newUser := persistence.UserProfile{
		Email:        req.Email,
		Password:     "",
		Name:         req.Name,
		Surname:      req.Surname,
		CompanyId:    companyId,
		Role:         req.Role,
		Status:       req.Status,
		OldPasswords: []string{},
		TempPassword: hashedPassword,
		CreateBy:     token.GetFullName(),
		CreateDt:     time.Now(),
		LastEditDt:   time.Now(),
	}

	userId, err := s.UserProfileDb.Insert(ctx, newUser)
	if err != nil {
		if err == persistence.DuplicateEmailError {
			return CreateResponse{}, app.ErrorMap.InvalidRequest.AddDebug(err.Error())
		}
		return CreateResponse{}, app.ErrorMap.InternalServerError.From(err).AddDebug("error when creating profile")
	}

	err = s.EmailUtil.SendCreateUserEmail(req.Email, req.Name, req.Surname, tempPassword)
	if err != nil {
		return CreateResponse{}, app.ErrorMap.InternalServerError.AddDebug("cannot send email")
	}

	return CreateResponse{
		Header:   pheader.Response{},
		UserId:   userId,
		Password: tempPassword,
	}, nil
}

func checkCanCreateUser(newUserRoleName, newUserCompanyId string, token jwtutil.JwtToken) error {
	if rbac.IsCompanyAdmin(newUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.AdminUserAdd, newUserCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to create profile", err.Error())
		}
	} else if rbac.IsViewerCheckerLabeler(newUserRoleName) {
		if err := rbac.IsHavePermission(token, rbac.UserAdd, newUserCompanyId); err != nil {
			return app.ErrorMap.InvalidRequest.AddDebug("do not have permission to create profile", err.Error())
		}
	} else {
		return app.ErrorMap.InvalidRequest.AddDebug(fmt.Sprintf("new user role not supported: %s", newUserRoleName))
	}

	return nil
}
