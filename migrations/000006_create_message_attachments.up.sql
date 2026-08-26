CREATE TABLE message_attachments (
    id                  SERIAL PRIMARY KEY,
    message_id          INT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    type                TEXT,
    url                 TEXT NOT NULL,
    filename            TEXT NOT NULL,
    size                BIGINT NOT NULL
);

CREATE INDEX idx_message_attachments_message_id ON message_attachments (message_id);