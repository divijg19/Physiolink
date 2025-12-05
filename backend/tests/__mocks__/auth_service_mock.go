package __mocks__

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type AuthServiceMock struct {
	Users map[string]struct {
		ID   uuid.UUID
		Hash string
		Role string
	}
}

var ErrInvalidCredentials = errors.New("Invalid Credentials")
var ErrUserExists = errors.New("User already exists")

func NewAuthServiceMock() *AuthServiceMock {
	return &AuthServiceMock{Users: make(map[string]struct {
		ID   uuid.UUID
		Hash string
		Role string
	})}
}

func (m *AuthServiceMock) Register(ctx context.Context, email, password, role string) (uuid.UUID, string, error) {
	if _, ok := m.Users[email]; ok {
		return uuid.Nil, "", ErrUserExists
	}
	id := uuid.New()
	m.Users[email] = struct {
		ID   uuid.UUID
		Hash string
		Role string
	}{ID: id, Hash: password, Role: role}
	return id, role, nil
}

func (m *AuthServiceMock) Authenticate(ctx context.Context, email, password string) (uuid.UUID, string, error) {
	u, ok := m.Users[email]
	if !ok || u.Hash != password {
		return uuid.Nil, "", ErrInvalidCredentials
	}
	return u.ID, u.Role, nil
}
