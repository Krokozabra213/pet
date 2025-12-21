CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_image_message_id
ON images(message_id);