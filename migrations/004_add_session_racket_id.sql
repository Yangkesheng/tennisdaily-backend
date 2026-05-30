ALTER TABLE tennis_sessions
  ADD COLUMN racket_id BIGINT NOT NULL DEFAULT 0 AFTER cost,
  ADD INDEX idx_tennis_sessions_user_racket_deleted (user_id, racket_id, deleted_at);
