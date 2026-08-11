-- 意见反馈表
CREATE TABLE IF NOT EXISTS `feedback` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '反馈ID，主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `content` TEXT NOT NULL COMMENT '反馈内容',
  `contact` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '联系方式（选填）',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '处理状态：0待处理 1已处理',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_feedback_user_deleted_created` (`user_id`, `deleted_at`, `created_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='意见反馈表';
