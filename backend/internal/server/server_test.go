package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/divijg19/physiolink/backend/internal/config"
)

func TestNewRouter_ReturnsHandler(t *testing.T) {
	cfg := &config.Config{BindAddr: ":8080"}
	handler := NewRouter(cfg)
	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{BindAddr: ":8080"}
	handler := NewRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestNew_SetsTimeouts(t *testing.T) {
	cfg := &config.Config{BindAddr: ":0"}
	s := New(cfg)
	if s == nil {
		t.Fatal("expected non-nil server")
	}
	if s.httpServer.ReadHeaderTimeout != 10*time.Second {
		t.Fatalf("expected ReadHeaderTimeout 10s, got %v", s.httpServer.ReadHeaderTimeout)
	}
	if s.httpServer.ReadTimeout != 30*time.Second {
		t.Fatalf("expected ReadTimeout 30s, got %v", s.httpServer.ReadTimeout)
	}
	if s.httpServer.WriteTimeout != 30*time.Second {
		t.Fatalf("expected WriteTimeout 30s, got %v", s.httpServer.WriteTimeout)
	}
	if s.httpServer.IdleTimeout != 60*time.Second {
		t.Fatalf("expected IdleTimeout 60s, got %v", s.httpServer.IdleTimeout)
	}
}

func TestListenAndServe_Shutdown(t *testing.T) {
	cfg := &config.Config{BindAddr: ":0"}
	s := New(cfg)

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.ListenAndServe()
	}()

	time.Sleep(50 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("unexpected ListenAndServe error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for server to stop")
	}
}
