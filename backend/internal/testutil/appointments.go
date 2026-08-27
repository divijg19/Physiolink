package testutil

import (
	"context"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// CreateAvailability calls the AppointmentService to insert availability slots for a therapist.
func CreateAvailability(ctx context.Context, database *db.DB, therapistID uuid.UUID, slots []struct{ StartTs, EndTs string }) error {
	apptSvc := service.NewAppointmentService(database, nil)
	return apptSvc.CreateAvailability(ctx, therapistID, slots)
}

// BookFirstAvailableSlot finds the first open slot for a therapist and books it for the patient.
// Returns appointment ID and booked slot ID.
func BookFirstAvailableSlot(ctx context.Context, database *db.DB, therapistID, patientID uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	apptSvc := service.NewAppointmentService(database, nil)
	slots, err := apptSvc.GetTherapistAvailability(ctx, therapistID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	if len(slots) == 0 {
		return uuid.Nil, uuid.Nil, nil
	}
	slotID := slots[0].ID
	apptID, err := apptSvc.BookAppointment(ctx, slotID, patientID)
	return apptID, slotID, err
}
