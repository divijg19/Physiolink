package testutil

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/divijg19/physiolink/backend/internal/service"
)

// CreateUserAndToken registers a user via the AuthService and returns the user ID and a signed JWT token.
func CreateUserAndToken(ctx context.Context, database *db.DB, cfg *config.Config, email, password, role string) (uuid.UUID, string, error) {
	authSvc := service.NewAuthService(database, cfg)
	id, _, err := authSvc.Register(ctx, email, password, role)
	if err != nil {
		return uuid.Nil, "", err
	}
	// create token matching handlers' format
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"id":   id.String(),
			"role": role,
		},
		"exp": time.Now().Add(5 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return uuid.Nil, "", err
	}
	return id, signed, nil
}
