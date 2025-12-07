-- name: GetAppointmentForReview :one
-- params: appointment_id uuid, patient_id uuid
SELECT id
FROM appointments
WHERE id = $1 AND patient_id = $2 AND status = 'confirmed';

-- name: CreateReview :one
-- params: appointment_id uuid, patient_id uuid, rating int, comment text
INSERT INTO reviews (appointment_id, patient_id, rating, comment)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetTherapistAverageRating :one
-- params: therapist_id uuid
SELECT AVG(r.rating)::float as avg_rating
FROM reviews r
JOIN appointments a ON a.id = r.appointment_id
WHERE a.therapist_id = $1;

-- name: UpdateProfileRating :exec
-- params: user_id uuid, rating float
UPDATE profiles
SET rating = $2
WHERE user_id = $1;

-- name: GetReviewsForTherapist :many
-- params: therapist_id uuid
SELECT r.id::text as review_id, r.patient_id, r.rating, r.comment,
       p.display_name as patient_name, r.created_at
FROM reviews r
JOIN appointments a ON a.id = r.appointment_id
LEFT JOIN profiles p ON p.user_id = r.patient_id
WHERE a.therapist_id = $1
ORDER BY r.created_at DESC;
