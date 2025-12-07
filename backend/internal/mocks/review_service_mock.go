package __mocks__

import (
	"context"

	"github.com/google/uuid"
)

type ReviewServiceMock struct {
	CreateResp map[string]interface{}
	CreateErr  error
	ListResp   []map[string]interface{}
	ListErr    error
}

func (m *ReviewServiceMock) CreateReview(ctx context.Context, patientID, therapistID uuid.UUID, rating int, comment string) (map[string]interface{}, error) {
	return m.CreateResp, m.CreateErr
}

func (m *ReviewServiceMock) GetReviewsForTherapist(ctx context.Context, therapistID uuid.UUID) ([]map[string]interface{}, error) {
	return m.ListResp, m.ListErr
}
