-- Create users table with new fields
CREATE TABLE IF NOT EXISTS "sessions" (
  "id" uuid PRIMARY KEY,
  "owner_id" bigint NOT NULL,
  "user_agent" TEXT NOT NULL,
  "refresh_token" TEXT NOT NULL,
  "client_ip" varchar(20),
  "is_blocked" BOOLEAN DEFAULT FALSE,
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);


ALTER TABLE "sessions" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");
