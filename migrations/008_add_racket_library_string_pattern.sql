ALTER TABLE racket_library
  ADD COLUMN string_pattern VARCHAR(64) NULL DEFAULT NULL COMMENT '穿线模式，例如 16x19、18x20、16/19' AFTER head_size;
