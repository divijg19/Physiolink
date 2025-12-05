-- Add rating to profiles for aggregated review scores
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS rating NUMERIC;