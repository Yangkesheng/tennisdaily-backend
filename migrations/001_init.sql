CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  openid VARCHAR(128) NOT NULL UNIQUE,
  nickname VARCHAR(128) NOT NULL DEFAULT '',
  avatar_url TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='用户表';

CREATE TABLE IF NOT EXISTS tennis_sessions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,

  date DATE NOT NULL,
  duration_minutes INT NOT NULL DEFAULT 120,
  rating SMALLINT NOT NULL DEFAULT 3,

  type SMALLINT NOT NULL DEFAULT 1,
  match_rank SMALLINT NOT NULL DEFAULT 0,

  court_name VARCHAR(128) NOT NULL DEFAULT '',
  cost DECIMAL(10,2) NOT NULL DEFAULT 0,
  racket_name VARCHAR(128) NOT NULL DEFAULT '',
  shoe_name VARCHAR(128) NOT NULL DEFAULT '',
  note TEXT NOT NULL,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,

  CONSTRAINT tennis_sessions_type_check
    CHECK (type IN (1, 2, 3, 4, 5)),
  CONSTRAINT tennis_sessions_match_rank_check
    CHECK (match_rank IN (0, 1, 2, 3, 4, 5)),
  CONSTRAINT tennis_sessions_rating_check
    CHECK (rating >= 1 AND rating <= 5),

  INDEX idx_tennis_sessions_user_deleted_date (user_id, deleted_at, date DESC),
  INDEX idx_tennis_sessions_user_deleted_created (user_id, deleted_at, created_at DESC)
) COMMENT='打球记录表';
