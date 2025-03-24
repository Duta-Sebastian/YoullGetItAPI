package api

import (
	queries "YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"YoullGetItAPI/models"
	"encoding/json"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"net/http"
	"time"
)

func RegisterSyncPullRoute(router *http.ServeMux) {
	router.Handle("/api/sync/pull", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*middleware.CustomClaims)
			userId := claims.UserId

			if !claims.HasScope("sync:pull") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}

			sinceParam := r.URL.Query().Get("since")
			var sinceTime *time.Time = nil
			if sinceParam != "" {
				parsedTime, err := time.Parse(time.RFC3339, sinceParam)

				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(`{"message":"Invalid 'since' timestamp format"}`))
					return
				}
				sinceTime = &parsedTime
			}

			db, dbConnectionErr := queries.GetDBConnection("app_user")
			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			records, err := queries.GetSyncPullData(db, sinceTime, userId)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Database query failed."}`))
				return
			}

			if len(records) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(records)
		})))
}

func RegisterSyncPushRoutes(router *http.ServeMux) {
	router.Handle("/api/sync/push", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			claims := token.CustomClaims.(*middleware.CustomClaims)
			userId := claims.UserId

			if !claims.HasScope("sync:push") {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"message":"Insufficient scope."}`))
				return
			}

			var records []models.JobRecord
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&records); err != nil {
				http.Error(w, "Error decoding JSON", http.StatusBadRequest)
				return
			}

			db, dbConnectionErr := queries.GetDBConnection("app_user")
			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			err := queries.PostSyncPushData(db, records, userId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message":"Database query failed."}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Push successful."}`))
		})))
}
