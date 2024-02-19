package session

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"api-hotel-booking/internal/app"
	"api-hotel-booking/internal/app/thirdparty/hash"
	"api-hotel-booking/internal/app/utils"
	"api-hotel-booking/internal/app/utils/header/pheader"
	"github.com/aead/ecdh"
)

type LoginRequest struct {
	Header       pheader.Request `json:"header"`
	Email        string          `json:"email"`
	Password     string          `json:"password"`
	ClientPublic string          `json:"clientPublic"`
}

type LoginResponse struct {
	Header        pheader.Response `json:"header"`
	SessionId     string           `json:"sessionId"`
	ServerPublic  string           `json:"serverPublic"`
	Token         string           `json:"token"`
	RequireAction string           `json:"requireAction,omitempty"`
}

func (r *LoginRequest) Validate() error {
	if err := r.Header.Validate(false); err != nil {
		return err
	}
	if r.Email == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("email cannot be empty")
	}
	if r.Password == "" {
		return app.ErrorMap.InvalidRequest.AddDebug("password cannot be empty")
	}
	if r.ClientPublic == "" && app.CFG.Service.Session.IsRequireEncryption {
		return app.ErrorMap.InvalidRequest.AddDebug("client public key cannot be empty")
	}

	return nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	if err := req.Validate(); err != nil {
		return LoginResponse{}, err
	}
	password := req.Password

	profile, err := s.UserProfileDb.GetByEmail(ctx, req.Email)
	if err != nil {
		return LoginResponse{}, app.ErrorMap.LoginFailed.AddDebug(err.Error(), "cannot find user profile")
	}

	// check status
	if profile.Status != s.AllowLoginStatus {
		return LoginResponse{}, app.ErrorMap.LoginFailed.AddDebug(fmt.Sprintf("profile status not %s", s.AllowLoginStatus))
	}

	// check Password
	if isMatch, err := hash.ComparePassword(profile.GetActivePassword(), password); err != nil {
		return LoginResponse{}, app.ErrorMap.LoginFailed.AddDebug("error when validate password")
	} else if isMatch == false {
		if i, err := s.UserProfileDb.WrongPasswordCounterIncrease(profile.Id, s.WrongPasswordLockLimit, s.WrongPasswordLockToStatus); err != nil {
			return LoginResponse{}, app.ErrorMap.PasswordNotMatch.From(err).AddDebug(fmt.Sprintf("the password is wrong %d times", i))
		} else {
			return LoginResponse{}, app.ErrorMap.PasswordNotMatch.AddDebug(fmt.Sprintf("the password is wrong %d times", i))
		}
	} else if profile.WrongPassword != 0 {
		if err := s.UserProfileDb.WrongPasswordCounterReset(profile.Id); err != nil {
			return LoginResponse{}, app.ErrorMap.InternalServerError.From(err).AddDebug("cannot reset wrong password count")
		}
	}

	// disable handshake
	sharedSecretB64 := "noKey"
	serverPublicB64 := "noKey"
	if app.CFG.Service.Session.IsRequireEncryption {
		x25519 := ecdh.X25519()
		clientPublic, err := base64.StdEncoding.DecodeString(req.ClientPublic)
		if err != nil {
			return LoginResponse{}, app.ErrorMap.InvalidRequest.AddDebug(err.Error(), "cannot base64 decode key")
		}

		if err := x25519.Check(clientPublic); err != nil {
			return LoginResponse{}, app.ErrorMap.InvalidRequest.AddDebug(err.Error(), "cannot use this key")
		}

		serverPrivate, serverPublic, err := x25519.GenerateKey(nil)
		if err != nil {
			return LoginResponse{}, app.ErrorMap.InternalServerError.AddDebug(err.Error(), "cannot generate x25519 key")
		}

		serverPublicBytes := serverPublic.([32]byte)
		serverPublicB64 = base64.StdEncoding.EncodeToString(serverPublicBytes[:])

		sharedSecretBytes := x25519.ComputeSecret(serverPrivate, clientPublic)
		sharedSecretB64 = base64.StdEncoding.EncodeToString(sharedSecretBytes)
	}

	timeout := time.Now().Add(time.Duration(s.Timeout) * time.Second)
	sessionId, err := s.SessionDb.Insert(ctx, profile.Id, sharedSecretB64, timeout)
	if err != nil {
		return LoginResponse{}, err
	}

	claims := s.JwtUtil.NewClaims(sessionId, profile.Id, profile.Email, profile.Name, profile.Surname, profile.CompanyId, profile.Role, timeout.Format(utils.RFC3339))
	token, err := claims.GenerateToken()
	if err != nil {
		return LoginResponse{}, app.ErrorMap.InternalServerError.AddDebug(err.Error(), "error when generating token")
	}

	if profile.TempPassword != "" {
		return LoginResponse{
			Header:        pheader.Response{},
			SessionId:     sessionId,
			ServerPublic:  serverPublicB64,
			Token:         token,
			RequireAction: "changePassword",
		}, nil
	}

	return LoginResponse{
		Header:       pheader.Response{},
		SessionId:    sessionId,
		ServerPublic: serverPublicB64,
		Token:        token,
	}, nil
}
