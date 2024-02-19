package session

import (
	"context"
	"time"

	"api-hotel-booking/internal/app/utils/encryptutil"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/utils/jwtutil"
)

type ServiceSession interface {
	Login(ctx context.Context, req LoginRequest) (LoginResponse, error)
	ExtendSession(ctx context.Context, req ExtendRequest) (ExtendResponse, error)
	Logout(req LogoutRequest) error
}

type service struct {
	UserProfileDb
	SessionDb
	app.SessionConfig
	JwtUtil
	EncryptUtil
}

type UserProfileDb interface {
	GetByEmail(ctx context.Context, email string) (persistence.UserProfile, error)
	WrongPasswordCounterReset(userId string) error
	WrongPasswordCounterIncrease(userId string, lockIfOver int, changeStatusTo string) (int, error)
}

type SessionDb interface {
	Insert(ctx context.Context, userId, key string, timeOut time.Time) (string, error)
	Get(ctx context.Context, sessionId string) (persistence.Session, error)
	Delete(sessionId string) error
	UpdateTimeout(ctx context.Context, sessionId string, timeout time.Time) error
}

type JwtUtil interface {
	NewClaims(sessionId, userId, email, name, surname, companyId, role, expire string) jwtutil.JwtToken
	NewBlankClaims() jwtutil.JwtToken
	ParseAndVerifyToken(token string) (jwtutil.JwtToken, error)
}

type EncryptUtil interface {
	NewEncryptor(key string) encryptutil.Encryptor
}

func NewServiceSession(config app.SessionConfig, userProfileDb UserProfileDb, sessionDb SessionDb) ServiceSession {
	return &service{
		UserProfileDb: userProfileDb,
		SessionDb:     sessionDb,
		SessionConfig: config,
		JwtUtil:       jwtutil.NewUtil(),
		EncryptUtil:   encryptutil.NewUtil(),
	}
}
