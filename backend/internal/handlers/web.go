package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/divijg19/physiolink/backend/internal/views"
)

// Render views
func LoginPage(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	views.Register().Render(r.Context(), w)
}

func TherapistsPage(w http.ResponseWriter, r *http.Request) {
	// Fetch therapists
	// We can reuse GetAllTherapists logic or call service directly
	// Calling service directly is better
	result, err := therapistService.GetAllTherapists(r.Context(), service.TherapistQueryParams{
		Page:  1,
		Limit: 20,
	})
	if err != nil {
		http.Error(w, "Failed to fetch therapists", http.StatusInternalServerError)
		return
	}

	_, ok := r.Context().Value(middleware.UserIDKey).(string)
	views.TherapistsList(result.Data, ok).Render(r.Context(), w)
}

func TherapistDetailPage(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	tData, err := therapistService.GetTherapistByID(r.Context(), idStr, "")
	if err != nil {
		http.Error(w, "Therapist not found", http.StatusNotFound)
		return
	}

	tid, _ := uuid.Parse(idStr)
	slots, err := apptService.GetTherapistAvailability(r.Context(), tid)
	if err != nil {
		slots = []service.Slot{}
	}

	profile, _ := tData["profile"].(map[string]interface{})
	firstName, _ := profile["firstName"].(string)
	lastName, _ := profile["lastName"].(string)
	specialty, _ := profile["specialty"].(string)
	bio, _ := profile["bio"].(string)
	email, _ := tData["email"].(string)

	var slotViews []views.SlotView
	for _, s := range slots {
		start, _ := time.Parse(time.RFC3339, s.StartTs)
		slotViews = append(slotViews, views.SlotView{
			ID:        s.ID.String(),
			StartTime: start,
			IsBooked:  s.Status == "booked",
		})
	}

	detail := views.TherapistDetailView{
		ID:        idStr,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Specialty: specialty,
		Bio:       bio,
		Slots:     slotViews,
	}

	_, isLoggedIn := r.Context().Value(middleware.UserIDKey).(string)
	views.TherapistDetail(detail, isLoggedIn).Render(r.Context(), w)
}

func BookAppointmentWeb(w http.ResponseWriter, r *http.Request) {
	slotIDStr := chi.URLParam(r, "id")
	slotID, _ := uuid.Parse(slotIDStr)
	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	_, err := apptService.BookAppointment(r.Context(), slotID, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error booking"))
		return
	}

	w.Write([]byte(`<button class="w-full bg-green-600 text-white py-2 px-4 rounded cursor-default text-sm font-medium" disabled>Booked</button>`))
}

func GetReviewsWeb(w http.ResponseWriter, r *http.Request) {
	tidStr := chi.URLParam(r, "therapistId")
	tid, _ := uuid.Parse(tidStr)

	rawReviews, err := reviewService.GetReviewsForTherapist(r.Context(), tid)
	if err != nil {
		rawReviews = []map[string]interface{}{}
	}

	var reviewViews []views.ReviewView
	for _, r := range rawReviews {
		patient, _ := r["patient"].(map[string]interface{})
		profile, _ := patient["profile"].(map[string]interface{})
		firstName, _ := profile["firstName"].(string)

		rating, _ := r["rating"].(float64) // JSON numbers are float64
		comment, _ := r["comment"].(string)
		createdAtStr, _ := r["createdAt"].(string)
		createdAt, _ := time.Parse(time.RFC3339, createdAtStr)

		reviewViews = append(reviewViews, views.ReviewView{
			ID:          r["_id"].(string),
			PatientName: firstName,
			Rating:      int(rating),
			Comment:     comment,
			CreatedAt:   createdAt,
		})
	}

	_, isLoggedIn := r.Context().Value(middleware.UserIDKey).(string)
	views.ReviewsList(reviewViews, tidStr, isLoggedIn).Render(r.Context(), w)
}

func GetReviewFormWeb(w http.ResponseWriter, r *http.Request) {
	tidStr := chi.URLParam(r, "therapistId")
	views.ReviewForm(tidStr).Render(r.Context(), w)
}

func PostReviewWeb(w http.ResponseWriter, r *http.Request) {
	tidStr := chi.URLParam(r, "therapistId")
	tid, _ := uuid.Parse(tidStr)

	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	ratingStr := r.FormValue("rating")
	comment := r.FormValue("comment")

	// Simple conversion
	rating := 5
	if ratingStr == "1" {
		rating = 1
	}
	if ratingStr == "2" {
		rating = 2
	}
	if ratingStr == "3" {
		rating = 3
	}
	if ratingStr == "4" {
		rating = 4
	}

	_, err := reviewService.CreateReview(r.Context(), userID, tid, rating, comment)
	if err != nil {
		// In a real app, return the form with error
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("<div class='bg-red-100 text-red-700 p-4 rounded mb-4'>Error: %s</div>", err.Error())))
		return
	}

	// Return the updated list
	GetReviewsWeb(w, r)
}

func GetProfileFormWeb(w http.ResponseWriter, r *http.Request) {
	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)

	// In a real app, fetch current profile data to pre-fill
	// For now, just render the form
	views.ProfileForm(userIDStr, role).Render(r.Context(), w)
}

func GetProfileWeb(w http.ResponseWriter, r *http.Request) {
	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)
	views.ProfileView(userIDStr, role).Render(r.Context(), w)
}

func PutProfileWeb(w http.ResponseWriter, r *http.Request) {
	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	bio := r.FormValue("bio")

	update := service.NodeProfileUpdate{
		FirstName: firstName,
		LastName:  lastName,
		Bio:       bio,
	}

	_, err := profileService.UpsertProfile(r.Context(), userID, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error updating profile"))
		return
	}

	// Return the updated profile view
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)
	views.ProfileView(userIDStr, role).Render(r.Context(), w)
}

// Form handlers
func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, token, err := authService.Authenticate(r.Context(), email, password)
	if err != nil {
		// In a real HTMX app, we'd return a partial with the error message
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("<div class='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative' role='alert'><strong class='font-bold'>Error!</strong> <span class='block sm:inline'>" + err.Error() + "</span></div>"))
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// HTMX redirect via header
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	_, token, err := authService.Register(r.Context(), email, password, role)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative' role='alert'><strong class='font-bold'>Error!</strong> <span class='block sm:inline'>" + err.Error() + "</span></div>"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
