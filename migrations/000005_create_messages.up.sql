CREATE TABLE messages (
    id                  SERIAL PRIMARY KEY,
    conversation_id     INT NOT NULL REFERENCES conversations(id),
    sender_id           INT NOT NULL REFERENCES users(id),
    reply_to_id         INT REFERENCES messages(id),
    body                TEXT,
    deleted_at          TIMESTAMPTZ,
    edited_at           TIMESTAMPTZ,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation_id ON messages (conversation_id, id DESC);

ALTER TABLE conversations
ADD CONSTRAINT fk_conversations_last_message_id
FOREIGN KEY (last_message_id) REFERENCES messages(id);