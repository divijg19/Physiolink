package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/auth"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/service"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

type authResponse struct {
	Token   string                 `json:"token"`
	Profile map[string]interface{} `json:"profile,omitempty"`
}

type errorResponse struct {
	Msg string `json:"msg"`
}

// Define a minimal interface for the auth service to enable testing.
type AuthService interface {
	Register(ctx context.Context, email, password, role string) (uuid.UUID, string, error)
	Authenticate(ctx context.Context, email, password string) (uuid.UUID, string, error)
}

var authService AuthService
var cfg *config.Config

func InitAuth(s AuthService, c *config.Config) {
	authService = s
	cfg = c
}

// tokenIssuer builds an issuer from the package-level config. Every token in
// the system is minted here so the claim shape has one definition.
func tokenIssuer() *auth.Issuer {
	return auth.NewIssuer(cfg.JWTSecret, auth.TTL)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid request"})
		return
	}
	ctx := r.Context()
	id, role, err := authService.Register(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "User already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server error"})
		return
	}
	signed, err := tokenIssuer().Issue(id, role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server error"})
		return
	}
	// Create empty profile linked to user to mirror Node behavior
	var prof map[string]interface{}
	if profileService != nil {
		if p, err := profileService.CreateEmptyProfile(ctx, id); err == nil {
			prof = p
		}
	}
	writeJSON(w, http.StatusOK, authResponse{Token: signed, Profile: prof})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "invalid request"})
		return
	}
	ctx := r.Context()
	id, role, err := authService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Msg: "Invalid Credentials"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server error"})
		return
	}
	signed, err := tokenIssuer().Issue(id, role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Msg: "Server error"})
		return
	}
	writeJSON(w, http.StatusOK, authResponse{Token: signed})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
