package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/testutil"
)

func TestE2E_RegisterCreateAvailabilityBook(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	// start server in-process on an httptest server
	// ensure environment variables for config are set
	os.Setenv("DATABASE_URL", dbURL)
	os.Setenv("JWT_SECRET", "testsecret")
	cfg := config.New()

	// Connect test database and initialize services with a deterministic fake clock
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	database, err := db.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("db connect failed: %v", err)
	}
	defer database.Close()

	// fixed time for deterministic reminders
	fixed := time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC)
	clk := clock.NewFake(fixed)

	handler := testutil.NewRouterWithServices(cfg, database, clk)
	ts := httptest.NewServer(handler)
	defer ts.Close()
	client := ts.Client()
	// use ts.URL as base
	base := ts.URL

	// helper to register a user
	register := func(email, password, role string) (token string) {
		body := map[string]string{"email": email, "password": password, "role": role}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, base+"/api/auth/register", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("register request failed: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("register status: %d", resp.StatusCode)
		}
		var out struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("decode register: %v", err)
		}
		return out.Token
	}

	// register therapist and patient
	thToken := register("therapist@example.com", "pass1234", "therapist")
	ptToken := register("patient@example.com", "pass1234", "patient")

	// decode therapist token to get id
	parsed, _ := jwt.Parse(thToken, func(token *jwt.Token) (interface{}, error) { return []byte("testsecret"), nil })
	claims := parsed.Claims.(jwt.MapClaims)
	user := claims["user"].(map[string]interface{})
	thID := user["id"].(string)

	// therapist creates availability
	slotBody := map[string]interface{}{"slots": []map[string]string{{"startTime": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339), "endTime": time.Now().Add(24*time.Hour + 30*time.Minute).UTC().Format(time.RFC3339)}}}
	sb, _ := json.Marshal(slotBody)
	req, _ := http.NewRequest(http.MethodPost, base+"/api/appointments/availability", bytes.NewReader(sb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+thToken)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create availability failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	// patient queries availability
	req2, _ := http.NewRequest(http.MethodGet, base+"/api/appointments/availability?ptId="+thID, nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("get availability failed: %v", err)
	}
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
	var slots []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&slots); err != nil {
		t.Fatalf("decode slots: %v", err)
	}
	resp2.Body.Close()
	if len(slots) == 0 {
		t.Fatalf("no slots returned")
	}

	// patient books first slot
	slotID := slots[0].ID
	req3, _ := http.NewRequest(http.MethodPut, base+"/api/appointments/"+slotID+"/book", nil)
	req3.Header.Set("Authorization", "Bearer "+ptToken)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("book failed: %v", err)
	}
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 booking, got %d", resp3.StatusCode)
	}
	var bookResp map[string]string
	if err := json.NewDecoder(resp3.Body).Decode(&bookResp); err != nil {
		t.Fatalf("decode book resp: %v", err)
	}
	resp3.Body.Close()
	if _, err := uuid.Parse(bookResp["id"]); err != nil {
		t.Fatalf("invalid appointment id: %v", err)
	}
	apptID := bookResp["id"]

	// therapist confirms the appointment (status -> confirmed)
	statusBody := map[string]string{"status": "confirmed"}
	sb2, _ := json.Marshal(statusBody)
	req4, _ := http.NewRequest(http.MethodPut, base+"/api/appointments/"+apptID+"/status", bytes.NewReader(sb2))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Authorization", "Bearer "+thToken)
	resp4, err := client.Do(req4)
	if err != nil {
		t.Fatalf("confirm failed: %v", err)
	}
	if resp4.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 confirm, got %d", resp4.StatusCode)
	}
	resp4.Body.Close()

	// patient checks reminders
	req5, _ := http.NewRequest(http.MethodGet, base+"/api/reminders/me", nil)
	req5.Header.Set("Authorization", "Bearer "+ptToken)
	resp5, err := client.Do(req5)
	if err != nil {
		t.Fatalf("get reminders failed: %v", err)
	}
	if resp5.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 reminders, got %d", resp5.StatusCode)
	}
	var rems []struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(resp5.Body).Decode(&rems); err != nil {
		t.Fatalf("decode reminders: %v", err)
	}
	resp5.Body.Close()
	if len(rems) == 0 {
		t.Fatalf("expected at least 1 reminder after confirm")
	}
}
