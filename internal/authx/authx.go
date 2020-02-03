package authx

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/httpx"
	"net/http"
)

type contextKey string

const (
	identityKey contextKey = "identity"
)

// Create a struct that will be encoded to a JWT
type Claims struct {
	Identity string `json:"identity"`
	jwt.StandardClaims
}

type AuthxConfig struct {
	SecretKey             string
	AccessTokenExpireTime int
}

type Authx struct {
	userRepo AuthRepo
	config   *AuthxConfig
}

// AuthableUser is identified by a password
type AuthableUser interface {
	GetEmail() (email string)
	GetPassword() (password string)
	PutPassword(password string)
}

// AuthUser is identified by a password
type AuthUser interface {
	GetId() (id int)
	GetEmail() (email string)
}

type AuthRepo interface {
	ExistsByEmail(ctx context.Context, identity string) bool
	GetByEmail(ctx context.Context, identity string) (AuthUser, error)
}

func New(userRepo AuthRepo, config *AuthxConfig) *Authx {
	return &Authx{userRepo: userRepo, config: config}
}

func (ax *Authx) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := fromAuthHeader(r)
		if err != nil {
			var ae *AuthError
			if errors.As(err, &ae) {
				httpx.ResponseJSONError(w, r, ae.Status, ae.Code, ae.Message)
			} else {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		parsedToken, err := parseToken(token, ax.config.SecretKey)
		if err != nil {
			var ae *AuthError
			if errors.As(err, &ae) {
				httpx.ResponseJSONError(w, r, ae.Status, ae.Code, ae.Message)
			} else {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		if !parsedToken.Valid {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "Token is invalid")
			return
		}

		claims, ok := parsedToken.Claims.(*Claims)
		if !ok {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
			return
		}
		ax.setCurrentUserAndServe(w, r, next, claims.Identity)
	})
}

func (ax *Authx) getUser(ctx context.Context, identity string) (AuthUser, error) {
	if !ax.userRepo.ExistsByEmail(ctx, identity) {
		return nil, errorx.ErrorNotFound
	}
	u, err := ax.userRepo.GetByEmail(ctx, identity)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (ax *Authx) GetCurrentUser(r *http.Request) (AuthUser, error) {
	ctx := r.Context()
	u, ok := ctx.Value(identityKey).(AuthUser)
	if !ok {
		return nil, errorx.ErrInternalServer
	}
	return u, nil
}

func (ax *Authx) setCurrentUserAndServe(w http.ResponseWriter, r *http.Request, next http.Handler, identity string) {
	ctx := r.Context()
	if ctx == nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	u, err := ax.getUser(ctx, identity)
	if err != nil {
		if errors.Is(err, errorx.ErrorNotFound) {
			httpx.ResponseJSONError(w, r, http.StatusNotFound, "user not found", err)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	ctx = context.WithValue(r.Context(), identityKey, u)
	next.ServeHTTP(w, r.WithContext(ctx))
}

func (ax *Authx) GenerateToken(identity string) (string, error) {
	tokenString, err := createToken(identity, ax.config.SecretKey, ax.config.AccessTokenExpireTime)
	return tokenString, err
}

// VerifyPassword uses mechanisms to check that a password is correct.
// Returns nil on success otherwise there will be an error. Simply a helper
// to do the bcrypt comparison.
func (ax *Authx) VerifyPassword(user AuthableUser, password string) bool {
	return comparePasswords(user.GetPassword(), password)
}

func (ax *Authx) HashPassword(password string) (string, error) {
	return hashPassword(password)
}
