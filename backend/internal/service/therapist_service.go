package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/db"
)

type TherapistService struct {
	db *db.DB
}

func NewTherapistService(d *db.DB) *TherapistService { return &TherapistService{db: d} }

type TherapistQueryParams struct {
	Specialty string
	Location  string
	Page      int
	Limit     int
	Sort      string
	Date      string
	Available bool
}

type TherapistSummary struct {
	ID             string                 `json:"_id"`
	Email          string                 `json:"email"`
	Profile        map[string]interface{} `json:"profile"`
	AvailableSlots int                    `json:"availableSlotsCount"`
	ReviewCount    int                    `json:"reviewCount"`
}

type TherapistListResult struct {
	Data       []TherapistSummary `json:"data"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	TotalPages int                `json:"totalPages"`
}

// GetAllTherapists returns paginated list of PT users with optional filters.
func (s *TherapistService) GetAllTherapists(ctx context.Context, p TherapistQueryParams) (TherapistListResult, error) {
	page := p.Page
	if page < 1 {
		page = 1
	}
	limit := p.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Use sqlc queries for therapist list and counts
	specParam := "%" + p.Specialty + "%"
	locParam := "%" + p.Location + "%"
	total64, err := s.db.Queries.GetTherapistCount(ctx, db.GetTherapistCountParams{Column1: specParam, Column2: locParam})
	if err != nil {
		return TherapistListResult{}, err
	}
	total := int(total64)

	therapists, err := s.db.Queries.GetTherapists(ctx, db.GetTherapistsParams{Column1: specParam, Column2: locParam, Limit: int32(limit), Offset: int32(offset)})
	if err != nil {
		return TherapistListResult{}, err
	}

	var out []TherapistSummary
	var ids []uuid.UUID
	for _, t := range therapists {
		displayName := t.DisplayName
		firstName, lastName := displayName, ""
		if idx := strings.Index(displayName, " "); idx > 0 {
			firstName = displayName[:idx]
			lastName = strings.TrimSpace(displayName[idx+1:])
		}
		specialty := ""
		if len(t.Specialties) > 0 {
			specialty = t.Specialties[0]
		}
		prof := map[string]interface{}{
			"firstName": firstName,
			"lastName":  lastName,
			"specialty": specialty,
		}
		if t.Address.Valid {
			prof["address"] = string(t.Address.RawMessage)
		}
		if t.Rating.Valid {
			// rating stored as string in generated type for compatibility; attempt parse
			prof["rating"] = t.Rating.String
		}
		out = append(out, TherapistSummary{ID: t.ID.String(), Email: t.Email, Profile: prof, AvailableSlots: 0, ReviewCount: 0})
		ids = append(ids, t.ID)
	}

	// Aggregate available slot counts
	if len(ids) > 0 {
		// convert []uuid.UUID to []uuid.UUID param for sqlc
		availRows, _ := s.db.Queries.GetAvailabilityCounts(ctx, db.GetAvailabilityCountsParams{Column1: ids, Column2: p.Date})
		mAvail := map[string]int64{}
		for _, r := range availRows {
			mAvail[r.TherapistID] = r.Count
		}
		revRows, _ := s.db.Queries.GetReviewCounts(ctx, ids)
		mRev := map[string]int64{}
		for _, r := range revRows {
			// column name from generator is ATherapistID
			mRev[r.ATherapistID] = r.Count
		}
		for i := range out {
			if c, ok := mAvail[out[i].ID]; ok {
				out[i].AvailableSlots = int(c)
			}
			if c, ok := mRev[out[i].ID]; ok {
				out[i].ReviewCount = int(c)
			}
		}
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}
	return TherapistListResult{Data: out, Total: total, Page: page, TotalPages: totalPages}, nil
}

// GetTherapistByID returns a single PT with basic profile and placeholder slot/review counts.
func (s *TherapistService) GetTherapistByID(ctx context.Context, id string, date string) (map[string]interface{}, error) {
	q := `SELECT u.id, u.email, COALESCE(p.display_name,''), COALESCE(p.specialties, ARRAY[]::text[]), p.address, COALESCE(p.bio,''), p.rating
          FROM users u LEFT JOIN profiles p ON p.user_id = u.id
          WHERE u.id = $1 AND u.role = 'pt'`
	var email, displayName, bio string
	var uid string
	var specialties []string
	var address []byte
	var rating sql.NullFloat64
	if err := s.db.Pool.QueryRow(ctx, q, id).Scan(&uid, &email, &displayName, &specialties, &address, &bio, &rating); err != nil {
		return nil, fmt.Errorf("therapist not found")
	}
	firstName, lastName := displayName, ""
	if idx := strings.Index(displayName, " "); idx > 0 {
		firstName = displayName[:idx]
		lastName = strings.TrimSpace(displayName[idx+1:])
	}
	specialty := ""
	if len(specialties) > 0 {
		specialty = specialties[0]
	}
	prof := map[string]interface{}{
		"firstName": firstName,
		"lastName":  lastName,
		"specialty": specialty,
		"bio":       bio,
	}
	if address != nil {
		prof["address"] = string(address)
	}
	if rating.Valid {
		prof["rating"] = rating.Float64
	}

	// available slots
	var slotsQuery string
	var rowsArgs []interface{}
	if date != "" {
		slotsQuery = `SELECT id::text, start_ts, end_ts FROM availability_slots
					  WHERE therapist_id = $1::uuid AND status = 'open' AND start_ts::date = $2::date
					  ORDER BY start_ts ASC LIMIT 20`
		rowsArgs = []interface{}{id, date}
	} else {
		slotsQuery = `SELECT id::text, start_ts, end_ts FROM availability_slots
					  WHERE therapist_id = $1::uuid AND status = 'open'
					  ORDER BY start_ts ASC LIMIT 20`
		rowsArgs = []interface{}{id}
	}
	rows, err := s.db.Pool.Query(ctx, slotsQuery, rowsArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	slots := make([]map[string]interface{}, 0, 20)
	for rows.Next() {
		var sid string
		var startTs, endTs interface{}
		if err := rows.Scan(&sid, &startTs, &endTs); err != nil {
			return nil, err
		}
		slots = append(slots, map[string]interface{}{
			"_id":       sid,
			"startTime": startTs,
			"endTime":   endTs,
		})
	}
	// review count
	qrev := `SELECT COUNT(r.id)
			 FROM appointments a JOIN reviews r ON r.appointment_id = a.id
			 WHERE a.therapist_id = $1::uuid`
	var reviewCount int
	_ = s.db.Pool.QueryRow(ctx, qrev, id).Scan(&reviewCount)

	out := map[string]interface{}{
		"_id":            uid,
		"email":          email,
		"profile":        prof,
		"availableSlots": slots,
		"reviewCount":    reviewCount,
	}
	return out, nil
}
