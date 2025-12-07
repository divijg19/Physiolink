-- name: GetUpcomingReminders :many
-- params: patient_id uuid
SELECT r.id, r.appointment_id, r.scheduled_for, r.payload,
       a.therapist_id, s.start_ts as appointment_start
FROM reminders r
JOIN appointments a ON a.id = r.appointment_id
JOIN availability_slots s ON s.id = a.slot_id
WHERE a.patient_id = $1 AND r.sent_at IS NULL AND r.scheduled_for <= now()
ORDER BY r.scheduled_for ASC;

-- name: GetUpcomingRemindersBefore :many
-- params: patient_id uuid, before timestamptz
SELECT r.id, r.appointment_id, r.scheduled_for, r.payload,
       a.therapist_id, s.start_ts as appointment_start
FROM reminders r
JOIN appointments a ON a.id = r.appointment_id
JOIN availability_slots s ON s.id = a.slot_id
WHERE a.patient_id = $1 AND r.sent_at IS NULL AND r.scheduled_for <= $2
ORDER BY r.scheduled_for ASC;
