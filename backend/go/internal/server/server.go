package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	stdmw "github.com/go-chi/chi/v5/middleware"

	"github.com/divijg19/physiolink/backend/go/internal/config"
	"github.com/divijg19/physiolink/backend/go/internal/handlers"
	mware "github.com/divijg19/physiolink/backend/go/internal/middleware"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config) *Server {
	r := chi.NewRouter()
	r.Use(stdmw.RequestID)
	r.Use(stdmw.RealIP)
	r.Use(stdmw.Logger)
	r.Use(stdmw.Recoverer)

	// health
	r.Get("/health", handlers.Health)

	// public (Templ SSR)
	r.Get("/", handlers.Home)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// auth
		r.Post("/auth/register", handlers.Register)
		r.Post("/auth/login", handlers.Login)

		// profile
		r.Group(func(r chi.Router) {
			r.Use(mware.JWTAuth(cfg))
			r.Get("/profile/me", handlers.GetMyProfile)
			r.Put("/profile/me", handlers.UpsertMyProfile)

			// appointments
			r.Post("/appointments/availability", handlers.CreateAvailability)
			r.Get("/appointments/availability", handlers.GetTherapistAvailability) // expects ptId query param
			r.Put("/appointments/{id}/book", handlers.BookAppointment)
		})
	})

	srv := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: r,
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
