package service

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"

	"github.com/divijg19/physiolink/backend/internal/db"
)

type ReviewService struct{ db *db.DB }

func NewReviewService(d *db.DB) *ReviewService { return &ReviewService{db: d} }

// CreateReview enforces that the patient has at least one appointment with the therapist
// before allowing a review. It links the review to the most recent appointment between them.
func (s *ReviewService) CreateReview(ctx context.Context, patientID, therapistID uuid.UUID, rating int, comment string) (map[string]interface{}, error) {
	// find recent appointment
	var apptID uuid.UUID
	err := s.db.Pool.QueryRow(ctx, `SELECT id FROM appointments
        WHERE patient_id=$1 AND therapist_id=$2
        ORDER BY created_at DESC LIMIT 1`, patientID, therapistID).Scan(&apptID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ForbiddenError{Msg: "Only patients with an appointment can leave a review"}
		}
		return nil, err
	}

	var reviewID uuid.UUID
	err = s.db.Pool.QueryRow(ctx, `INSERT INTO reviews (appointment_id, patient_id, rating, comment)
        VALUES ($1,$2,$3,$4) RETURNING id`, apptID, patientID, rating, comment).Scan(&reviewID)
	if err != nil {
		return nil, err
	}

	// Update average rating on therapist profile
	var avg sql.NullFloat64
	_ = s.db.Pool.QueryRow(ctx, `SELECT AVG(r.rating)::float
		FROM appointments a JOIN reviews r ON r.appointment_id = a.id
		WHERE a.therapist_id = $1`, therapistID).Scan(&avg)
	if avg.Valid {
		_, _ = s.db.Pool.Exec(ctx, `UPDATE profiles SET rating = $2 WHERE user_id = $1`, therapistID, avg.Float64)
	}

	out := map[string]interface{}{
		"_id":       reviewID.String(),
		"therapist": map[string]interface{}{"_id": therapistID.String()},
		"patient":   map[string]interface{}{"_id": patientID.String()},
		"rating":    rating,
		"comment":   comment,
	}
	return out, nil
}

func (s *ReviewService) GetReviewsForTherapist(ctx context.Context, therapistID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := s.db.Pool.Query(ctx, `SELECT r.id::text, r.patient_id, r.rating, r.comment,
		COALESCE(p.display_name,''), p.profile_extra
		FROM appointments a JOIN reviews r ON r.appointment_id = a.id
		LEFT JOIN profiles p ON p.user_id = r.patient_id
		WHERE a.therapist_id=$1 ORDER BY r.created_at DESC`, therapistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []map[string]interface{}
	for rows.Next() {
		var rid string
		var pid uuid.UUID
		var rating int
		var comment sql.NullString
		var displayName string
		var profileExtra []byte
		if err := rows.Scan(&rid, &pid, &rating, &comment, &displayName, &profileExtra); err != nil {
			return nil, err
		}
		firstName, lastName := displayName, ""
		if idx := strings.Index(displayName, " "); idx > 0 {
			firstName = displayName[:idx]
			lastName = strings.TrimSpace(displayName[idx+1:])
		}
		patient := map[string]interface{}{
			"_id": pid.String(),
			"profile": map[string]interface{}{
				"firstName": firstName,
				"lastName":  lastName,
			},
		}
		m := map[string]interface{}{
			"_id":     rid,
			"patient": patient,
			"rating":  rating,
		}
		if comment.Valid {
			m["comment"] = comment.String
		}
		out = append(out, m)
	}
	return out, nil
}

type ForbiddenError struct{ Msg string }

func (e *ForbiddenError) Error() string { return e.Msg }
