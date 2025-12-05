package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// ReviewService interface for handler tests.
type ReviewService interface {
	CreateReview(ctx context.Context, patientID, therapistID uuid.UUID, rating int, comment string) (map[string]interface{}, error)
	GetReviewsForTherapist(ctx context.Context, therapistID uuid.UUID) ([]map[string]interface{}, error)
}

var reviewService ReviewService

func InitReviews(s ReviewService) { reviewService = s }

type createReviewReq struct {
	TherapistID string `json:"therapistId"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment"`
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Msg: "unauthorized"})
		return
	}
	pid, err := uuid.Parse(sub)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid user"})
		return
	}
	var req createReviewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid request"})
		return
	}
	tid, err := uuid.Parse(req.TherapistID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid therapistId"})
		return
	}
	res, err := reviewService.CreateReview(r.Context(), pid, tid, req.Rating, req.Comment)
	if err != nil {
		if fe, ok := err.(*service.ForbiddenError); ok {
			writeJSON(w, http.StatusForbidden, errorResponse{Msg: fe.Msg})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server Error"})
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func GetReviewsForTherapist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "therapistId")
	tid, err := uuid.Parse(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid therapistId"})
		return
	}
	res, err := reviewService.GetReviewsForTherapist(r.Context(), tid)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server Error"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
