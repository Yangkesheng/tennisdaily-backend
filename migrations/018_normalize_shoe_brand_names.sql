-- 品牌名规范化：修正大小写不一致，避免 stats 按 (brand, brand_id) 分组时同一品牌被拆成多条。

-- 1. shoe_brands 主表：ASICS 官方写法为全大写
UPDATE `shoe_brands` SET `name` = 'ASICS' WHERE `id` = 3;

-- 2. K-Swiss 重复品牌合并：id=8（K-Swiss）未被任何数据引用，删除后把 id=37（KSwiss）规范为官方写法
DELETE FROM `shoe_brands`
WHERE `id` = 8
  AND NOT EXISTS (SELECT 1 FROM `shoe_library` WHERE `brand_id` = 8)
  AND NOT EXISTS (SELECT 1 FROM `shoe_series` WHERE `brand_id` = 8);
UPDATE `shoe_brands` SET `name` = 'K-Swiss', `slug` = 'k-swiss' WHERE `id` = 37;

-- 3. shoe_library.brand 与主表对齐（幂等，云端存在其他写法时也会被修正）
UPDATE `shoe_library` SET `brand` = 'ASICS' WHERE `brand_id` = 3;
UPDATE `shoe_library` SET `brand` = 'HEAD' WHERE `brand_id` = 17;
UPDATE `shoe_library` SET `brand` = 'K-Swiss' WHERE `brand_id` = 37;

-- 4. 若存在引用 brand_id=8 的历史行，一并归并到 37
UPDATE `shoe_library` SET `brand_id` = 37, `brand` = 'K-Swiss' WHERE `brand_id` = 8;
