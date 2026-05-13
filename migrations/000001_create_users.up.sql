CREATE TABLE users (
    id                  SERIAL PRIMARY KEY,
    username            VARCHAR(50) NOT NULL,
    password            TEXT NOT NULL,
    avatar_url          TEXT,
    last_seen_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_username_lower ON USERS (LOWER(username));