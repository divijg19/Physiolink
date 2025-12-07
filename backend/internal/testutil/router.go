package testutil

import (
	"net/http"

	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/handlers"
	"github.com/divijg19/physiolink/backend/internal/server"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// NewRouterWithServices wires up database-backed services and handlers using the provided clock.
// Caller is responsible for connecting the DB and closing it when done.
func NewRouterWithServices(cfg *config.Config, database *db.DB, clk clock.Clock) http.Handler {
	// create services
	authSvc := service.NewAuthService(database, cfg)
	profileSvc := service.NewProfileService(database, cfg)
	therapistSvc := service.NewTherapistService(database)
	reviewSvc := service.NewReviewService(database)
	reminderSvc := service.NewReminderService(database, clk)
	apptSvc := service.NewAppointmentService(database, nil)

	// register handlers
	handlers.InitAuth(authSvc, cfg)
	handlers.InitProfile(profileSvc)
	handlers.InitTherapists(therapistSvc)
	handlers.InitReviews(reviewSvc)
	handlers.InitAppointments(apptSvc)
	handlers.InitReminders(reminderSvc)

	return server.NewRouter(cfg)
}
