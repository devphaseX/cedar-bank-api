-- Create users table with new fields
CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(50) UNIQUE NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "fullname" varchar(255) NOT NULL,
  "hashed_password" text NOT NULL,
  "password_changed_at" timestamptz,
  "created_at" timestamptz DEFAULT (now())
);


-- Change column type for "currency"
ALTER TABLE "accounts"
  ALTER COLUMN "currency" TYPE varchar;

-- Create new indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON "users" ("username");
CREATE INDEX IF NOT EXISTS idx_users_email ON "users" ("email");
CREATE INDEX IF NOT EXISTS idx_accounts_currency ON "accounts" ("id", "currency");

ALTER TABLE "accounts"
ADD CONSTRAINT unique_owner_currency UNIQUE ("owner_id", "currency");
