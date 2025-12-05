-- name: CreateUser
-- params: email text, password_hash text, role text
-- result: id uuid
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetUserByEmail
-- params: email text
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID
-- params: id uuid
SELECT id, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateOrUpdateProfile
-- params: user_id uuid, display_name text, bio text, phone text, address jsonb, specialties text[], profile_extra jsonb
-- result: id uuid
INSERT INTO profiles (user_id, display_name, bio, phone, address, specialties, profile_extra)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) DO UPDATE SET
  display_name = EXCLUDED.display_name,
  bio = EXCLUDED.bio,
  phone = EXCLUDED.phone,
  address = EXCLUDED.address,
  specialties = EXCLUDED.specialties,
  profile_extra = EXCLUDED.profile_extra,
  updated_at = now()
RETURNING id;

-- name: GetProfileByUserID
-- params: user_id uuid
SELECT id, user_id, display_name, bio, phone, address, specialties, profile_extra, created_at, updated_at
FROM profiles
WHERE user_id = $1;
