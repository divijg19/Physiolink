package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/db"
)

type ReminderService struct {
	db *db.DB
}

func NewReminderService(d *db.DB) *ReminderService { return &ReminderService{db: d} }

type ReminderItem struct {
	ID       string `json:"_id"`
	Message  string `json:"message"`
	RemindAt string `json:"remindAt"`
}

func (s *ReminderService) ListForPatient(ctx context.Context, patientID uuid.UUID) ([]ReminderItem, error) {
	q := `SELECT r.id::text, r.scheduled_for, r.payload, s.start_ts
	      FROM reminders r
	      JOIN appointments a ON a.id = r.appointment_id
	      JOIN availability_slots s ON s.id = a.slot_id
	      WHERE a.patient_id = $1
	      ORDER BY r.scheduled_for ASC`
	rows, err := s.db.Pool.Query(ctx, q, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ReminderItem
	for rows.Next() {
		var id string
		var scheduled, startTs string
		var payload []byte
		if err := rows.Scan(&id, &scheduled, &payload, &startTs); err != nil {
			return nil, err
		}
		msg := "Reminder: appointment on " + startTs
		if len(payload) > 0 {
			var m map[string]interface{}
			_ = json.Unmarshal(payload, &m)
			if v, ok := m["message"].(string); ok && v != "" {
				msg = v
			}
		}
		out = append(out, ReminderItem{ID: id, Message: msg, RemindAt: scheduled})
	}
	return out, nil
}
