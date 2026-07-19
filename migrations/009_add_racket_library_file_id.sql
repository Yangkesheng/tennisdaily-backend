ALTER TABLE racket_library
  ADD COLUMN file_id VARCHAR(255) NULL DEFAULT NULL COMMENT '球拍图片文件ID' AFTER string_pattern;
