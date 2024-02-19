package user

import (
	"context"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/thirdparty/hash"
	"api-hotel-booking/internal/app/thirdparty/logger"
	"api-hotel-booking/internal/app/utils/jwtutil"
	"api-hotel-booking/internal/app/utils/passwordutil"
	"api-hotel-booking/internal/app/utils/rbac"
	"github.com/google/uuid"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, req CreateRequest) (CreateResponse, error)
	EditUser(ctx context.Context, req EditRequest) (EditResponse, error)
	DeleteUser(ctx context.Context, req DeleteRequest) (DeleteResponse, error)
	GetUser(ctx context.Context, req GetRequest) (GetResponse, error)
	ChangePassword(ctx context.Context, req ChangePassRequest) (ChangePassResponse, error)
	ListUser(ctx context.Context, req ListRequest) (ListResponse, error)
	ForgetPassword(ctx context.Context, req ForgetPassRequest) error
	ResetPassword(ctx context.Context, req ResetPassRequest) (ResetPassResponse, error)
}

type service struct {
	UserProfileDb
	app.SessionConfig
	app.ForgetPWConfig
	JwtUtil
	PasswordUtil
	EmailUtil
	UserProfileMemorizer userGetter
}

type UserProfileDb interface {
	Insert(ctx context.Context, profile persistence.UserProfile) (string, error)
	GetByEmail(ctx context.Context, email string) (persistence.UserProfile, error)
	GetById(ctx context.Context, userId string) (persistence.UserProfile, error)
	GetByCompanyId(ctx context.Context, companyId string, filter persistence.UserFilter) ([]persistence.UserProfile, error)
	GetAll(ctx context.Context, filter persistence.UserFilter) ([]persistence.UserProfile, error)
	Update(ctx context.Context, userId string, profile persistence.EditUserProfile) error
	UpdateInfo(ctx context.Context, updateReq persistence.EditUserProfileInfo) error
	WrongPasswordCounterReset(userId string) error
	WrongPasswordCounterIncrease(userId string, lockIfOver int, changeStatusTo string) (int, error)
}

type userGetter interface {
	Get(ctx context.Context, userId string) (persistence.UserProfile, error)
	Delete(userId string)
}

type JwtUtil interface {
	NewClaims(sessionId, userId, email, name, surname, companyId, role, expire string) jwtutil.JwtToken
	ParseAndVerifyToken(token string) (jwtutil.JwtToken, error)
	NewForgetPWClaims(userId, email, expire string) jwtutil.ForgetPWJwtToken
	ParseAndVerifyForgetPWToken(token string) (jwtutil.ForgetPWJwtToken, error)
}

type PasswordUtil interface {
	GeneratePassword() string
}

type EmailUtil interface {
	SendForgetPasswordEmail(email, name, surname, token string) error
	SendCreateUserEmail(email, name, surname, tempPassword string) error
}

func NewServiceUser(sessionConfig app.SessionConfig, forgetPWConfig app.ForgetPWConfig, userGetter userGetter, userProfileDb UserProfileDb, emailUtil EmailUtil) ServiceUser {
	if admin, err := userProfileDb.GetByEmail(context.TODO(), "super@admin.com"); err != nil {
		hashedPassword, _ := hash.NewHashPassword("super")
		if err != persistence.NotFoundError {
			panic(err)
		} else {
			user := persistence.UserProfile{
				Id:           "user_" + uuid.NewString(),
				Email:        "super@admin.com",
				Password:     hashedPassword,
				Name:         "super",
				Surname:      "admin",
				CompanyId:    "kbtg",
				Role:         rbac.RoleKBankAdmin.Name(),
				Status:       StatusActive,
				OldPasswords: nil,
				TempPassword: "",
				CreateBy:     "Test",
				CreateDt:     time.Now(),
			}
			if _, err := userProfileDb.Insert(context.TODO(), user); err != nil {
				panic(err)
			}
			logger.Info("create admin user success id:%s", admin.Id)
		}
	} else {
		logger.Info("found admin user with id:%s", admin.Id)
	}

	return &service{
		UserProfileDb:        userProfileDb,
		SessionConfig:        sessionConfig,
		ForgetPWConfig:       forgetPWConfig,
		JwtUtil:              jwtutil.NewUtil(),
		PasswordUtil:         passwordutil.NewUtil(),
		EmailUtil:            emailUtil,
		UserProfileMemorizer: userGetter,
	}
}
