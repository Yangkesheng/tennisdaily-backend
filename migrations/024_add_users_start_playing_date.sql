-- 用户表新增“开始打球年月”
ALTER TABLE `users`
  ADD COLUMN `start_playing_date` INT DEFAULT NULL COMMENT '开始打球年月，格式 YYYYMM，如 202308' AFTER `avatar_url`;
