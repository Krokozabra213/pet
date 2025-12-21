CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_text_message_id
ON texts(message_id);