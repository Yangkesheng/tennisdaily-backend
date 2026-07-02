ALTER TABLE racket_stringing_record
  ADD COLUMN vertical_tension DECIMAL(4,1) DEFAULT NULL COMMENT '竖线磅数' AFTER tension,
  ADD COLUMN horizontal_tension DECIMAL(4,1) DEFAULT NULL COMMENT '横线磅数' AFTER vertical_tension;
