package main

import (
	"context"
	"log"
	"time"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/server"
	"github.com/divijg19/physiolink/backend/internal/service"
)

func main() {
	cfg := config.New()
	// connect DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	database, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer database.Close()

	// init services
	authSvc := service.NewAuthService(database, cfg)
	profileSvc := service.NewProfileService(database, cfg)
	therapistSvc := service.NewTherapistService(database)
	reviewSvc := service.NewReviewService(database)
	reminderSvc := service.NewReminderService(database)
	// temporal client (optional in dev)
	tcl, err := service.NewTemporalClient()
	if err != nil {
		log.Printf("temporal client init failed: %v", err)
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
	log.Printf("starting server on %s", cfg.BindAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
