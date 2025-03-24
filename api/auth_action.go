package api

import (
	queries "YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/http"
)

// RegisterCreateUserRoute registers the create user route
func RegisterCreateUserRoute(router *http.ServeMux) {
	router.Handle("POST /api/auth_action/create_user", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*middleware.CustomClaims)
			userId := claims.UserId

			if !claims.HasScope("user:create") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}
			db, dbConnectionErr := queries.GetDBConnection("auth_action")

			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			userCreationErr := queries.CreateUser(db, userId)
			if userCreationErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message":"Creating new user failed."}`))
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message":"User created successfully."}`))
			return
		}),
	))
}
