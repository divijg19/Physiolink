package testutil

import (
	"context"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/auth"
	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// CreateUserAndToken registers a user via the AuthService and returns the user ID and a signed JWT token.
func CreateUserAndToken(ctx context.Context, database *db.DB, cfg *config.Config, email, password, role string) (uuid.UUID, string, error) {
	authSvc := service.NewAuthService(database, cfg)
	id, actualRole, err := authSvc.Register(ctx, email, password, role)
	if err != nil {
		return uuid.Nil, "", err
	}
	// Mint the token through the same issuer the handlers use, so a test can
	// never pass against a token shape the server would reject.
	signed, err := auth.NewIssuer(cfg.JWTSecret, auth.TTL).Issue(id, actualRole)
	if err != nil {
		return uuid.Nil, "", err
	}
	return id, signed, nil
}
