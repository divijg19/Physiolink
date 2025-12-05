package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
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

type AppointmentBrief struct {
	ID      string                 `json:"_id"`
	PT      map[string]interface{} `json:"pt"`
	Patient map[string]interface{} `json:"patient"`
	Start   string                 `json:"startTime"`
	End     string                 `json:"endTime"`
	Status  string                 `json:"status"`
}

func (s *AppointmentService) ListMyAppointments(ctx context.Context, userID uuid.UUID, role string) ([]AppointmentBrief, error) {
	q := `SELECT a.id::text, a.therapist_id, a.patient_id, a.status, s.start_ts, s.end_ts,
				 p_pt.display_name, p_pt.profile_extra,
				 p_pa.display_name, p_pa.profile_extra
		  FROM appointments a
		  JOIN availability_slots s ON s.id = a.slot_id
		  LEFT JOIN profiles p_pt ON p_pt.user_id = a.therapist_id
		  LEFT JOIN profiles p_pa ON p_pa.user_id = a.patient_id
		  WHERE CASE WHEN $2='pt' THEN a.therapist_id=$1 ELSE a.patient_id=$1 END
		  ORDER BY s.start_ts ASC`
	rows, err := s.db.Pool.Query(ctx, q, userID, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AppointmentBrief
	for rows.Next() {
		var id string
		var ptID, paID uuid.UUID
		var status string
		var startTs, endTs string
		var ptDisplay, paDisplay string
		var ptExtra, paExtra []byte
		if err := rows.Scan(&id, &ptID, &paID, &status, &startTs, &endTs, &ptDisplay, &ptExtra, &paDisplay, &paExtra); err != nil {
			return nil, err
		}
		// derive names
		ptFirst, ptLast := splitDisplayName(ptDisplay)
		paFirst, paLast := splitDisplayName(paDisplay)

		a := AppointmentBrief{
			ID:     id,
			Start:  startTs,
			End:    endTs,
			Status: status,
			PT: map[string]interface{}{
				"_id": ptID.String(),
				"profile": map[string]interface{}{
					"firstName": ptFirst,
					"lastName":  ptLast,
				},
			},
			Patient: map[string]interface{}{
				"_id": paID.String(),
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
	var therapistID uuid.UUID
	err := s.db.Pool.QueryRow(ctx, `SELECT therapist_id FROM appointments WHERE id=$1`, appointmentID).Scan(&therapistID)
	if err != nil {
		return out, err
	}
	if therapistID != ptID {
		return out, errors.New("forbidden")
	}

	// update status
	if _, err := s.db.Pool.Exec(ctx, `UPDATE appointments SET status=$2, updated_at=now() WHERE id=$1`, appointmentID, status); err != nil {
		return out, err
	}

	// create reminder when confirmed: schedule 24h before start
	if status == "confirmed" {
		var slotID uuid.UUID
		var startTs string
		err := s.db.Pool.QueryRow(ctx, `SELECT a.slot_id, s.start_ts FROM appointments a JOIN availability_slots s ON s.id=a.slot_id WHERE a.id=$1`, appointmentID).Scan(&slotID, &startTs)
		if err == nil {
			// Insert reminder with payload message
			payload := map[string]interface{}{"message": "Reminder: appointment on " + startTs}
			b, _ := json.Marshal(payload)
			// schedule for 24h before
			_, _ = s.db.Pool.Exec(ctx, `INSERT INTO reminders (appointment_id, scheduled_for, payload) VALUES ($1, $2::timestamptz - INTERVAL '24 hours', $3)`, appointmentID, startTs, b)
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
