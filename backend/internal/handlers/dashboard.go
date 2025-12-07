package handlers

import (
	"net/http"
	"time"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/views"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {
	// Middleware ensures these are present
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)

	// In a real app, we'd fetch user details from DB using userID
	// For now, we'll just display the ID and Role

	views.Dashboard(userID, role).Render(r.Context(), w)
}

func DashboardAppointments(w http.ResponseWriter, r *http.Request) {
	// Simulate fetching from DB
	// In a real app, use the service layer: service.GetMyAppointments(ctx, userID)
	time.Sleep(500 * time.Millisecond) // Simulate network delay to show loading state

	appointments := []views.AppointmentView{
		{
			ID:        "1",
			OtherName: "Dr. Smith",
			StartTime: time.Now().Add(24 * time.Hour),
			Status:    "confirmed",
		},
		{
			ID:        "2",
			OtherName: "Dr. Jones",
			StartTime: time.Now().Add(48 * time.Hour),
			Status:    "pending",
		},
	}

	views.AppointmentsList(appointments).Render(r.Context(), w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
