DROP INDEX idx_messages_conversation_id;
ALTER TABLE conversations DROP CONSTRAINT fk_conversations_last_message_id;
DROP TABLE messages;