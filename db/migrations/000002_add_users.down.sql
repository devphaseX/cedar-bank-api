-- Drop indexes created in the up migration
DROP INDEX IF EXISTS idx_accounts_currency;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Drop foreign key constraint added in the up migration
ALTER TABLE "accounts"
  DROP CONSTRAINT IF EXISTS fk_accounts_users;

-- Revert the "currency" column to its original type if necessary
-- Note: You need to replace 'original_type' with the actual original data type of the currency column
-- ALTER TABLE "accounts" ALTER COLUMN "currency" TYPE original_type;

-- Remove the "owner_id" column added in the up migration
ALTER TABLE "accounts"
  DROP COLUMN IF EXISTS "owner_id";

-- Restore the "owner" column if it was originally dropped
-- Note: You need to replace 'original_type' with the actual data type of the owner column
-- ALTER TABLE "accounts" ADD COLUMN "owner" original_type;

-- Drop the "users" table created in the up migration
DROP TABLE IF EXISTS "users";

ALTER TABLE "accounts"
DROP CONSTRAINT unique_owner_currency;
