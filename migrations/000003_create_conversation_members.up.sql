CREATE TABLE conversation_members (
    conversation_id             INT NOT NULL REFERENCES conversations(id),
    user_id                     INT NOT NULL REFERENCES users(id),
    joined_at                   TIMESTAMPTZ DEFAULT NOW(),
    after_cursor                INT,
    deleted_at                  TIMESTAMPTZ,
    last_message_read           INT,
    PRIMARY KEY (conversation_id, user_id)         
);

CREATE INDEX idx_conversation_members_user_id ON conversation_members (user_id);