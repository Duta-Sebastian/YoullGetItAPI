package api

import (
	"YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"YoullGetItAPI/models"
	"YoullGetItAPI/util"
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"log"
	"net/http"
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

			tableParam := r.URL.Query().Get("table")

			if !util.IsTableAllowed(tableParam) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"Invalid 'table' parameter. Must be one of: job_cart, auth_user, cv."}`))
			}

			db, dbConnectionErr := queries.GetDBConnection("app_user")
			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			var records interface{}
			var err error

			switch tableParam {
			case "job_cart":
				//records, err = queries.GetJobCartSyncPullData(db, userId, sinceTime)
			case "auth_user":
				records, err = queries.GetUserSyncPullData(db, userId)
			case "cv":
				records, err = queries.GetCvSyncPullData(db, userId)
			}

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Database query failed."}`))
				return
			}

			if util.IsRecordEmpty(records) {
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

			tableParam := r.URL.Query().Get("table")

			if !util.IsTableAllowed(tableParam) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"message":"Invalid 'table' parameter. Must be one of: job_cart, auth_user, cv."}`))
				return
			}

			db, dbConnectionErr := queries.GetDBConnection("app_user")
			if dbConnectionErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"message":"Connection to database failed."}`))
				return
			}

			var jobRecords []models.JobRecord
			var userRecords []models.UserRecord
			var cvRecords []models.CvRecord

			switch tableParam {
			case "job_cart":
				{
					decoder := json.NewDecoder(r.Body)
					if err := decoder.Decode(&jobRecords); err != nil {
						http.Error(w, "Error decoding JSON", http.StatusBadRequest)
						return
					}
				}
			case "auth_user":
				{
					decoder := json.NewDecoder(r.Body)
					if err := decoder.Decode(&userRecords); err != nil {
						log.Println("Error decoding JSON:", err)
						http.Error(w, "Error decoding JSON", http.StatusBadRequest)
						return
					}
					if err := queries.PostUserSyncPushData(db, userId, userRecords); err != nil {
						log.Println("Error decoding JSON:", err)
						http.Error(w, "Error posting data to database", http.StatusBadRequest)
						return
					}
				}
			case "cv":
				{
					decoder := json.NewDecoder(r.Body)
					if err := decoder.Decode(&cvRecords); err != nil {
						log.Println("Error decoding JSON:", err)
						http.Error(w, "Error decoding JSON", http.StatusBadRequest)
						return
					}
					if err := queries.PostCvSyncPushData(db, userId, cvRecords); err != nil {
						log.Println("Error decoding JSON:", err)
						http.Error(w, "Error posting data to database", http.StatusBadRequest)
						return
					}
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"message":"%s push successful."}`, tableParam)))
		})))
}
