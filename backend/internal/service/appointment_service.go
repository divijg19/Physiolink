package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go.temporal.io/sdk/client"

	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/workflows"
)

// ErrConflict is returned when attempting to book a slot that is no longer open.
var ErrConflict = errors.New("conflict")

type AppointmentService struct {
	db  *db.DB
	tcl client.Client
}

func NewAppointmentService(d *db.DB, tcl client.Client) *AppointmentService {
	return &AppointmentService{db: d, tcl: tcl}
}

type Slot struct {
	ID          uuid.UUID
	TherapistID uuid.UUID
	StartTs     string
	EndTs       string
	Status      string
}

func (s *AppointmentService) CreateAvailability(ctx context.Context, therapistID uuid.UUID, slots []struct{ StartTs, EndTs string }) error {
	for _, sl := range slots {
		parsed, err := time.Parse(time.RFC3339, sl.StartTs)
		if err != nil {
			return err
		}
		parsedEnd, err := time.Parse(time.RFC3339, sl.EndTs)
		if err != nil {
			return err
		}
		arg := db.CreateAvailabilitySlotsParams{
			TherapistID: therapistID,
			StartTs:     parsed,
			EndTs:       parsedEnd,
		}
		if err := s.db.Queries.CreateAvailabilitySlots(ctx, arg); err != nil {
			return err
		}
	}
	return nil
}

func (s *AppointmentService) GetTherapistAvailability(ctx context.Context, therapistID uuid.UUID) ([]Slot, error) {
	rows, err := s.db.Queries.GetTherapistOpenSlots(ctx, therapistID)
	if err != nil {
		return nil, err
	}
	var out []Slot
	for _, r := range rows {
		out = append(out, Slot{
			ID:          r.ID,
			TherapistID: r.TherapistID,
			StartTs:     r.StartTs.Format(time.RFC3339),
			EndTs:       r.EndTs.Format(time.RFC3339),
			Status:      r.Status,
		})
	}
	return out, nil
}

func (s *AppointmentService) BookAppointment(ctx context.Context, appointmentID uuid.UUID, patientID uuid.UUID) (uuid.UUID, error) {
	// Start a transaction using database/sql
	tx, err := s.db.SQL.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// use sqlc queries with transaction
	qtx := s.db.Queries.WithTx(tx)

	// lock the slot within the transaction
	slot, err := qtx.BookAppointmentTxLockSlot(ctx, appointmentID)
	if err != nil {
		return uuid.Nil, err
	}
	if slot.Status != "open" {
		return uuid.Nil, ErrConflict
	}

	// insert appointment
	apptID, err := qtx.InsertAppointment(ctx, db.InsertAppointmentParams{
		SlotID:      uuid.NullUUID{UUID: slot.ID, Valid: true},
		PatientID:   patientID,
		TherapistID: slot.TherapistID,
		Status:      "booked",
		Notes:       sql.NullString{Valid: false},
	})
	if err != nil {
		return uuid.Nil, err
	}

	// mark slot reserved
	if err := qtx.UpdateSlotStatus(ctx, db.UpdateSlotStatusParams{
		ID:     slot.ID,
		Status: "reserved",
	}); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	// Kick off a Temporal workflow for side effects (non-blocking).
	if s.tcl != nil {
		_, _ = s.tcl.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
			TaskQueue: "appointment-task-queue",
		}, workflows.BookingWorkflow, workflows.BookingWorkflowParam{
			AppointmentID: apptID.String(),
			PatientID:     patientID.String(),
			TherapistID:   slot.TherapistID.String(),
		})
	}

	return apptID, nil
}

func (s *AppointmentService) ListMyAppointments(ctx context.Context, userID uuid.UUID, role string) ([]AppointmentBrief, error) {
	rows, err := s.db.Queries.ListMyAppointmentsWithDetails(ctx, db.ListMyAppointmentsWithDetailsParams{
		TherapistID: userID,
		Column2:     role,
	})
	if err != nil {
		return nil, err
	}

	var out []AppointmentBrief
	for _, r := range rows {
		// derive names
		ptFirst, ptLast := splitDisplayName(r.PtDisplayName.String)
		paFirst, paLast := splitDisplayName(r.PaDisplayName.String)

		a := AppointmentBrief{
			ID:     r.ID.String(),
			Start:  r.StartTs.Format(time.RFC3339),
			End:    r.EndTs.Format(time.RFC3339),
			Status: r.Status,
			PT: map[string]interface{}{
				"_id": r.TherapistID.String(),
				"profile": map[string]interface{}{
					"firstName": ptFirst,
					"lastName":  ptLast,
				},
			},
			Patient: map[string]interface{}{
				"_id": r.PatientID.String(),
				"profile": map[string]interface{}{
					"firstName": paFirst,
					"lastName":  paLast,
				},
			},
		}
		out = append(out, a)
	}
	return out, nil
}

type AppointmentBrief struct {
	ID      string                 `json:"_id"`
	PT      map[string]interface{} `json:"pt"`
	Patient map[string]interface{} `json:"patient"`
	Start   string                 `json:"startTime"`
	End     string                 `json:"endTime"`
	Status  string                 `json:"status"`
}

func splitDisplayName(s string) (string, string) {
	if s == "" {
		return "", ""
	}
	for i, r := range s {
		if r == ' ' {
			if i+1 < len(s) {
				return s[:i], s[i+1:]
			}
			return s[:i], ""
		}
	}
	return s, ""
}

func (s *AppointmentService) UpdateAppointmentStatus(ctx context.Context, appointmentID uuid.UUID, ptID uuid.UUID, status string) (AppointmentBrief, error) {
	var out AppointmentBrief
	if status != "confirmed" && status != "rejected" {
		return out, errors.New("invalid status")
	}

	// ensure appointment exists and owned by pt
	therapistID, err := s.db.Queries.GetAppointmentTherapistID(ctx, appointmentID)
	if err != nil {
		return out, err
	}
	if therapistID != ptID {
		return out, errors.New("forbidden")
	}

	// update status
	if err := s.db.Queries.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
		ID:     appointmentID,
		Status: status,
	}); err != nil {
		return out, err
	}

	// create reminder when confirmed: schedule 24h before start
	if status == "confirmed" {
		slotInfo, err := s.db.Queries.GetAppointmentSlotStartTime(ctx, appointmentID)
		if err == nil {
			payload := map[string]interface{}{"message": "Reminder: appointment on " + slotInfo.StartTs.Format(time.RFC3339)}
			b, _ := json.Marshal(payload)
			scheduledFor := slotInfo.StartTs.Add(-24 * time.Hour)
			_ = s.db.Queries.InsertReminder(ctx, db.InsertReminderParams{
				AppointmentID: appointmentID,
				ScheduledFor:  scheduledFor,
				Payload:       pqtype.NullRawMessage{RawMessage: b, Valid: len(b) > 0},
			})
		}
	}

	// return populated brief
	list, err := s.ListMyAppointments(ctx, ptID, "pt")
	if err != nil {
		return out, err
	}
	for _, a := range list {
		if a.ID == appointmentID.String() {
			return a, nil
		}
	}
	return out, errors.New("not found")
}
