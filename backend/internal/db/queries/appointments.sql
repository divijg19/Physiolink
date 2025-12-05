-- name: CreateAvailabilitySlots
-- params: therapist_id uuid, start_ts timestamptz, end_ts timestamptz
INSERT INTO availability_slots (therapist_id, start_ts, end_ts, status)
VALUES ($1, $2, $3, 'open')
ON CONFLICT DO NOTHING;

-- name: GetTherapistOpenSlots
-- params: therapist_id uuid
SELECT id, therapist_id, start_ts, end_ts, status
FROM availability_slots
WHERE therapist_id = $1 AND status = 'open'
ORDER BY start_ts ASC;

-- name: BookAppointmentTxLockSlot
-- params: slot_id uuid
SELECT id, therapist_id, start_ts, end_ts, status
FROM availability_slots
WHERE id = $1
FOR UPDATE;

-- name: InsertAppointment
-- params: slot_id uuid, patient_id uuid, therapist_id uuid, status text, notes text
INSERT INTO appointments (slot_id, patient_id, therapist_id, status, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateSlotStatus
-- params: slot_id uuid, status text
UPDATE availability_slots
SET status = $2
WHERE id = $1;

-- name: GetMySchedule
-- params: user_id uuid, role text
SELECT a.id, a.slot_id, a.patient_id, a.therapist_id, a.status, a.notes, a.created_at, a.updated_at
FROM appointments a
WHERE CASE WHEN $2 = 'pt' THEN a.therapist_id = $1 ELSE a.patient_id = $1 END
ORDER BY a.created_at ASC;
