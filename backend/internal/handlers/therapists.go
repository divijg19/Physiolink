package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

// TherapistService interface for handler tests.
type TherapistService interface {
	GetAllTherapists(ctx context.Context, params service.TherapistQueryParams) (service.TherapistListResult, error)
	GetTherapistByID(ctx context.Context, id, date string) (map[string]interface{}, error)
}

var therapistService TherapistService

func InitTherapists(s TherapistService) { therapistService = s }

func GetAllTherapists(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	available := q.Get("available") == "true"
	params := service.TherapistQueryParams{
		Specialty: q.Get("specialty"),
		Location:  q.Get("location"),
		Page:      page,
		Limit:     limit,
		Sort:      q.Get("sort"),
		Date:      q.Get("date"),
		Available: available,
	}
	res, err := therapistService.GetAllTherapists(r.Context(), params)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server Error"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func GetTherapistByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	date := r.URL.Query().Get("date")
	res, err := therapistService.GetTherapistByID(r.Context(), id, date)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Msg: "Therapist not found"})
		return
	}
	writeJSON(w, http.StatusOK, res)
}
