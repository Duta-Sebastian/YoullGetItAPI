package api

import (
	"YoullGetItAPI/database"
	"YoullGetItAPI/middleware"
	"YoullGetItAPI/models"
	"YoullGetItAPI/util"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

// RegisterSyncPullRoute handles data sync pull requests
func RegisterSyncPullRoute(router *http.ServeMux, db *sql.DB) {
	router.Handle("/api/sync/pull", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := util.GetClaimsFromRequest(r)
			userId := claims.UserId

			if !claims.HasScope("sync:pull") {
				util.RespondWithError(w, http.StatusForbidden, "Insufficient scope.")
				return
			}

			tableParam := r.URL.Query().Get("table")
			if !util.IsTableAllowed(tableParam) {
				util.RespondWithError(w, http.StatusBadRequest,
					"Invalid 'table' parameter. Must be one of: job_cart, auth_user, cv.")
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
			defer cancel()

			var records interface{}
			var err error

			switch tableParam {
			case "job_cart":
				records, err = database.GetJobCartSyncPullData(ctx, db, userId)
			case "auth_user":
				records, err = database.GetUserSyncPullData(ctx, db, userId)
			case "cv":
				records, err = database.GetCvSyncPullData(ctx, db, userId)
			}

			if err != nil {
				log.Printf("Database query failed for table %s: %v", tableParam, err)
				util.RespondWithError(w, http.StatusServiceUnavailable, "Database query failed.")
				return
			}

			if util.IsRecordEmpty(records) {
				util.RespondWithJSON(w, http.StatusNoContent, nil)
				return
			}

			util.RespondWithJSON(w, http.StatusOK, records)
		})))
}

// RegisterSyncPushRoutes handles data sync push requests
func RegisterSyncPushRoutes(router *http.ServeMux, db *sql.DB) {
	router.Handle("/api/sync/push", middleware.EnsureValidToken()(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := util.GetClaimsFromRequest(r)
			userId := claims.UserId

			if !claims.HasScope("sync:push") {
				util.RespondWithError(w, http.StatusForbidden, "Insufficient scope.")
				return
			}

			tableParam := r.URL.Query().Get("table")
			if !util.IsTableAllowed(tableParam) {
				util.RespondWithError(w, http.StatusBadRequest,
					"Invalid 'table' parameter. Must be one of: job_cart, auth_user, cv.")
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
			defer cancel()

			var err error

			switch tableParam {
			case "job_cart":
				var jobRecords []models.JobRecord
				if err = util.DecodeRequestBody(r, &jobRecords); err != nil {
					util.RespondWithError(w, http.StatusBadRequest, "Error decoding job cart JSON")
					return
				}

				if err = database.PostJobCartSyncPushData(ctx, db, userId, jobRecords); err != nil {
					util.RespondWithError(w, http.StatusServiceUnavailable, "Error posting job cart data")
					return
				}

			case "auth_user":
				var userRecords []models.UserRecord
				if err = util.DecodeRequestBody(r, &userRecords); err != nil {
					util.RespondWithError(w, http.StatusBadRequest, "Error decoding user JSON")
					return
				}

				if err = database.PostUserSyncPushData(ctx, db, userId, userRecords); err != nil {
					log.Printf("Error posting user data: %v", err)
					util.RespondWithError(w, http.StatusInternalServerError, "Error posting user data")
					return
				}

			case "cv":
				var cvRecords []models.CvRecord
				if err = util.DecodeRequestBody(r, &cvRecords); err != nil {
					util.RespondWithError(w, http.StatusBadRequest, "Error decoding cv JSON")
					return
				}

				if err = database.PostCvSyncPushData(ctx, db, userId, cvRecords); err != nil {
					log.Printf("Error posting CV data: %v", err)
					util.RespondWithError(w, http.StatusInternalServerError, "Error posting cv data")
					return
				}
			}

			util.RespondWithJSON(w, http.StatusOK, map[string]string{
				"message": fmt.Sprintf("%s push successful.", tableParam),
			})
		})))
}
