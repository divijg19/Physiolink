package service

import (
	"context"
	"encoding/json"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/google/uuid"
)

type ProfileService struct {
	db  *db.DB
	cfg *config.Config
}

func NewProfileService(d *db.DB, cfg *config.Config) *ProfileService {
	return &ProfileService{db: d, cfg: cfg}
}

// NodeProfileUpdate mirrors the payload expected from the existing
// React Native app (parity with the Node backend).
type NodeProfileUpdate struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Age             *int   `json:"age"`
	Gender          string `json:"gender"`
	Condition       string `json:"condition"`
	Goals           string `json:"goals"`
	Specialty       string `json:"specialty"`
	Bio             string `json:"bio"`
	Credentials     string `json:"credentials"`
	Location        string `json:"location"`
	ProfileImageURL string `json:"profileImageUrl"`
}

func (s *ProfileService) UpsertProfile(ctx context.Context, userID uuid.UUID, p NodeProfileUpdate) (uuid.UUID, error) {
	// Persist key fields across existing columns
	displayName := ""
	if p.FirstName != "" || p.LastName != "" {
		if p.FirstName != "" && p.LastName != "" {
			displayName = p.FirstName + " " + p.LastName
		} else if p.FirstName != "" {
			displayName = p.FirstName
		} else {
			displayName = p.LastName
		}
	}

	// Map specialty to first element of specialties array for compatibility
	var specialties []string
	if p.Specialty != "" {
		specialties = []string{p.Specialty}
	}

	// Store the remaining flexible fields in profile_extra JSONB
	extraMap := map[string]interface{}{}
	if p.Age != nil {
		extraMap["age"] = *p.Age
	}
	if p.Gender != "" {
		extraMap["gender"] = p.Gender
	}
	if p.Condition != "" {
		extraMap["condition"] = p.Condition
	}
	if p.Goals != "" {
		extraMap["goals"] = p.Goals
	}
	if p.Credentials != "" {
		extraMap["credentials"] = p.Credentials
	}
	if p.Location != "" {
		extraMap["location"] = p.Location
	}
	if p.ProfileImageURL != "" {
		extraMap["profileImageUrl"] = p.ProfileImageURL
	}
	extra, err := json.Marshal(extraMap)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	q := `INSERT INTO profiles (user_id, display_name, bio, specialties, profile_extra)
		  VALUES ($1, $2, $3, $4, $5)
		  ON CONFLICT (user_id) DO UPDATE SET display_name = EXCLUDED.display_name, bio = EXCLUDED.bio, specialties = EXCLUDED.specialties, profile_extra = EXCLUDED.profile_extra, updated_at = now()
		  RETURNING id`
	row := s.db.Pool.QueryRow(ctx, q, userID, displayName, p.Bio, pqStringArray(specialties), extra)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Join with users to include email and role like Node's populate
	q := `SELECT p.id, p.user_id, p.display_name, p.bio, p.specialties, p.profile_extra, COALESCE(p.rating, 0), u.email, u.role
		  FROM profiles p
		  JOIN users u ON u.id = p.user_id
		  WHERE p.user_id = $1`
	row := s.db.Pool.QueryRow(ctx, q, userID)
	var id uuid.UUID
	var uid uuid.UUID
	var displayName, bio string
	var specialties []string
	var profileExtra []byte
	var rating float64
	var email, role string
	if err := row.Scan(&id, &uid, &displayName, &bio, &specialties, &profileExtra, &rating, &email, &role); err != nil {
		return nil, err
	}

	// Derive firstName/lastName from display_name conservatively
	firstName := ""
	lastName := ""
	if displayName != "" {
		// simple split on first space
		for i, r := range displayName {
			if r == ' ' {
				firstName = displayName[:i]
				if i+1 < len(displayName) {
					lastName = displayName[i+1:]
				}
				break
			}
		}
		if firstName == "" { // no space
			firstName = displayName
		}
	}

	// Extract extra fields
	extra := map[string]interface{}{}
	if len(profileExtra) > 0 {
		_ = json.Unmarshal(profileExtra, &extra)
	}

	// Build Node-like response
	out := make(map[string]interface{})
	out["id"] = id.String()
	out["user"] = map[string]interface{}{"email": email, "role": role}
	out["firstName"] = firstName
	out["lastName"] = lastName
	out["bio"] = bio
	if len(specialties) > 0 {
		out["specialty"] = specialties[0]
	} else {
		out["specialty"] = ""
	}
	out["rating"] = rating

	// Copy select extras if present
	if v, ok := extra["age"]; ok {
		out["age"] = v
	}
	if v, ok := extra["gender"]; ok {
		out["gender"] = v
	}
	if v, ok := extra["condition"]; ok {
		out["condition"] = v
	}
	if v, ok := extra["goals"]; ok {
		out["goals"] = v
	}
	if v, ok := extra["credentials"]; ok {
		out["credentials"] = v
	}
	if v, ok := extra["location"]; ok {
		out["location"] = v
	}
	if v, ok := extra["profileImageUrl"]; ok {
		out["profileImageUrl"] = v
	}
	if v, ok := extra["isVerified"]; ok {
		out["isVerified"] = v
	} else {
		out["isVerified"] = false
	}

	return out, nil
}

// pqStringArray is a simple helper to convert []string to Postgres array input via text representation.
func pqStringArray(s []string) interface{} {
	if s == nil {
		return nil
	}
	return s
}

// CreateEmptyProfile creates a minimal profile row for a newly registered user
// to mirror the Node behavior (firstName/lastName empty initially).
func (s *ProfileService) CreateEmptyProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	// Ensure a row exists
	q := `INSERT INTO profiles (user_id, display_name, bio)
		  VALUES ($1, '', '')
		  ON CONFLICT (user_id) DO NOTHING`
	if _, err := s.db.Pool.Exec(ctx, q, userID); err != nil {
		return nil, err
	}
	// Return Node-like shape
	return s.GetProfile(ctx, userID)
}
