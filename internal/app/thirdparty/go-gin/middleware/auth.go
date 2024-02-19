package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// auth middleware fo go-gin
// @maintainer: Minh Bach (minh.b@kbtg.tech)
//
//go:generate mockgen -source auth.go -package middleware -destination auth_mock.go
type (
	Authenticator struct {
		cfg          AuthConfig
		userStore    UserStore
		tokenStore   TokenStore
		em           AuthErrorMapTemplate
		keyFn        jwt.Keyfunc
		rsaKeys      rsaKeyData
		compareToken func(inputToken, storedToken string) bool
	}

	rsaKeyData struct {
		PublicKey  *rsa.PublicKey
		PrivateKey *rsa.PrivateKey
	}

	UserStore interface {
		CheckUserIdentity(ctx context.Context, username, password string) (UserData, error)
	}

	TokenStore interface {
		GetTokenByUser(ctx context.Context, username string) (string, error)
	}

	UserData struct {
		Username string   `json:"username" bson:"username" binding:"required,min=3,max=20"`
		Group    string   `json:"group,omitempty" bson:"group,omitempty" binding:"required,uppercase"`
		Scopes   []string `json:"scopes,omitempty" bson:"scopes,omitempty"`
	}

	UserClaims struct {
		*jwt.RegisteredClaims
		Group  string   `json:"group,omitempty"`
		Scopes []string `json:"scopes,omitempty"`
	}

	AuthConfig struct {
		Header     string        `mapstructure:"Header"`
		Expiration time.Duration `mapstructure:"Expiration"`
		Issuer     string        `mapstructure:"Issuer"`
		SignMethod string        `mapstructure:"SignMethod"`
		SignKey    SignKeyData   `mapstructure:"SignKeyPath"`
	}

	SignKeyData struct {
		RSAPublicKey  string
		RSAPrivateKey string
		SecretKey     string
	}

	AuthError struct {
		Code         string `mapstructure:"Code" json:"code"`
		ErrorMessage string `mapstructure:"Message" json:"message"`
		HttpCode     int    `mapstructure:"HttpCode" json:"-"`
	}

	AuthErrorMapTemplate struct {
		InvalidHeaderToken        AuthError `mapstructure:"InvalidHeaderToken"`
		TokenStoreNotFound        AuthError `mapstructure:"TokenStoreNotFound"`
		InvalidClaimsStructure    AuthError `mapstructure:"InvalidClaimsStructure"`
		ParseBearerTokenError     AuthError `mapstructure:"ParseJWTTokenError"`
		ParseBasicTokenError      AuthError `mapstructure:"ParseBasicTokenError"`
		InvalidUserIdentity       AuthError `mapstructure:"InvalidUserIdentity"`
		UserClaimsContextNotFound AuthError `mapstructure:"UserClaimsContextNotFound"`
		UserNotAuthorized         AuthError `mapstructure:"UserNotAuthorized"`
	}
)

const (
	UserClaimsContextKey = "user"
	SignMethodRSA        = "RSA"
	SignMethodHMACSHA    = "HMAC-SHA"
	ScopePermitAll       = "permit_all"
	ScopeAuthenticated   = "authenticated"
)

func NewAuthenticator(cfg AuthConfig, userStore UserStore, tokenStore TokenStore, em AuthErrorMapTemplate) (*Authenticator, error) {
	m := &Authenticator{
		cfg:        cfg,
		userStore:  userStore,
		tokenStore: tokenStore,
		em:         em,
		compareToken: func(inputToken, storedToken string) bool {
			return inputToken == storedToken
		},
	}

	switch cfg.SignMethod {
	case SignMethodRSA:
		// Load RSA private key
		privateKey, err := ioutil.ReadFile(cfg.SignKey.RSAPrivateKey)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to read rsa key at %s", cfg.SignKey.RSAPrivateKey)
		}
		privateRSA, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to parse rsa key at %s", cfg.SignKey.RSAPrivateKey)
		}

		// Load RSA public key
		publicKey, err := ioutil.ReadFile(cfg.SignKey.RSAPublicKey)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to read rsa key at %s", cfg.SignKey.RSAPublicKey)
		}
		publicRSA, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to parse rsa key at %s", cfg.SignKey.RSAPublicKey)
		}
		m.rsaKeys = rsaKeyData{
			PublicKey:  publicRSA,
			PrivateKey: privateRSA,
		}
		m.keyFn = m.RSAKeyFunction
	case SignMethodHMACSHA:
		m.keyFn = m.HMACKeyFunction
	default:
		return nil, errors.New("invalid sign method in auth properties, 'RSA' or 'HMAC-SHA' required")
	}
	return m, nil
}

func (a *Authenticator) GenerateAuthToken(user UserData) (string, error) {
	claims := UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(a.cfg.Expiration)},
			Issuer:    a.cfg.Issuer,
			ID:        user.Username,
			Subject:   user.Username,
		},
		Group:  user.Group,
		Scopes: user.Scopes,
	}
	var signMethod jwt.SigningMethod
	var signKey interface{}
	switch a.cfg.SignMethod {
	case SignMethodRSA:
		signMethod = jwt.SigningMethodRS256
		signKey = a.rsaKeys.PrivateKey
	case SignMethodHMACSHA:
		signMethod = jwt.SigningMethodHS256
		signKey = []byte(a.cfg.SignKey.SecretKey)
	}

	token := jwt.NewWithClaims(signMethod, claims)
	return token.SignedString(signKey)
}

func (a *Authenticator) Authenticate(ctx *gin.Context) {
	// extract token from header
	authHeader := ctx.GetHeader(a.cfg.Header)
	if authHeader == "" {
		return
	}

	arr := strings.SplitN(authHeader, " ", 2)
	if len(arr) != 2 {
		a.abortError(ctx, a.em.InvalidHeaderToken)
		return
	}

	switch strings.ToUpper(arr[0]) {
	case "BEARER":
		a.BearerAuthHandler(ctx, arr[1])
	case "BASIC":
		a.BasicAuthHandler(ctx, arr[1])
	default:
		a.abortError(ctx, a.em.InvalidHeaderToken)
		return
	}
}

func (a *Authenticator) Authorize(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if requiredScope == ScopePermitAll {
			return
		}

		if u, err := GetContextUserClaims(ctx); err == nil && u != nil {
			if requiredScope == ScopeAuthenticated {
				return
			}

			for _, s := range u.Scopes {
				if s == requiredScope {
					return
				}
			}
			a.abortError(ctx, a.em.UserNotAuthorized)
			return
		}
		a.abortError(ctx, a.em.UserClaimsContextNotFound)
	}
}

func (a *Authenticator) BearerAuthHandler(ctx *gin.Context, tokenString string) {
	// parse token using defined signing method
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, a.keyFn)
	if err != nil || !token.Valid {
		a.abortError(ctx, a.em.ParseBearerTokenError)
		return
	}

	// verify whether token still in token store
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		a.abortError(ctx, a.em.InvalidClaimsStructure)
		return
	}

	if storedToken, err := a.tokenStore.GetTokenByUser(ctx, claims.Username()); err != nil || !a.compareToken(tokenString, storedToken) {
		a.abortError(ctx, a.em.TokenStoreNotFound)
		return
	}

	// set token claims into context for further usage
	ctx.Set(UserClaimsContextKey, claims)
}

func (a *Authenticator) BasicAuthHandler(ctx *gin.Context, tokenString string) {
	b, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		a.abortError(ctx, a.em.ParseBasicTokenError)
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		a.abortError(ctx, a.em.ParseBasicTokenError)
		return
	}

	user, err := a.userStore.CheckUserIdentity(ctx, pair[0], pair[1])
	if err != nil {
		a.abortError(ctx, a.em.InvalidUserIdentity)
		return
	}

	ctx.Set(UserClaimsContextKey, &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{Subject: user.Username},
		Scopes:           user.Scopes,
		Group:            user.Group,
	})
}

func (a *Authenticator) HMACKeyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
	}
	return []byte(a.cfg.SignKey.SecretKey), nil
}

func (a *Authenticator) RSAKeyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
	}
	return a.rsaKeys.PublicKey, nil
}

func (a *Authenticator) abortError(ctx *gin.Context, err AuthError) {
	if err.Code == "" {
		err = AuthError{
			Code:         "EC9998",
			ErrorMessage: "Undefined authentication error",
			HttpCode:     http.StatusUnauthorized,
		}
	}

	ctx.AbortWithStatusJSON(err.HttpCode, gin.H{"error": err})
}

func (ae AuthError) Error() string {
	return fmt.Sprintf("%s:%s", ae.Code, ae.ErrorMessage)
}

func (u *UserClaims) Username() string {
	return u.Subject
}

func GetContextUserClaims(ctx context.Context) (*UserClaims, error) {
	if claims, ok := ctx.Value(UserClaimsContextKey).(*UserClaims); ok {
		return claims, nil
	}
	return nil, errors.New("no valid user_claims in context")
}
