CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  openid VARCHAR(128) NOT NULL UNIQUE,
  nickname VARCHAR(128) NOT NULL DEFAULT '',
  avatar_url TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tennis_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),

  date DATE NOT NULL,
  duration_minutes INT NOT NULL DEFAULT 120,
  rating SMALLINT NOT NULL DEFAULT 3,

  type SMALLINT NOT NULL DEFAULT 1,
  match_rank SMALLINT NOT NULL DEFAULT 0,

  court_name VARCHAR(128) NOT NULL DEFAULT '',
  cost NUMERIC(10,2) NOT NULL DEFAULT 0,
  racket_name VARCHAR(128) NOT NULL DEFAULT '',
  shoe_name VARCHAR(128) NOT NULL DEFAULT '',
  note TEXT NOT NULL DEFAULT '',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_tennis_sessions_user_date
  ON tennis_sessions(user_id, date DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tennis_sessions_user_created
  ON tennis_sessions(user_id, created_at DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tennis_sessions_user_deleted
  ON tennis_sessions(user_id, deleted_at);

ALTER TABLE tennis_sessions
  DROP CONSTRAINT IF EXISTS tennis_sessions_type_check,
  ADD CONSTRAINT tennis_sessions_type_check
  CHECK (type IN (1, 2, 3, 4, 5));

ALTER TABLE tennis_sessions
  DROP CONSTRAINT IF EXISTS tennis_sessions_match_rank_check,
  ADD CONSTRAINT tennis_sessions_match_rank_check
  CHECK (match_rank IN (0, 1, 2, 3, 4, 5));

ALTER TABLE tennis_sessions
  DROP CONSTRAINT IF EXISTS tennis_sessions_rating_check,
  ADD CONSTRAINT tennis_sessions_rating_check
  CHECK (rating >= 1 AND rating <= 5);
