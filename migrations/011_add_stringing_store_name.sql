ALTER TABLE racket_stringing_record
  ADD COLUMN store_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '穿线门店' AFTER string_name;
