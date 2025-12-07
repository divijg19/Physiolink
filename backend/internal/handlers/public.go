package handlers

import (
	"net/http"

	"github.com/divijg19/physiolink/backend/internal/middleware"
	"github.com/divijg19/physiolink/backend/internal/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(string)
	views.Home("Physiolink", ok).Render(r.Context(), w)
}
