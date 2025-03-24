package api

import (
	queries "YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"encoding/json"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/http"
)

type CreateUserRequest struct {
	UserId string `json:"userId"`
}

// RegisterCreateUserRoute registers the create user route
func RegisterCreateUserRoute(router *http.ServeMux) {
	router.Handle("POST /api/auth_action/create_user", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*middleware.CustomClaims)

			if !claims.HasScope("user:create") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}

			var req CreateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"message":"Invalid JSON body."}`, http.StatusBadRequest)
				return
			}

			if req.UserId == "" {
				http.Error(w, `{"message":"Missing userId."}`, http.StatusBadRequest)
				return
			}

			db, dbConnectionErr := queries.GetDBConnection("auth_action")

			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			userCreationErr := queries.CreateUser(db, req.UserId)
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
