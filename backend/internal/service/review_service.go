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
	rows, err := s.db.Queries.GetReviewsForTherapist(ctx, therapistID)
	if err != nil {
		return nil, err
	}
	var out []map[string]interface{}
	for _, r := range rows {
		firstName, lastName := "", ""
		if r.PatientName.Valid {
			displayName := r.PatientName.String
			if idx := strings.Index(displayName, " "); idx > 0 {
				firstName = displayName[:idx]
				lastName = strings.TrimSpace(displayName[idx+1:])
			} else {
				firstName = displayName
			}
		}
		patient := map[string]interface{}{
			"_id": r.PatientID.String(),
			"profile": map[string]interface{}{
				"firstName": firstName,
				"lastName":  lastName,
			},
		}
		m := map[string]interface{}{
			"_id":     r.ReviewID,
			"patient": patient,
			"rating":  r.Rating,
		}
		if r.Comment.Valid {
			m["comment"] = r.Comment.String
		}
		out = append(out, m)
	}
	return out, nil
}

type ForbiddenError struct{ Msg string }

func (e *ForbiddenError) Error() string { return e.Msg }
