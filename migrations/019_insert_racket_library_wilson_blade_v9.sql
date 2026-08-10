-- TennisDaily：向球拍库插入 Wilson Blade V9 系列
-- 排除型号：Blade 100L V9、Blade 104 V9、Blade 98 (16x19) V9、
--          Blade 98L V9、Blade 98S V9、Noir Blade 98 V9
-- 脚本可重复执行：品牌/系列不存在时自动创建；同型号已存在时更新规格与图片。

USE tennis_diary;

-- 1. 确保品牌与系列存在
INSERT IGNORE INTO racket_brands (name) VALUES ('Wilson');

INSERT IGNORE INTO racket_series (brand_id, name)
SELECT id, 'Blade'
FROM racket_brands
WHERE name = 'Wilson';

-- 2. 插入球拍库（品牌/系列 ID 动态关联，不需要手工指定）
INSERT INTO racket_library
  (brand_id, brand, series_id, series, model, release_year, weight, head_size, string_pattern, file_id)
SELECT
  b.id,
  'Wilson',
  s.id,
  'Blade',
  t.model,
  2024,
  t.weight,
  t.head_size,
  t.string_pattern,
  'cloud://prod-d6gkg1meqce4aaa7e.7072-prod-d6gkg1meqce4aaa7e-1440725645/racket_library/Wilson/Blade/Wilson-Blade-V9.jpg'
FROM racket_brands b
JOIN racket_series s
  ON s.brand_id = b.id
 AND s.name = 'Blade'
CROSS JOIN (
  SELECT 'Blade 98 (18x20) V9' AS model, 305 AS weight, 98 AS head_size, '18x20' AS string_pattern
  UNION ALL SELECT 'Blade 100 V9', 300, 100, '16x19'
  UNION ALL SELECT 'Blade 100UL V9', 265, 100, '16x19'
  UNION ALL SELECT 'Blade 101L V9', 274, 101, '16x20'
  UNION ALL SELECT 'Blade Pro 98 (16x19) V9', 305, 98, '16x19'
  UNION ALL SELECT 'Blade Pro 98 (18x20) V9', 305, 98, '18x20'
) t
WHERE b.name = 'Wilson'
ON DUPLICATE KEY UPDATE
  brand_id = VALUES(brand_id),
  series_id = VALUES(series_id),
  weight = VALUES(weight),
  head_size = VALUES(head_size),
  string_pattern = VALUES(string_pattern),
  file_id = VALUES(file_id);
