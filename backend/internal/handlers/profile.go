package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// ProfileService interface to enable testing mocks.
type ProfileService interface {
	UpsertProfile(ctx context.Context, userID uuid.UUID, p service.NodeProfileUpdate) (uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
	CreateEmptyProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error)
}

var profileService ProfileService

func InitProfile(s ProfileService) {
	profileService = s
}

func GetMyProfile(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	data, err := profileService.GetProfile(r.Context(), userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"msg": "There is no profile for this user"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func UpsertMyProfile(w http.ResponseWriter, r *http.Request) {
	sub, _ := r.Context().Value(middleware.UserIDKey).(string)
	if sub == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	var p service.NodeProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	_, err = profileService.UpsertProfile(r.Context(), userID, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the full profile like the Node service does
	data, err := profileService.GetProfile(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
