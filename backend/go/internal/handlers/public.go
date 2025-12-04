package handlers

import (
	"net/http"

	"github.com/divijg19/physiolink/backend/go/internal/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	views.Home("Physiolink").Render(r.Context(), w)
}
