package activities

import (
	"context"
	"testing"
)

func TestSendConfirmationEmail(t *testing.T) {
	result, err := SendConfirmationEmail(context.Background(), "appt-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Email sent" {
		t.Fatalf("expected 'Email sent', got %q", result)
	}
}
