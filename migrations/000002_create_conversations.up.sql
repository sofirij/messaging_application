CREATE TABLE conversations (
    id                      SERIAL PRIMARY KEY,
    type                    TEXT,
    name                    TEXT,
    avatar_url              TEXT,
    created_by              INT NOT NULL REFERENCES users(id),
    created_at              TIMESTAMPTZ DEFAULT NOW()
);
