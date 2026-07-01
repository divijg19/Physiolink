package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/db"
)

type mockReminderQueries struct {
	rows []db.GetUpcomingRemindersBeforeRow
	err  error
}

func (m *mockReminderQueries) GetUpcomingRemindersBefore(_ context.Context, _ db.GetUpcomingRemindersBeforeParams) ([]db.GetUpcomingRemindersBeforeRow, error) {
	return m.rows, m.err
}

func TestListForPatient_ReturnsReminders(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	apptStart := time.Date(2025, 6, 2, 10, 0, 0, 0, time.UTC)

	mockQ := &mockReminderQueries{
		rows: []db.GetUpcomingRemindersBeforeRow{
			{
				ID:               uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				AppointmentID:    uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				ScheduledFor:     time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC),
				Payload:          pqtype.NullRawMessage{Valid: false},
				TherapistID:      uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				AppointmentStart: apptStart,
			},
		},
	}

	svc := NewReminderService(mockQ, clock.NewFake(now))
	pid := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	items, err := svc.ListForPatient(context.Background(), pid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Message != "Reminder: appointment on 2025-06-02 10:00" {
		t.Fatalf("unexpected message: %q", items[0].Message)
	}
}

func TestListForPatient_UsesPayloadMessage(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	payload := pqtype.NullRawMessage{
		RawMessage: []byte(`{"message": "Custom reminder message"}`),
		Valid:      true,
	}

	mockQ := &mockReminderQueries{
		rows: []db.GetUpcomingRemindersBeforeRow{
			{
				ID:               uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				ScheduledFor:     now,
				Payload:          payload,
				AppointmentStart: now,
			},
		},
	}

	svc := NewReminderService(mockQ, clock.NewFake(now))
	items, err := svc.ListForPatient(context.Background(), uuid.MustParse("44444444-4444-4444-4444-444444444444"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if items[0].Message != "Custom reminder message" {
		t.Fatalf("expected custom message, got %q", items[0].Message)
	}
}

func TestListForPatient_QueryError(t *testing.T) {
	mockQ := &mockReminderQueries{err: errors.New("db error")}
	svc := NewReminderService(mockQ, clock.NewFake(time.Now()))

	_, err := svc.ListForPatient(context.Background(), uuid.MustParse("44444444-4444-4444-4444-444444444444"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListForPatient_EmptyRows(t *testing.T) {
	mockQ := &mockReminderQueries{rows: []db.GetUpcomingRemindersBeforeRow{}}
	svc := NewReminderService(mockQ, clock.NewFake(time.Now()))

	items, err := svc.ListForPatient(context.Background(), uuid.MustParse("44444444-4444-4444-4444-444444444444"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
