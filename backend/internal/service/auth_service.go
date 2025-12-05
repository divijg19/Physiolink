package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	// hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, "", err
	}

	// Use raw query execution; when using sqlc replace with generated call
	var id uuid.UUID
	q := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	row := s.db.Pool.QueryRow(ctx, q, email, string(hash), role)
	if err := row.Scan(&id); err != nil {
		// Map unique violation to ErrUserExists (constraint name may vary)
		// Fallback to generic error otherwise
		return uuid.Nil, "", ErrUserExists
	}
	_ = time.Now()
	return id, role, nil
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (uuid.UUID, string, error) {
	var id uuid.UUID
	var pwHash string
	var role string
	q := `SELECT id, password_hash, role FROM users WHERE email = $1`
	row := s.db.Pool.QueryRow(ctx, q, email)
	if err := row.Scan(&id, &pwHash, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, "", ErrInvalidCredentials
		}
		return uuid.Nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(password)); err != nil {
		return uuid.Nil, "", ErrInvalidCredentials
	}
	return id, role, nil
}
