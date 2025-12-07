-- name: GetTherapistCount :one
-- params: specialty text, location text
SELECT COUNT(*)
FROM users u
LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.role = 'pt'
  AND ($1 = '' OR p.specialties::text ILIKE $1)
  AND ($2 = '' OR p.address::text ILIKE $2);

-- name: GetTherapists :many
-- params: specialty text, location text, limit int, offset int
SELECT u.id, u.email, COALESCE(p.display_name,'') AS display_name,
       COALESCE(p.specialties, ARRAY[]::text[]) AS specialties,
       p.address, p.rating
FROM users u
LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.role = 'pt'
  AND ($1 = '' OR p.specialties::text ILIKE $1)
  AND ($2 = '' OR p.address::text ILIKE $2)
ORDER BY u.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetAvailabilityCounts :many
-- params: therapist_ids text[], date text
SELECT therapist_id::text, COUNT(*)
FROM availability_slots
WHERE therapist_id = ANY($1::uuid[])
  AND status = 'open'
  AND ($2 = '' OR start_ts::date = $2::date)
GROUP BY therapist_id;

-- name: GetReviewCounts :many
-- params: therapist_ids text[]
SELECT a.therapist_id::text, COUNT(r.id)
FROM appointments a
JOIN reviews r ON r.appointment_id = a.id
WHERE a.therapist_id = ANY($1::uuid[])
GROUP BY a.therapist_id;

-- name: GetTherapistByID :one
-- params: id uuid
SELECT u.id, u.email, COALESCE(p.display_name,''), COALESCE(p.specialties, ARRAY[]::text[]), p.address, COALESCE(p.bio,''), p.rating
FROM users u LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.id = $1 AND u.role = 'pt';

-- name: GetTherapistAvailabilitySlots :many
-- params: therapist_id uuid, date text, limit int
SELECT id::text, start_ts, end_ts
FROM availability_slots
WHERE therapist_id = $1 AND status = 'open' AND ($2 = '' OR start_ts::date = $2::date)
ORDER BY start_ts ASC
LIMIT $3;

-- name: GetTherapistReviewCount :one
-- params: therapist_id uuid
SELECT COUNT(r.id)
FROM appointments a JOIN reviews r ON r.appointment_id = a.id
WHERE a.therapist_id = $1;
