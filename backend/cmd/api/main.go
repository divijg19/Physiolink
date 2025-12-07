package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/server"
	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Try loading .env from current directory or parent directory (for monorepo structure)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	cfg := config.New()
	// connect DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database, err := db.Connect(ctx, cfg)
	if err != nil {
		slog.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// init services
	authSvc := service.NewAuthService(database, cfg)
	profileSvc := service.NewProfileService(database, cfg)
	therapistSvc := service.NewTherapistService(database)
	reviewSvc := service.NewReviewService(database)
	reminderSvc := service.NewReminderService(database, clock.NewReal())
	// temporal client (optional in dev)
	tcl, err := service.NewTemporalClient()
	if err != nil {
		slog.Warn("temporal client init failed", "error", err)
	}
	defer func() {
		if tcl != nil {
			tcl.Close()
		}
	}()

	apptSvc := service.NewAppointmentService(database, tcl)

	// init handlers
	handlers.InitAuth(authSvc, cfg)
	handlers.InitProfile(profileSvc)
	handlers.InitTherapists(therapistSvc)
	handlers.InitReviews(reviewSvc)
	handlers.InitAppointments(apptSvc)
	handlers.InitReminders(reminderSvc)

	srv := server.New(cfg)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("starting server", "addr", cfg.BindAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
