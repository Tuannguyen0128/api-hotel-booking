package service

import (
	"context"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/persistence"
	"api-hotel-booking/internal/app/service/user"
	"api-hotel-booking/internal/app/utils/cache"
	"api-hotel-booking/internal/app/utils/memorizer"

	"api-hotel-booking/internal/app/service/session"
	"api-hotel-booking/internal/app/thirdparty/logger"
)

var log = logger.WithModule("service")

type Service struct { // This struct must implement Service interface in handler package (route.go)
	sSession session.ServiceSession
	sUser    user.ServiceUser
}

func NewService(a *Adapters) *Service {
	userRepo := a.UserRepo()
	sessionRepo := a.SessionRepo()
	userProfileMemorizer := memorizer.NewMemorizer[persistence.UserProfile](
		app.CFG.Service.Cache.Company.MemorizerTimeInMs,
		cache.NewCache[persistence.UserProfile](
			app.CFG.Service.Cache.Company.CacheSize,
			app.CFG.Service.Cache.Company.CacheTimeInSec),
		userRepo.GetById,
	)
	return &Service{
		sSession: session.NewServiceSession(app.CFG.Service.Session, userRepo, sessionRepo),
		sUser:    user.NewServiceUser(app.CFG.Service.Session, app.CFG.Service.ForgetPW, userProfileMemorizer, userRepo, a.emailUtil),
	}
}

func (s *Service) Login(ctx context.Context, req session.LoginRequest) (session.LoginResponse, error) {
	return s.sSession.Login(ctx, req)
}

func (s *Service) ExtendSession(ctx context.Context, req session.ExtendRequest) (session.ExtendResponse, error) {
	return s.sSession.ExtendSession(ctx, req)
}

func (s *Service) Logout(req session.LogoutRequest) error {
	return s.sSession.Logout(req)
}

func (s *Service) CreateUser(ctx context.Context, req user.CreateRequest) (user.CreateResponse, error) {
	return s.sUser.CreateUser(ctx, req)
}

func (s *Service) EditUser(ctx context.Context, req user.EditRequest) (user.EditResponse, error) {
	return s.sUser.EditUser(ctx, req)
}

func (s *Service) GetUser(ctx context.Context, req user.GetRequest) (user.GetResponse, error) {
	return s.sUser.GetUser(ctx, req)
}

func (s *Service) ChangePassword(ctx context.Context, req user.ChangePassRequest) (user.ChangePassResponse, error) {
	return s.sUser.ChangePassword(ctx, req)
}

func (s *Service) ListUser(ctx context.Context, req user.ListRequest) (user.ListResponse, error) {
	return s.sUser.ListUser(ctx, req)
}

func (s *Service) DeleteUser(ctx context.Context, req user.DeleteRequest) (user.DeleteResponse, error) {
	return s.sUser.DeleteUser(ctx, req)
}

func (s *Service) ForgetPassword(ctx context.Context, req user.ForgetPassRequest) error {
	return s.sUser.ForgetPassword(ctx, req)
}

func (s *Service) ResetPassword(ctx context.Context, req user.ResetPassRequest) (user.ResetPassResponse, error) {
	return s.sUser.ResetPassword(ctx, req)
}
