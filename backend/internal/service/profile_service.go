package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/divijg19/physiolink/backend/internal/db"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
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

	var specialties []string
	if p.Specialty != "" {
		specialties = []string{p.Specialty}
	}

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

	// Build args for sqlc-generated query
	arg := db.CreateOrUpdateProfileParams{
		UserID:      userID,
		DisplayName: sql.NullString{String: displayName, Valid: displayName != ""},
		Bio:         sql.NullString{String: p.Bio, Valid: p.Bio != ""},
		Phone:       sql.NullString{Valid: false},
		Address:     pqtype.NullRawMessage{},
		Specialties: specialties,
		ProfileExtra: pqtype.NullRawMessage{
			RawMessage: extra,
			Valid:      len(extra) > 0,
		},
	}

	id, err := s.db.Queries.CreateOrUpdateProfile(ctx, arg)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	row, err := s.db.Queries.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.db.Queries.GetProfileWithUserInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	// derive first/last name
	displayName := ""
	if row.DisplayName.Valid {
		displayName = row.DisplayName.String
	}
	firstName := ""
	lastName := ""
	if displayName != "" {
		for i, ch := range displayName {
			if ch == ' ' {
				firstName = displayName[:i]
				if i+1 < len(displayName) {
					lastName = displayName[i+1:]
				}
				break
			}
		}
		if firstName == "" {
			firstName = displayName
		}
	}

	// parse extra
	extra := map[string]interface{}{}
	if row.ProfileExtra.Valid {
		_ = json.Unmarshal(row.ProfileExtra.RawMessage, &extra)
	}

	out := make(map[string]interface{})
	out["id"] = row.ID.String()
	out["user"] = map[string]interface{}{"email": userInfo.Email, "role": userInfo.Role}
	out["firstName"] = firstName
	out["lastName"] = lastName
	out["bio"] = ""
	if row.Bio.Valid {
		out["bio"] = row.Bio.String
	}
	if len(row.Specialties) > 0 {
		out["specialty"] = row.Specialties[0]
	} else {
		out["specialty"] = ""
	}
	out["rating"] = userInfo.Rating

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

// CreateEmptyProfile creates a minimal profile row for a newly registered user
// to mirror the Node behavior (firstName/lastName empty initially).
func (s *ProfileService) CreateEmptyProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	if err := s.db.Queries.CreateEmptyProfile(ctx, userID); err != nil {
		return nil, err
	}
	return s.GetProfile(ctx, userID)
}
