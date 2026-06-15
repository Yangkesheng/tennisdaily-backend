ALTER TABLE tennis_sessions
  ADD COLUMN partner VARCHAR(128) NOT NULL DEFAULT '' COMMENT '搭档' AFTER court_name;
