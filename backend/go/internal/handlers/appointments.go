package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/go/internal/middleware"
	"github.com/divijg19/physiolink/backend/go/internal/service"
)

var apptService *service.AppointmentService

func InitAppointments(s *service.AppointmentService) { apptService = s }

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
	// path: /api/appointments/{id}/book
	// simple parse for last segment
	parts := splitPath(r.URL.Path)
	if len(parts) < 4 {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	slotID, err := uuid.Parse(parts[3])
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

func splitPath(p string) []string {
	var out []string
	start := 0
	for i := 0; i < len(p); i++ {
		if p[i] == '/' {
			if i > start {
				out = append(out, p[start:i])
			}
			start = i + 1
		}
	}
	if start < len(p) {
		out = append(out, p[start:])
	}
	return out
}
