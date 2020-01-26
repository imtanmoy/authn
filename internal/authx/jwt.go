package authx

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/imtanmoy/httpx"
	"net/http"
	"strings"
	"time"
)

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

func createToken(identity, secreteKey string, expireTime int) (string, error) {
	now := time.Now()
	expirationTime := now.Add(time.Duration(expireTime) * time.Minute)
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
	tokenString, err := token.SignedString([]byte(secreteKey))
	return tokenString, err
}

func parseToken(token, secretKey string) (*jwt.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		message := ""
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			message := fmt.Sprintf("Unexpected signing method: %v", parsedToken.Header["alg"])
			return nil, &AuthError{Message: message, Code: http.StatusBadRequest, Status: http.StatusBadRequest}
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
			return nil, &AuthError{Message: message, Code: http.StatusBadRequest, Status: http.StatusBadRequest}
		}
		return nil, err
	}
	return parsedToken, nil
}
