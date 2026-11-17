CREATE TABLE IF NOT EXISTS `analytics_installations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `installation_id` VARCHAR(64) NOT NULL,
  `user_id` BIGINT DEFAULT NULL,
  `first_seen_at` DATETIME NOT NULL,
  `last_seen_at` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_analytics_installation_id` (`installation_id`),
  KEY `idx_analytics_installations_last_seen_user` (`last_seen_at`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='匿名小程序安装活跃标识';

CREATE TABLE IF NOT EXISTS `analytics_installation_days` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `installation_id` VARCHAR(64) NOT NULL,
  `activity_date` DATE NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_analytics_installation_day` (`installation_id`, `activity_date`),
  KEY `idx_analytics_installation_days_date` (`activity_date`, `installation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='匿名安装每日去重打开';
