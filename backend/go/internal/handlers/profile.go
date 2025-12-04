package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/go/internal/middleware"
	"github.com/divijg19/physiolink/backend/go/internal/service"
)

var profileService *service.ProfileService

func InitProfile(s *service.ProfileService) {
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
		http.Error(w, "not found", http.StatusNotFound)
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
	var p service.ProfileUpdate
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, err := profileService.UpsertProfile(r.Context(), userID, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}
