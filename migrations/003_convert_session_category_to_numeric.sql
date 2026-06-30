UPDATE tennis_sessions
SET category = CASE category
    WHEN '1' THEN '1'
    WHEN '2' THEN '2'
    WHEN '3' THEN '3'
    WHEN 'training' THEN '2'
    WHEN 'match' THEN '3'
    WHEN 'daily' THEN '1'
    ELSE CASE type
      WHEN 3 THEN '2'
      WHEN 4 THEN '3'
      WHEN 5 THEN '3'
      ELSE '1'
    END
  END,
  sub_category = CASE sub_category
    WHEN '1' THEN '1'
    WHEN '2' THEN '2'
    WHEN '3' THEN '3'
    WHEN '4' THEN '4'
    WHEN 'singles' THEN '1'
    WHEN 'doubles' THEN '2'
    WHEN 'serve' THEN '3'
    WHEN 'other' THEN '4'
    ELSE CASE type
      WHEN 1 THEN '2'
      WHEN 2 THEN '1'
      WHEN 3 THEN '4'
      WHEN 4 THEN '1'
      WHEN 5 THEN '2'
      ELSE '2'
    END
  END;

ALTER TABLE tennis_sessions
  MODIFY COLUMN category ENUM('1','2','3') NOT NULL DEFAULT '1' COMMENT '一级类型:1日常球局 2训练 3比赛',
  MODIFY COLUMN sub_category ENUM('1','2','3','4') NOT NULL DEFAULT '2' COMMENT '二级类型:1单打/打单 2双打 3发球 4其他';
