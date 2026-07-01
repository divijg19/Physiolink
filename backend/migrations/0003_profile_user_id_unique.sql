-- Add unique index on profiles.user_id for ON CONFLICT support in CreateEmptyProfile and CreateOrUpdateProfile
CREATE UNIQUE INDEX IF NOT EXISTS ux_profiles_user_id ON profiles(user_id);
