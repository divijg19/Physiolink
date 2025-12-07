package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/divijg19/physiolink/backend/tests/__mocks__"
)

func setupReviewsRouter(service *__mocks__.ReviewServiceMock) *chi.Mux {
	r := chi.NewRouter()
	handlers.InitReviews(service)
	r.Route("/reviews", func(r chi.Router) {
		r.Use(mockAuth)
		r.Post("/", handlers.CreateReview)
		r.Get("/{therapistId}", handlers.GetReviewsForTherapist)
	})
	return r
}

// MockAuth middleware for tests in this package
func mockAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// default authenticated user as patient
		ctx := r.Context()
		ctx = context.WithValue(ctx, middleware.UserIDKey, uuid.New().String())
		ctx = context.WithValue(ctx, middleware.UserRoleKey, "patient")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestCreateReview_Success(t *testing.T) {
	m := &__mocks__.ReviewServiceMock{
		CreateResp: map[string]interface{}{
			"id": uuid.New().String(),
		},
	}
	r := chi.NewRouter()
	handlers.InitReviews(m)
	r.Route("/reviews", func(r chi.Router) {
		r.Use(mockAuth)
		r.Post("/", handlers.CreateReview)
	})

	therapistID := uuid.New().String()
	body := map[string]interface{}{"therapistId": therapistID, "rating": 5, "comment": "great"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("expected 201/200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetReviewsForTherapist_Success(t *testing.T) {
	m := &__mocks__.ReviewServiceMock{
		ListResp: []map[string]interface{}{
			{"rating": 4, "comment": "good"},
			{"rating": 5, "comment": "great"},
		},
	}
	r := setupReviewsRouter(m)

	therapistID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/reviews/"+therapistID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateReview_ForbiddenForTherapistRole(t *testing.T) {
	m := &__mocks__.ReviewServiceMock{
		CreateResp: map[string]interface{}{"id": uuid.New().String()},
		CreateErr:  &service.ForbiddenError{Msg: "forbidden"},
	}
	r := chi.NewRouter()
	handlers.InitReviews(m)
	r.Route("/reviews", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler { // inject therapist role
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				ctx = context.WithValue(ctx, middleware.UserIDKey, uuid.New().String())
				ctx = context.WithValue(ctx, middleware.UserRoleKey, "therapist")
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})
		r.Post("/", handlers.CreateReview)
	})

	therapistID := uuid.New().String()
	body := map[string]interface{}{"therapistId": therapistID, "rating": 5, "comment": "great"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 403/401, got %d: %s", w.Code, w.Body.String())
	}
}
