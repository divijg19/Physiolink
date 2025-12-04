package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	"github.com/divijg19/physiolink/backend/go/internal/db"
	"github.com/divijg19/physiolink/backend/go/internal/workflows"
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
		_, err := s.db.Pool.Exec(ctx, `INSERT INTO availability_slots (therapist_id, start_ts, end_ts, status) VALUES ($1,$2,$3,'open') ON CONFLICT DO NOTHING`, therapistID, sl.StartTs, sl.EndTs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *AppointmentService) GetTherapistAvailability(ctx context.Context, therapistID uuid.UUID) ([]Slot, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT id, therapist_id, start_ts, end_ts, status FROM availability_slots WHERE therapist_id=$1 AND status='open' ORDER BY start_ts ASC`, therapistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Slot
	for rows.Next() {
		var sl Slot
		if err := rows.Scan(&sl.ID, &sl.TherapistID, &sl.StartTs, &sl.EndTs, &sl.Status); err != nil {
			return nil, err
		}
		out = append(out, sl)
	}
	return out, nil
}

func (s *AppointmentService) BookAppointment(ctx context.Context, appointmentID uuid.UUID, patientID uuid.UUID) (uuid.UUID, error) {
	// Start a transaction -- FOR UPDATE requires a transaction to take row locks.
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() {
		// if not committed, rollback to free locks
		_ = tx.Rollback(ctx)
	}()

	// lock the slot within the transaction
	var slotID uuid.UUID
	var therapistID uuid.UUID
	var status string
	err = tx.QueryRow(ctx, `SELECT id, therapist_id, status FROM availability_slots WHERE id=$1 FOR UPDATE`, appointmentID).Scan(&slotID, &therapistID, &status)
	if err != nil {
		return uuid.Nil, err
	}
	if status != "open" {
		return uuid.Nil, ErrConflict
	}

	// insert appointment
	var apptID uuid.UUID
	err = tx.QueryRow(ctx, `INSERT INTO appointments (slot_id, patient_id, therapist_id, status) VALUES ($1,$2,$3,'booked') RETURNING id`, slotID, patientID, therapistID).Scan(&apptID)
	if err != nil {
		return uuid.Nil, err
	}

	// mark slot reserved
	if _, err := tx.Exec(ctx, `UPDATE availability_slots SET status='reserved' WHERE id=$1`, slotID); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	// Kick off a Temporal workflow for side effects (non-blocking).
	if s.tcl != nil {
		_, _ = s.tcl.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
			TaskQueue: "appointment-task-queue",
		}, workflows.BookingWorkflow, workflows.BookingWorkflowParam{
			AppointmentID: apptID.String(),
			PatientID:     patientID.String(),
			TherapistID:   therapistID.String(),
		})
	}

	return apptID, nil
}

func (s *AppointmentService) GetMySchedule(ctx context.Context, userID uuid.UUID, role string) ([]uuid.UUID, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT id FROM appointments WHERE CASE WHEN $2='pt' THEN therapist_id=$1 ELSE patient_id=$1 END ORDER BY created_at ASC`, userID, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
