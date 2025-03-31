package util

import (
	"YoullGetItAPI/middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/http"
)

// GetClaimsFromRequest extracts the custom claims from the JWT token in the request context
func GetClaimsFromRequest(r *http.Request) *middleware.CustomClaims {
	token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	return claims
}

// GetUserIDFromRequest extracts the user ID from JWT token in the request context
func GetUserIDFromRequest(claims *middleware.CustomClaims) (string, error) {
	return claims.UserId, nil
}

// ValidateScope checks if the user has the required scope
func ValidateScope(r *http.Request, requiredScope string) bool {
	token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := token.CustomClaims.(*middleware.CustomClaims)
	return claims.HasScope(requiredScope)
}
