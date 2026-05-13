CREATE TABLE message_attachments (
    id                  SERIAL PRIMARY KEY,
    message_id          INT NOT NULL REFERENCES messages(id),
    type                VARCHAR(10) NOT NULL,
    url                 TEXT NOT NULL,
    filename            TEXT,
    size                BIGINT
);

CREATE INDEX idx_message_attachments_message_id ON message_attachments (message_id);