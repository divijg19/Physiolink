package __mocks__

import (
	"context"

	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/google/uuid"
)

type ReminderServiceMock struct {
	ListResp []service.ReminderItem
	ListErr  error
}

func (m *ReminderServiceMock) ListForPatient(ctx context.Context, patientID uuid.UUID) ([]service.ReminderItem, error) {
	return m.ListResp, m.ListErr
}
