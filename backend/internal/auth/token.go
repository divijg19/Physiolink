// Package auth issues and verifies the session tokens used by both the JSON
// API and the cookie-authenticated web portal.
//
// Token construction was previously duplicated in three places
// (handlers/auth.go twice and testutil/users.go), which let the web portal
// write the user's *role* into the auth_token cookie instead of a signed JWT
// and silently break every cookie-authenticated route. Keeping issue and
// verify here means the format has exactly one definition.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned when a token fails to parse or its signature
// does not verify.
var ErrInvalidToken = errors.New("invalid token")

// TTL is how long an issued token remains valid. The auth cookie must not
// outlive it, or the portal breaks once the token expires while the cookie
// still claims a live session.
const TTL = 5 * time.Hour

// Claims is the token payload. The user id and role are nested under "user"
// because that is the shape the middleware and the Flutter client read.
type Claims struct {
	User struct {
		ID   string `json:"id"`
		Role string `json:"role"`
	} `json:"user"`
	ExpiresAt int64 `json:"exp"`
}

// jwt.Claims requires these accessors. Only the expiry is meaningful for this
// token; the rest are absent claims and report zero values.
func (c Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	if c.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

func (c Claims) GetIssuedAt() (*jwt.NumericDate, error)  { return nil, nil }
func (c Claims) GetNotBefore() (*jwt.NumericDate, error) { return nil, nil }
func (c Claims) GetIssuer() (string, error)              { return "", nil }
func (c Claims) GetSubject() (string, error)             { return "", nil }
func (c Claims) GetAudience() (jwt.ClaimStrings, error)  { return nil, nil }

// Issuer signs and verifies tokens with a shared secret.
type Issuer struct {
	secret []byte
	ttl    time.Duration
}

// NewIssuer returns an Issuer using the given secret and TTL.
func NewIssuer(secret string, ttl time.Duration) *Issuer {
	if ttl <= 0 {
		ttl = TTL
	}
	return &Issuer{secret: []byte(secret), ttl: ttl}
}

// Issue signs a token for the given user and role.
func (i *Issuer) Issue(id uuid.UUID, role string) (string, error) {
	claims := Claims{ExpiresAt: time.Now().Add(i.ttl).Unix()}
	claims.User.ID = id.String()
	claims.User.Role = role

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Parse verifies a token's signature and expiry and returns its claims.
func (i *Issuer) Parse(raw string) (*Claims, error) {
	parsed, err := jwt.NewParser(jwt.WithValidMethods([]string{"HS256"})).
		ParseWithClaims(raw, &Claims{}, func(*jwt.Token) (interface{}, error) {
			return i.secret, nil
		})
	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || claims.User.ID == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
