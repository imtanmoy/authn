package authlib

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/httpx"
	"net/http"
	"strings"
	"time"
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

func GenerateToken(identity string) (string, error) {
	if auth == nil {
		return "", errors.New("authlib is not initiated")
	}
	now := time.Now()
	expirationTime := now.Add(time.Duration(auth.accessTokenExpireTime) * time.Minute)
	claims := &Claims{
		Identity: identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Subject:   identity,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(auth.secretKey))
	return tokenString, err
}

func AuthMiddleware(next http.Handler) http.Handler {
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
		if token == "" {
			errorMsg := "Required authorization token not found"
			httpx.ResponseJSONError(w, r, 400, 400, errorMsg)
			return
		}
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return auth.secretKey, nil
		})
		if err != nil {
			message := ""
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				message := fmt.Sprintf("Unexpected signing method: %v", parsedToken.Header["alg"])
				httpx.ResponseJSONError(w, r, http.StatusBadRequest, message)
				return
			}
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "malformed token"
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					message = "token is expired"
				} else {
					message = ve.Error()
				}
				if message == "" {
					message = httpx.ErrInternalServerError.Error()
				}
				httpx.ResponseJSONError(w, r, http.StatusBadRequest, message)
				return
			}
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
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
		setCurrentUserAndServe(w, r, next, claims.Identity)
	})
}

// fromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func fromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", &AuthError{Message: "No authorization header present", Code: http.StatusBadRequest, Status: http.StatusBadRequest}
	}
	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", &AuthError{Message: "authorization header format must be bearer type", Code: http.StatusBadRequest, Status: http.StatusBadRequest}
	}
	return authHeaderParts[1], nil
}

func setCurrentUserAndServe(w http.ResponseWriter, r *http.Request, next http.Handler, identity string) {
	ctx := r.Context()
	if ctx == nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	u, err := auth.GetUser(ctx, identity)
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

func GetCurrentUser(r *http.Request) (*models.User, error) {
	ctx := r.Context()
	u, ok := ctx.Value(identityKey).(*models.User)
	if !ok {
		return nil, errorx.ErrInternalServer
	}
	return u, nil
}
