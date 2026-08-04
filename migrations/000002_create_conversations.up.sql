CREATE TABLE conversations (
    id                      SERIAL PRIMARY KEY,
    type                    VARCHAR(10) NOT NULL,
    name                    TEXT,
    avatar_url              TEXT,
    created_by              INT NOT NULL REFERENCES users(id),
    created_at              TIMESTAMPTZ DEFAULT NOW()
);
