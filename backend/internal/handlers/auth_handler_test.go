package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/handlers"
	mocks "github.com/divijg19/physiolink/backend/tests/__mocks__"
)

func TestRegister_ReturnsTokenAndProfile(t *testing.T) {
	cfg := config.New()
	handlers.InitAuth(mocks.NewAuthServiceMock(), cfg)
	// Patch profile service only for creating empty profile
	handlers.InitProfile(mocks.NewProfileServiceMock())

	body := map[string]string{"email": "a@b.com", "password": "pw", "role": "patient"}
	b, _ := json.Marshal(body)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(b))
	handlers.Register(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp struct {
		Token   string                 `json:"token"`
		Profile map[string]interface{} `json:"profile"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("expected token")
	}
	// sanity: token decodes
	if _, _, err := new(jwt.Parser).ParseUnverified(resp.Token, jwt.MapClaims{}); err != nil {
		t.Fatalf("token not parseable: %v", err)
	}
}
