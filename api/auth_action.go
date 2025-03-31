package api

import (
	"YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"YoullGetItAPI/models"
	"YoullGetItAPI/util"
	"context"
	"database/sql"
	"net/http"
	"time"
)

// RegisterCreateUserRoute registers the create user route
func RegisterCreateUserRoute(router *http.ServeMux, db *sql.DB) {
	router.Handle("POST /api/auth_action/create_user", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := util.GetClaimsFromRequest(r)

			if !claims.HasScope("user:create") {
				util.RespondWithError(w, http.StatusForbidden, "Insufficient scope.")
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
			defer cancel()

			var req models.CreateUserModel
			if err := util.DecodeRequestBody(r, &req); err != nil {
				util.RespondWithError(w, http.StatusBadRequest, "Invalid JSON body.")
				return
			}

			if req.UserId == "" {
				util.RespondWithError(w, http.StatusBadRequest, "Missing userId.")
				return
			}

			userCreationErr := database.CreateUser(ctx, db, req.UserId)
			if userCreationErr != nil {
				util.RespondWithError(w, http.StatusInternalServerError, "Creating new user failed.")
				return
			}

			util.RespondWithJSON(w, http.StatusCreated, map[string]string{
				"message": "User created successfully.",
			})
			return
		}),
	))
}
