package service

import (
	"context"
	"encoding/json"

	"github.com/divijg19/physiolink/backend/go/internal/config"
	"github.com/divijg19/physiolink/backend/go/internal/db"
	"github.com/google/uuid"
)

type ProfileService struct {
	db  *db.DB
	cfg *config.Config
}

func NewProfileService(d *db.DB, cfg *config.Config) *ProfileService {
	return &ProfileService{db: d, cfg: cfg}
}

type ProfileUpdate struct {
	DisplayName string                 `json:"display_name"`
	Bio         string                 `json:"bio"`
	Phone       string                 `json:"phone"`
	Address     map[string]interface{} `json:"address"`
	Specialties []string               `json:"specialties"`
	Extra       map[string]interface{} `json:"profile_extra"`
}

func (s *ProfileService) UpsertProfile(ctx context.Context, userID uuid.UUID, p ProfileUpdate) (uuid.UUID, error) {
	addr := []byte("null")
	extra := []byte("null")
	var err error
	if p.Address != nil {
		addr, err = json.Marshal(p.Address)
		if err != nil {
			return uuid.Nil, err
		}
	}
	if p.Extra != nil {
		extra, err = json.Marshal(p.Extra)
		if err != nil {
			return uuid.Nil, err
		}
	}
	var id uuid.UUID
	q := `INSERT INTO profiles (user_id, display_name, bio, phone, address, specialties, profile_extra)
          VALUES ($1, $2, $3, $4, $5, $6, $7)
          ON CONFLICT (user_id) DO UPDATE SET display_name = EXCLUDED.display_name, bio = EXCLUDED.bio, phone = EXCLUDED.phone, address = EXCLUDED.address, specialties = EXCLUDED.specialties, profile_extra = EXCLUDED.profile_extra, updated_at = now()
          RETURNING id`
	row := s.db.Pool.QueryRow(ctx, q, userID, p.DisplayName, p.Bio, p.Phone, addr, pqStringArray(p.Specialties), extra)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *ProfileService) GetProfile(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	q := `SELECT id, user_id, display_name, bio, phone, address, specialties, profile_extra, created_at, updated_at FROM profiles WHERE user_id = $1`
	row := s.db.Pool.QueryRow(ctx, q, userID)
	var id uuid.UUID
	var uid uuid.UUID
	var displayName, bio, phone string
	var address []byte
	var specialties []string
	var profileExtra []byte
	if err := row.Scan(&id, &uid, &displayName, &bio, &phone, &address, &specialties, &profileExtra, new(interface{}), new(interface{})); err != nil {
		return nil, err
	}
	out := make(map[string]interface{})
	out["id"] = id.String()
	out["user_id"] = uid.String()
	out["display_name"] = displayName
	out["bio"] = bio
	out["phone"] = phone
	if address != nil {
		var a interface{}
		_ = json.Unmarshal(address, &a)
		out["address"] = a
	}
	out["specialties"] = specialties
	if profileExtra != nil {
		var e interface{}
		_ = json.Unmarshal(profileExtra, &e)
		out["profile_extra"] = e
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
