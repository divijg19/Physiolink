package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	stdmw "github.com/go-chi/chi/v5/middleware"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/handlers"
	mware "github.com/divijg19/physiolink/backend/internal/middleware"
)

type Server struct {
	httpServer *http.Server
}

// NewRouter builds and returns an http.Handler configured with all routes.
func NewRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(stdmw.RequestID)
	r.Use(stdmw.RealIP)
	r.Use(stdmw.Logger)
	r.Use(stdmw.Recoverer)

	// health
	r.Get("/health", handlers.Health)

	// Public routes with optional auth (for navbar state)
	r.Group(func(r chi.Router) {
		r.Use(mware.OptionalCookieAuth(cfg))
		r.Get("/", handlers.Home)
		r.Get("/login", handlers.LoginPage)
		r.Get("/register", handlers.RegisterPage)
		r.Get("/therapists", handlers.TherapistsPage)
		r.Get("/therapists/{id}", handlers.TherapistDetailPage)
		r.Get("/web/reviews/{therapistId}", handlers.GetReviewsWeb)
	})

	// Protected routes (Web)
	r.Group(func(r chi.Router) {
		r.Use(mware.CookieAuth(cfg))
		r.Get("/dashboard", handlers.DashboardPage)
		r.Get("/dashboard/appointments", handlers.DashboardAppointments)
		r.Put("/web/appointments/{id}/book", handlers.BookAppointmentWeb)
		r.Get("/web/reviews/{therapistId}/form", handlers.GetReviewFormWeb)
		r.Post("/web/reviews/{therapistId}", handlers.PostReviewWeb)
		r.Get("/web/profile", handlers.GetProfileWeb)
		r.Get("/web/profile/edit", handlers.GetProfileFormWeb)
		r.Put("/web/profile", handlers.PutProfileWeb)
		r.Post("/auth/logout", handlers.Logout)
	})

	// HTMX Form submissions
	r.Post("/auth/login-form", handlers.LoginSubmit)
	r.Post("/auth/register-form", handlers.RegisterSubmit)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// auth
		r.Post("/auth/register", handlers.Register)
		r.Post("/auth/login", handlers.Login)

		// therapists (private)
		r.Group(func(r chi.Router) {
			r.Use(mware.JWTAuth(cfg))
			r.Get("/therapists", handlers.GetAllTherapists)
			r.Get("/therapists/{id}", handlers.GetTherapistByID)
		})

		// reviews (private)
		r.Group(func(r chi.Router) {
			r.Use(mware.JWTAuth(cfg))
			r.Post("/reviews", handlers.CreateReview)
			r.Get("/reviews/{therapistId}", handlers.GetReviewsForTherapist)
		})

		// profile
		r.Group(func(r chi.Router) {
			r.Use(mware.JWTAuth(cfg))
			r.Get("/profile/me", handlers.GetMyProfile)
			r.Put("/profile/me", handlers.UpsertMyProfile)
			r.Post("/profile", handlers.UpsertMyProfile)

			// appointments
			r.Post("/appointments/availability", handlers.CreateAvailability)
			r.Get("/appointments/availability", handlers.GetTherapistAvailability)        // supports ptId query
			r.Get("/appointments/availability/{ptId}", handlers.GetTherapistAvailability) // supports path param
			r.Get("/appointments/me", handlers.GetMyAppointments)
			r.Put("/appointments/{id}/book", handlers.BookAppointment)
			r.Put("/appointments/{id}/status", handlers.UpdateAppointmentStatus)
		})

		// reminders (private)
		r.Group(func(r chi.Router) {
			r.Use(mware.JWTAuth(cfg))
			r.Get("/reminders/me", handlers.GetMyReminders)
		})
	})

	return r
}

// New returns a Server that wraps the configured router and listens on cfg.BindAddr.
func New(cfg *config.Config) *Server {
	srv := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: NewRouter(cfg),
	}
	return &Server{httpServer: srv}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
