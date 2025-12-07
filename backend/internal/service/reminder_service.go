package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/clock"
	"github.com/divijg19/physiolink/backend/internal/db"
)

type ReminderService struct {
	db  *db.DB
	clk clock.Clock
}

func NewReminderService(d *db.DB, clk clock.Clock) *ReminderService {
	return &ReminderService{db: d, clk: clk}
}

type ReminderItem struct {
	ID       string `json:"_id"`
	Message  string `json:"message"`
	RemindAt string `json:"remindAt"`
}

func (s *ReminderService) ListForPatient(ctx context.Context, patientID uuid.UUID) ([]ReminderItem, error) {
	// use parameterized query so tests can control "now"
	rows, err := s.db.Queries.GetUpcomingRemindersBefore(ctx, db.GetUpcomingRemindersBeforeParams{PatientID: patientID, ScheduledFor: s.clk.Now()})
	if err != nil {
		return nil, err
	}
	var out []ReminderItem
	for _, r := range rows {
		msg := "Reminder: appointment on " + r.AppointmentStart.Format("2006-01-02 15:04")
		if r.Payload.Valid && len(r.Payload.RawMessage) > 0 {
			var m map[string]interface{}
			_ = json.Unmarshal(r.Payload.RawMessage, &m)
			if v, ok := m["message"].(string); ok && v != "" {
				msg = v
			}
		}
		out = append(out, ReminderItem{
			ID:       r.ID.String(),
			Message:  msg,
			RemindAt: r.ScheduledFor.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return out, nil
}
