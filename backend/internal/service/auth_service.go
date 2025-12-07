package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
)

type AuthService struct {
	db  *db.DB
	cfg *config.Config
}

func NewAuthService(d *db.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: d, cfg: cfg}
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *AuthService) Register(ctx context.Context, email, password, role string) (uuid.UUID, string, error) {
	if email == "" || password == "" {
		return uuid.Nil, "", errors.New("email and password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, "", err
	}

	arg := db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	id, err := s.db.Queries.CreateUser(ctx, arg)
	if err != nil {
		// Map unique violation to ErrUserExists when constraint violated
		return uuid.Nil, "", ErrUserExists
	}
	return id, role, nil
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (uuid.UUID, string, error) {
	user, err := s.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", ErrInvalidCredentials
		}
		return uuid.Nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, "", ErrInvalidCredentials
	}
	return user.ID, user.Role, nil
}
