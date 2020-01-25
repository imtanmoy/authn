package authlib

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/httpx"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	identityKey contextKey = "identity"
)

var jwtKey = []byte(config.Conf.SECRET_KEY)

// Create a struct that will be encoded to a JWT
type Claims struct {
	Identity string `json:"identity"`
	jwt.StandardClaims
}

func GenerateToken(identity string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(60 * time.Minute)
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
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
			return
		}
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
			return jwtKey, nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				message := fmt.Sprintf("Unexpected signing method: %v", parsedToken.Header["alg"])
				httpx.ResponseJSONError(w, r, http.StatusBadRequest, message)
				return
			}
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
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
		fmt.Println(claims.Identity)
		ctx = context.WithValue(r.Context(), identityKey, claims.Identity)
		next.ServeHTTP(w, r.WithContext(ctx))
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
