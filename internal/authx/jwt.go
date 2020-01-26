package authx

import (
	"net/http"
	"strings"
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
