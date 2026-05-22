CREATE TABLE conversations (
    id                      SERIAL PRIMARY KEY,
    type                    VARCHAR(10) NOT NULL,
    name                    TEXT,
    avatar_url              TEXT,
    created_by              INT NOT NULL REFERENCES users(id),
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    last_message_id         INT REFERENCES messages(id),
    last_message_read       INT
);

CREATE INDEX idx_conversations_last_message_at ON conversations (last_message_at DESC NULLS LAST);
