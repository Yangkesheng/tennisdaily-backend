-- 用户设置表（键值结构，用户 + 设置项 + 设置值）
CREATE TABLE IF NOT EXISTS `user_settings` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '设置ID，主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `setting_key` SMALLINT NOT NULL COMMENT '设置项枚举：1 个性签名，2 默认球场，3 默认时长',
  `setting_value` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '设置值',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_settings_user_key_deleted` (`user_id`, `setting_key`, `deleted_at`),
  KEY `idx_user_settings_user_deleted` (`user_id`, `deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='用户设置表（键值）';
