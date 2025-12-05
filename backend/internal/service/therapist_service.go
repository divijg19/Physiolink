package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

	// Base query: PT users joined with profiles
	// Note: location filtering depends on address JSON; we apply a simple ILIKE on display_name and specialties
	var args []interface{}
	where := []string{"u.role = 'pt'"}
	if p.Specialty != "" {
		args = append(args, "%"+p.Specialty+"%")
		where = append(where, "(p.specialties::text ILIKE $"+fmt.Sprint(len(args))+")")
	}
	if p.Location != "" {
		args = append(args, "%"+p.Location+"%")
		where = append(where, "(p.address::text ILIKE $"+fmt.Sprint(len(args))+")")
	}
	whereSQL := strings.Join(where, " AND ")

	countSQL := "SELECT COUNT(*) FROM users u LEFT JOIN profiles p ON p.user_id = u.id WHERE " + whereSQL
	var total int
	if err := s.db.Pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return TherapistListResult{}, err
	}

	q := `SELECT u.id, u.email, COALESCE(p.display_name,'') AS display_name,
				 COALESCE(p.specialties, ARRAY[]::text[]) AS specialties,
				 p.address, p.rating
          FROM users u
          LEFT JOIN profiles p ON p.user_id = u.id
          WHERE ` + whereSQL + `
          ORDER BY u.created_at DESC
          LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)

	args = append(args, limit, offset)
	rows, err := s.db.Pool.Query(ctx, q, args...)
	if err != nil {
		return TherapistListResult{}, err
	}
	defer rows.Close()

	var out []TherapistSummary
	var ids []string
	for rows.Next() {
		var id, email, displayName string
		var specialties []string
		var address []byte
		var rating sql.NullFloat64
		if err := rows.Scan(&id, &email, &displayName, &specialties, &address, &rating); err != nil {
			return TherapistListResult{}, err
		}
		// map display_name -> firstName/lastName heuristically; specialty -> first
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
		}
		if address != nil {
			prof["address"] = string(address)
		}
		if rating.Valid {
			prof["rating"] = rating.Float64
		}
		out = append(out, TherapistSummary{ID: id, Email: email, Profile: prof, AvailableSlots: 0, ReviewCount: 0})
		ids = append(ids, id)
	}

	// Aggregate available slot counts
	if len(ids) > 0 {
		// ANY($1) with text[]; cast to uuid in query
		// Optional date filter
		if p.Date != "" {
			qcnt := `SELECT therapist_id::text, COUNT(*)
					 FROM availability_slots
					 WHERE therapist_id = ANY($1::uuid[]) AND status = 'open' AND start_ts::date = $2::date
					 GROUP BY therapist_id`
			rows2, err := s.db.Pool.Query(ctx, qcnt, ids, p.Date)
			if err == nil {
				m := map[string]int{}
				for rows2.Next() {
					var tid string
					var c int
					_ = rows2.Scan(&tid, &c)
					m[tid] = c
				}
				rows2.Close()
				for i := range out {
					if c, ok := m[out[i].ID]; ok {
						out[i].AvailableSlots = c
					}
				}
			}
		} else {
			qcnt := `SELECT therapist_id::text, COUNT(*)
					 FROM availability_slots
					 WHERE therapist_id = ANY($1::uuid[]) AND status = 'open'
					 GROUP BY therapist_id`
			rows2, err := s.db.Pool.Query(ctx, qcnt, ids)
			if err == nil {
				m := map[string]int{}
				for rows2.Next() {
					var tid string
					var c int
					_ = rows2.Scan(&tid, &c)
					m[tid] = c
				}
				rows2.Close()
				for i := range out {
					if c, ok := m[out[i].ID]; ok {
						out[i].AvailableSlots = c
					}
				}
			}
		}

		// Aggregate review counts
		qrev := `SELECT a.therapist_id::text, COUNT(r.id)
				 FROM appointments a JOIN reviews r ON r.appointment_id = a.id
				 WHERE a.therapist_id = ANY($1::uuid[])
				 GROUP BY a.therapist_id`
		rows3, err := s.db.Pool.Query(ctx, qrev, ids)
		if err == nil {
			m := map[string]int{}
			for rows3.Next() {
				var tid string
				var c int
				_ = rows3.Scan(&tid, &c)
				m[tid] = c
			}
			rows3.Close()
			for i := range out {
				if c, ok := m[out[i].ID]; ok {
					out[i].ReviewCount = c
				}
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
