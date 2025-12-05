package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// AppointmentService interface for handler tests.
type AppointmentService interface {
	CreateAvailability(ctx context.Context, therapistID uuid.UUID, slots []struct{ StartTs, EndTs string }) error
	GetTherapistAvailability(ctx context.Context, therapistID uuid.UUID) ([]service.Slot, error)
	BookAppointment(ctx context.Context, slotID, patientID uuid.UUID) (uuid.UUID, error)
	ListMyAppointments(ctx context.Context, userID uuid.UUID, role string) ([]service.AppointmentBrief, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID, therapistID uuid.UUID, status string) (service.AppointmentBrief, error)
}

var apptService AppointmentService

func InitAppointments(s AppointmentService) { apptService = s }

type createAvailReq struct {
	Slots []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"slots"`
}

func CreateAvailability(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	tid, err := uuid.Parse(sub)
	if err != nil {
		http.Error(w, "bad user", http.StatusBadRequest)
		return
	}
	var req createAvailReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	slots := make([]struct{ StartTs, EndTs string }, 0, len(req.Slots))
	for _, s := range req.Slots {
		slots = append(slots, struct{ StartTs, EndTs string }{StartTs: s.StartTime, EndTs: s.EndTime})
	}
	if err := apptService.CreateAvailability(r.Context(), tid, slots); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetTherapistAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("ptId")
	if idStr == "" {
		idStr = chi.URLParam(r, "ptId")
	}
	if idStr == "" {
		http.Error(w, "missing ptId", http.StatusBadRequest)
		return
	}
	tid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	slots, err := apptService.GetTherapistAvailability(r.Context(), tid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(slots)
}

func BookAppointment(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	pid, err := uuid.Parse(sub)
	if err != nil {
		http.Error(w, "bad user", http.StatusBadRequest)
		return
	}
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	slotID, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	id, err := apptService.BookAppointment(r.Context(), slotID, pid)
	if err != nil {
		if errors.Is(err, service.ErrConflict) {
			http.Error(w, "conflict", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}

func GetMyAppointments(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Msg: "unauthorized"})
		return
	}
	uid, err := uuid.Parse(sub)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid user"})
		return
	}
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)
	if role == "" {
		role = "patient"
	}
	list, err := apptService.ListMyAppointments(r.Context(), uid, role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server Error"})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

type updateStatusReq struct {
	Status string `json:"status"`
}

func UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Msg: "unauthorized"})
		return
	}
	ptID, err := uuid.Parse(sub)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid user"})
		return
	}
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "bad path"})
		return
	}
	apptID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "bad id"})
		return
	}
	var req updateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid request"})
		return
	}
	updated, err := apptService.UpdateAppointmentStatus(r.Context(), apptID, ptID, req.Status)
	if err != nil {
		switch err.Error() {
		case "forbidden":
			writeJSON(w, http.StatusForbidden, errorResponse{Msg: "Forbidden"})
			return
		case "invalid status":
			writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "Invalid status"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server Error"})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}
