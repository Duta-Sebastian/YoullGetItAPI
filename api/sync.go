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
				errMsg := "Insufficient scope for sync:pull operation"
				log.Printf("Access denied for user %s: %s", userId, errMsg)
				util.RespondWithError(w, http.StatusForbidden, errMsg)
				return
			}

			tableParam := r.URL.Query().Get("table")
			if !util.IsTableAllowed(tableParam) {
				errMsg := fmt.Sprintf("Invalid 'table' parameter: %s. Must be one of: job_cart,"+
					" auth_user, cv, question.", tableParam)
				log.Printf("Invalid table request from user %s: %s", userId, errMsg)
				util.RespondWithError(w, http.StatusBadRequest, errMsg)
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
			case "question":
				records, err = database.GetQuestionSyncPullData(ctx, db, userId)
			}

			if err != nil {
				errMsg := fmt.Sprintf("Database query failed for table %s: %v", tableParam, err)
				log.Printf(errMsg)
				util.RespondWithError(w, http.StatusServiceUnavailable, errMsg)
				return
			}

			if util.IsRecordEmpty(records) {
				log.Printf("No records found for user %s in table %s", userId, tableParam)
				util.RespondWithJSON(w, http.StatusNoContent, nil)
				return
			}

			log.Printf("Successfully pulled data from %s for user %s", tableParam, userId)
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
				errMsg := "Insufficient scope for sync:push operation"
				log.Printf("Access denied for user %s: %s", userId, errMsg)
				util.RespondWithError(w, http.StatusForbidden, errMsg)
				return
			}

			tableParam := r.URL.Query().Get("table")
			if !util.IsTableAllowed(tableParam) {
				errMsg := fmt.Sprintf("Invalid 'table' parameter: %s. Must be one of: job_cart,"+
					" auth_user, cv, question.", tableParam)
				log.Printf("Invalid table request from user %s: %s", userId, errMsg)
				util.RespondWithError(w, http.StatusBadRequest, errMsg)
				return
			}

			log.Printf("Processing push request for table %s from user %s", tableParam, userId)

			ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
			defer cancel()

			var err error

			switch tableParam {
			case "job_cart":
				var jobRecords []models.JobRecord
				if err = util.DecodeRequestBody(r, &jobRecords); err != nil {
					errMsg := fmt.Sprintf("Error decoding job cart JSON: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusBadRequest, errMsg)
					return
				}

				log.Printf("Received %d job records for user %s", len(jobRecords), userId)
				for i, job := range jobRecords {
					log.Printf("Job Record #%d: %+v", i+1, job)
				}

				if err = database.PostJobCartSyncPushData(ctx, db, userId, jobRecords); err != nil {
					errMsg := fmt.Sprintf("Error posting job cart data: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusServiceUnavailable, errMsg)
					return
				}

			case "auth_user":
				var userRecords []models.UserRecord
				if err = util.DecodeRequestBody(r, &userRecords); err != nil {
					errMsg := fmt.Sprintf("Error decoding user JSON: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusBadRequest, errMsg)
					return
				}

				log.Printf("Received %d user records for user %s", len(userRecords), userId)
				for i, user := range userRecords {
					log.Printf("User Record #%d: %+v", i+1, user)
				}

				if err = database.PostUserSyncPushData(ctx, db, userId, userRecords); err != nil {
					errMsg := fmt.Sprintf("Error posting user data: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusInternalServerError, errMsg)
					return
				}

			case "cv":
				var cvRecords []models.CvRecord
				if err = util.DecodeRequestBody(r, &cvRecords); err != nil {
					errMsg := fmt.Sprintf("Error decoding cv JSON: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusBadRequest, errMsg)
					return
				}

				log.Printf("Received %d CV records for user %s", len(cvRecords), userId)
				for i, cv := range cvRecords {
					log.Printf("CV Record #%d: %+v", i+1, cv)
				}

				if err = database.PostCvSyncPushData(ctx, db, userId, cvRecords); err != nil {
					errMsg := fmt.Sprintf("Error posting CV data: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusInternalServerError, errMsg)
					return
				}

			case "question":
				var questionRecords []models.QuestionRecord
				if err = util.DecodeRequestBody(r, &questionRecords); err != nil {
					errMsg := fmt.Sprintf("Error decoding question JSON: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusBadRequest, errMsg)
					return
				}

				log.Printf("Received %d question records for user %s", len(questionRecords), userId)
				for i, question := range questionRecords {
					log.Printf("Question Record #%d: %+v", i+1, question)
				}

				if err = database.PostQuestionSyncPushData(ctx, db, userId, questionRecords); err != nil {
					errMsg := fmt.Sprintf("Error posting question data: %v", err)
					log.Printf(errMsg)
					util.RespondWithError(w, http.StatusInternalServerError, errMsg)
					return
				}
			}

			successMsg := fmt.Sprintf("%s push successful for user %s", tableParam, userId)
			log.Printf(successMsg)
			util.RespondWithJSON(w, http.StatusOK, map[string]string{
				"message": successMsg,
			})
		})))
}
