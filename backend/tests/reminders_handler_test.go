package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/divijg19/physiolink/backend/tests/__mocks__"
)

func mockAuthPatient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, uuid.New().String())
		ctx = context.WithValue(ctx, middleware.UserRoleKey, "patient")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestGetMyReminders_Success(t *testing.T) {
	mock := &__mocks__.ReminderServiceMock{
		ListResp: []service.ReminderItem{
			{ID: uuid.New().String(), Message: "Reminder: appointment on 2025-01-02T10:00:00Z", RemindAt: time.Now().UTC().Format(time.RFC3339)},
		},
	}
	r := chi.NewRouter()
	handlers.InitReminders(mock)
	r.Route("/reminders", func(r chi.Router) {
		r.Use(mockAuthPatient)
		r.Get("/me", handlers.GetMyReminders)
	})

	req := httptest.NewRequest(http.MethodGet, "/reminders/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp []service.ReminderItem
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 reminder, got %d", len(resp))
	}
}

func TestGetMyReminders_Unauthorized(t *testing.T) {
	mock := &__mocks__.ReminderServiceMock{}
	r := chi.NewRouter()
	handlers.InitReminders(mock)
	r.Get("/reminders/me", handlers.GetMyReminders)

	req := httptest.NewRequest(http.MethodGet, "/reminders/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetMyReminders_InvalidUser(t *testing.T) {
	mock := &__mocks__.ReminderServiceMock{}
	r := chi.NewRouter()
	handlers.InitReminders(mock)
	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, "not-a-uuid")
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}).Get("/reminders/me", handlers.GetMyReminders)

	req := httptest.NewRequest(http.MethodGet, "/reminders/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetMyReminders_ServiceError(t *testing.T) {
	mock := &__mocks__.ReminderServiceMock{ListErr: errors.New("db error")}
	r := chi.NewRouter()
	handlers.InitReminders(mock)
	r.With(mockAuthPatient).Get("/reminders/me", handlers.GetMyReminders)

	req := httptest.NewRequest(http.MethodGet, "/reminders/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
