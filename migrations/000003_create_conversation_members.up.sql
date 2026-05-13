CREATE TABLE conversation_members (
    conversation_id             INT NOT NULL REFERENCES conversations(id),
    user_id                     INT NOT NULL REFERENCES users(id),
    role                        VARCHAR(10) NOT NULL DEFAULT 'member',
    joined_at                   TIMESTAMPTZ DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ,
    PRIMARY KEY (conversation_id, user_id)         
);

CREATE INDEX idx_conversation_members_user_id ON conversation_members (user_id);