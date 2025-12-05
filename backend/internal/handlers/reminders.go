package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// ReminderService interface for handler tests
type ReminderService interface {
	ListForPatient(ctx context.Context, patientID uuid.UUID) ([]service.ReminderItem, error)
}

var reminderService ReminderService

func InitReminders(s ReminderService) { reminderService = s }

func GetMyReminders(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"msg": "unauthorized"})
		return
	}
	uid, err := uuid.Parse(sub)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"msg": "invalid user"})
		return
	}
	list, err := reminderService.ListForPatient(r.Context(), uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"msg": "Server Error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
