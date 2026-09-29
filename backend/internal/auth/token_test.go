package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const secret = "test-secret-at-least-32-characters-long"

func TestIssueParseRoundTrip(t *testing.T) {
	issuer := NewIssuer(secret, TTL)
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	raw, err := issuer.Issue(id, "pt")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if raw == "" {
		t.Fatal("Issue returned an empty token")
	}
	// A token must never be a bare role string: the portal used to write
	// "patient" straight into the cookie, which no parser could accept.
	if raw == "pt" || raw == "patient" {
		t.Fatalf("Issue returned the bare role %q instead of a JWT", raw)
	}

	claims, err := issuer.Parse(raw)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.User.ID != id.String() {
		t.Errorf("claims.User.ID = %q, want %q", claims.User.ID, id)
	}
	if claims.User.Role != "pt" {
		t.Errorf("claims.User.Role = %q, want %q", claims.User.Role, "pt")
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	issuer := NewIssuer(secret, TTL)
	raw, err := issuer.Issue(uuid.New(), "patient")
	if err != nil {
		t.Fatal(err)
	}

	other := NewIssuer("a-completely-different-secret-value!!", TTL)
	if _, err := other.Parse(raw); err == nil {
		t.Fatal("Parse accepted a token signed with a different secret")
	}
}

func TestParseRejectsExpired(t *testing.T) {
	// Signed directly with a past exp rather than via Issue, so the assertion
	// is deterministic. NewIssuer deliberately clamps a non-positive TTL to
	// the default, so a negative-TTL issuer cannot produce an expired token.
	issuer := NewIssuer(secret, TTL)
	expired := Claims{ExpiresAt: time.Now().Add(-time.Hour).Unix()}
	expired.User.ID = uuid.NewString()
	expired.User.Role = "patient"

	raw, err := jwt.NewWithClaims(jwt.SigningMethodHS256, expired).
		SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := issuer.Parse(raw); err == nil {
		t.Fatal("Parse accepted an expired token")
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	issuer := NewIssuer(secret, TTL)
	for _, bad := range []string{"", "patient", "therapist", "a.b.c", "not-a-jwt"} {
		if _, err := issuer.Parse(bad); err == nil {
			t.Errorf("Parse accepted %q", bad)
		}
	}
}

func TestParseRejectsNoneAlgorithm(t *testing.T) {
	// The classic JWT bypass: sign with alg=none and the empty secret.
	issuer := NewIssuer(secret, TTL)
	const unsigned = "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0." +
		"eyJ1c2VyIjp7ImlkIjoiMSIsInJvbGUiOiJwdCJ9LCJleHAiOjk5OTk5OTk5OTl9."
	if _, err := issuer.Parse(unsigned); err == nil {
		t.Fatal("Parse accepted an alg=none token")
	}
}

func TestDefaultTTLApplied(t *testing.T) {
	if got := NewIssuer(secret, 0).ttl; got != TTL {
		t.Errorf("ttl = %v, want %v", got, TTL)
	}
}
