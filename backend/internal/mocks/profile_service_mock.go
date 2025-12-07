package __mocks__

import (
	"context"

	"github.com/divijg19/physiolink/backend/internal/service"
	"github.com/google/uuid"
)

type ProfileServiceMock struct{}

func NewProfileServiceMock() *ProfileServiceMock { return &ProfileServiceMock{} }
func (m *ProfileServiceMock) UpsertProfile(ctx context.Context, userID uuid.UUID, p service.NodeProfileUpdate) (uuid.UUID, error) {
	return userID, nil
}

func (m *ProfileServiceMock) GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":         userID.String(),
		"firstName":  "",
		"lastName":   "",
		"bio":        "",
		"specialty":  "",
		"rating":     0,
		"user":       map[string]interface{}{"email": "", "role": ""},
		"isVerified": false,
	}, nil
}

func (m *ProfileServiceMock) CreateEmptyProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":         userID.String(),
		"firstName":  "",
		"lastName":   "",
		"bio":        "",
		"specialty":  "",
		"rating":     0,
		"user":       map[string]interface{}{"email": "", "role": ""},
		"isVerified": false,
	}, nil
}
