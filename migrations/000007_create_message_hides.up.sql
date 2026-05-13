CREATE TABLE message_hides (
    message_id              INT NOT NULL REFERENCES messages(id),
    user_id                 INT NOT NULL REFERENCES users(id),
    hidden_at               TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (message_id, user_id)
);

CREATE INDEX idx_message_hides_user_id ON message_hides (user_id, message_id);
