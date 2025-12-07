-- name: CreateAvailabilitySlots :exec
-- params: therapist_id uuid, start_ts timestamptz, end_ts timestamptz
INSERT INTO availability_slots (therapist_id, start_ts, end_ts, status)
VALUES ($1, $2, $3, 'open')
ON CONFLICT DO NOTHING;

-- name: GetTherapistOpenSlots :many
-- params: therapist_id uuid
SELECT id, therapist_id, start_ts, end_ts, status
FROM availability_slots
WHERE therapist_id = $1 AND status = 'open'
ORDER BY start_ts ASC;

-- name: BookAppointmentTxLockSlot :one
-- params: slot_id uuid
SELECT id, therapist_id, start_ts, end_ts, status
FROM availability_slots
WHERE id = $1
FOR UPDATE;

-- name: InsertAppointment :one
-- params: slot_id uuid, patient_id uuid, therapist_id uuid, status text, notes text
INSERT INTO appointments (slot_id, patient_id, therapist_id, status, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateSlotStatus :exec
-- params: slot_id uuid, status text
UPDATE availability_slots
SET status = $2
WHERE id = $1;

-- name: ListMyAppointmentsWithDetails :many
-- params: user_id uuid, role text
SELECT 
    a.id, 
    a.therapist_id, 
    a.patient_id, 
    a.status, 
    s.start_ts, 
    s.end_ts,
    p_pt.display_name as pt_display_name, 
    p_pt.profile_extra as pt_profile_extra,
    p_pa.display_name as pa_display_name, 
    p_pa.profile_extra as pa_profile_extra
FROM appointments a
JOIN availability_slots s ON s.id = a.slot_id
LEFT JOIN profiles p_pt ON p_pt.user_id = a.therapist_id
LEFT JOIN profiles p_pa ON p_pa.user_id = a.patient_id
WHERE CASE WHEN $2 = 'pt' THEN a.therapist_id = $1 ELSE a.patient_id = $1 END
ORDER BY s.start_ts ASC;

-- name: GetAppointmentTherapistID :one
-- params: appointment_id uuid
SELECT therapist_id
FROM appointments
WHERE id = $1;

-- name: UpdateAppointmentStatus :exec
-- params: appointment_id uuid, status text
UPDATE appointments
SET status = $2, updated_at = now()
WHERE id = $1;

-- name: GetAppointmentSlotStartTime :one
-- params: appointment_id uuid
SELECT a.slot_id, s.start_ts
FROM appointments a
JOIN availability_slots s ON s.id = a.slot_id
WHERE a.id = $1;

-- name: InsertReminder :exec
-- params: appointment_id uuid, scheduled_for timestamptz, payload jsonb
INSERT INTO reminders (appointment_id, scheduled_for, payload)
VALUES ($1, $2, $3);
