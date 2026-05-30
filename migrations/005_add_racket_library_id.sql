ALTER TABLE racket
  ADD COLUMN library_id BIGINT NOT NULL DEFAULT 0 AFTER user_id,
  ADD INDEX idx_racket_user_library (user_id, library_id);
