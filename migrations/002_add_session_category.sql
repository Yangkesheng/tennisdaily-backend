ALTER TABLE tennis_sessions
  ADD COLUMN category SMALLINT NOT NULL DEFAULT 1 COMMENT '一级类型:1日常球局 2训练 3比赛' AFTER type,
  ADD COLUMN sub_category SMALLINT NOT NULL DEFAULT 2 COMMENT '二级类型:1单打/打单 2双打 3发球 4其他' AFTER category;

UPDATE tennis_sessions
SET category = CASE type
    WHEN 3 THEN 2
    WHEN 4 THEN 3
    WHEN 5 THEN 3
    ELSE 1
  END,
  sub_category = CASE type
    WHEN 1 THEN 2
    WHEN 2 THEN 1
    WHEN 3 THEN 4
    WHEN 4 THEN 1
    WHEN 5 THEN 2
    ELSE 2
  END;
