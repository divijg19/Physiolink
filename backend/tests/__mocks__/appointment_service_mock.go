package __mocks__

import (
	"context"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/service"
)

type AppointmentServiceMock struct {
	CreateErr  error
	SlotsResp  []service.Slot
	BookResp   uuid.UUID
	BookErr    error
	ListResp   []service.AppointmentBrief
	ListErr    error
	UpdateResp service.AppointmentBrief
	UpdateErr  error
}

func (m *AppointmentServiceMock) CreateAvailability(ctx context.Context, therapistID uuid.UUID, slots []struct{ StartTs, EndTs string }) error {
	return m.CreateErr
}

func (m *AppointmentServiceMock) GetTherapistAvailability(ctx context.Context, therapistID uuid.UUID) ([]service.Slot, error) {
	return m.SlotsResp, nil
}

func (m *AppointmentServiceMock) BookAppointment(ctx context.Context, slotID, patientID uuid.UUID) (uuid.UUID, error) {
	return m.BookResp, m.BookErr
}

func (m *AppointmentServiceMock) ListMyAppointments(ctx context.Context, userID uuid.UUID, role string) ([]service.AppointmentBrief, error) {
	return m.ListResp, m.ListErr
}

func (m *AppointmentServiceMock) UpdateAppointmentStatus(ctx context.Context, appointmentID, therapistID uuid.UUID, status string) (service.AppointmentBrief, error) {
	return m.UpdateResp, m.UpdateErr
}
