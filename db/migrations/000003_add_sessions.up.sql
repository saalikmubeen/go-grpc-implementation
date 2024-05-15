CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL, -- References users.username
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL, -- User agent of the client (i.e browser, mobile app, etc.)
  "client_ip" varchar NOT NULL, -- IP address of the client
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL, -- Expiry date of the session or refresh token
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");