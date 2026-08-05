-- 移除 shoe_library 发布状态：取消 draft/published/archived 概念，
-- 库内数据默认全部可被 C 端读取，无需发布流程。
SET @col_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'shoe_library'
    AND COLUMN_NAME = 'status'
);
SET @ddl := IF(
  @col_exists > 0,
  'ALTER TABLE `shoe_library` DROP COLUMN `status`',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
