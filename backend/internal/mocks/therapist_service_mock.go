package __mocks__

import (
	"context"

	"github.com/divijg19/physiolink/backend/internal/service"
)

type TherapistServiceMock struct {
	ListResp   service.TherapistListResult
	DetailResp map[string]interface{}
}

func (m *TherapistServiceMock) GetAllTherapists(ctx context.Context, params service.TherapistQueryParams) (service.TherapistListResult, error) {
	return m.ListResp, nil
}

func (m *TherapistServiceMock) GetTherapistByID(ctx context.Context, id, date string) (map[string]interface{}, error) {
	if m.DetailResp != nil {
		return m.DetailResp, nil
	}
	return map[string]interface{}{"_id": id}, nil
}

// MakeTherapistListResult is a helper to create a minimal list result.
func MakeTherapistListResult(ids []string) service.TherapistListResult {
	out := make([]service.TherapistSummary, 0, len(ids))
	for _, id := range ids {
		out = append(out, service.TherapistSummary{ID: id, Email: "", Profile: map[string]interface{}{}, AvailableSlots: 0, ReviewCount: 0})
	}
	return service.TherapistListResult{Data: out, Total: len(out), Page: 1, TotalPages: 1}
}
