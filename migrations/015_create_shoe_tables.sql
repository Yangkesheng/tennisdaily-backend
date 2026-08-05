-- 球鞋库相关表（球鞋品牌、球鞋系列、球鞋库）
-- 业务主表为 shoe_library；colorway 不同配色各占一行，product_code 仅用于回源关联 source_shoes

CREATE TABLE IF NOT EXISTS `shoe_brands` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID，主键',
  `name` VARCHAR(128) NOT NULL COMMENT '品牌名称，例如 Nike、adidas、Asics',
  `slug` VARCHAR(128) NOT NULL COMMENT '品牌标识，例如 nike、adidas、asics',
  `file_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '品牌图片文件ID',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shoe_brands_name` (`name`),
  UNIQUE KEY `uk_shoe_brands_slug` (`slug`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='球鞋品牌表';

CREATE TABLE IF NOT EXISTS `shoe_series` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '系列ID，主键',
  `brand_id` BIGINT UNSIGNED NOT NULL COMMENT '品牌ID，仅逻辑关联 shoe_brands.id，不使用数据库外键',
  `gender` SMALLINT NOT NULL DEFAULT 0 COMMENT '性别：0未知 1男 2女 3童',
  `name` VARCHAR(128) NOT NULL COMMENT '系列名称，例如 Vapor、Court FF、Gel Resolution',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shoe_series_brand_gender_name` (`brand_id`, `gender`, `name`),
  KEY `idx_shoe_series_brand_id` (`brand_id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='球鞋系列表';

CREATE TABLE IF NOT EXISTS `shoe_library` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '库表ID，主键',
  `brand_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '品牌ID，仅逻辑关联 shoe_brands.id',
  `brand` VARCHAR(50) NOT NULL COMMENT '品牌名称冗余字段，例如 Nike',
  `series_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '系列ID，仅逻辑关联 shoe_series.id',
  `series` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '系列名称冗余字段，例如 Vapor',
  `model` VARCHAR(100) NOT NULL COMMENT '型号，例如 Vapor 12',
  `gender` SMALLINT NOT NULL DEFAULT 0 COMMENT '性别：0未知 1男 2女 3童',
  `colorway` VARCHAR(128) NULL DEFAULT NULL COMMENT '配色，例如 White/Orewood Brown，不同配色各占一行',
  `product_code` VARCHAR(64) NULL DEFAULT NULL COMMENT '来源商品代码，仅作为回源关联 source_shoes 的列，不是去重键',
  `release_year` INT NOT NULL DEFAULT 0 COMMENT '上市年份',
  `weight` VARCHAR(128) NULL DEFAULT NULL COMMENT '单只重量，例如 14.6 ounces (size 10.5)',
  `width` VARCHAR(128) NULL DEFAULT NULL COMMENT '鞋宽（版型）：Snug Medium、Wide、Narrow',
  `surface` VARCHAR(255) NULL DEFAULT NULL COMMENT '适用场地：Hard (all court)、Clay 等',
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '当前售价',
  `colorway_count` INT NOT NULL DEFAULT 0 COMMENT '该型号在售配色数量',
  `file_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '首图云存储fileID',
  `image_url` VARCHAR(500) NULL DEFAULT NULL COMMENT '首图访问URL',
  `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '状态：draft、published、archived',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_shoe_library_brand_id` (`brand_id`),
  KEY `idx_shoe_library_series_id` (`series_id`),
  KEY `idx_shoe_library_brand` (`brand`),
  KEY `idx_shoe_library_model` (`model`),
  KEY `idx_shoe_library_gender` (`gender`),
  KEY `idx_shoe_library_colorway` (`colorway`),
  KEY `idx_shoe_library_product_code` (`product_code`),
  KEY `idx_shoe_library_surface` (`surface`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='球鞋库表（业务主表）';

-- 我的球鞋表（用户维度：尺码、配色、状态、购买信息；软删除）
CREATE TABLE IF NOT EXISTS `my_shoes` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '球鞋ID，主键',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `library_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '球鞋库配色行ID，仅逻辑关联 shoe_library.id',

  PRIMARY KEY (`id`),
  `name` VARCHAR(100) NOT NULL COMMENT '球鞋名称',
  `brand` VARCHAR(50) DEFAULT NULL COMMENT '品牌',
  `model` VARCHAR(100) DEFAULT NULL COMMENT '型号',
  `status` TINYINT NOT NULL DEFAULT 2 COMMENT '状态:1主力鞋 2在用 3已退役',
  `size` VARCHAR(20) DEFAULT NULL COMMENT '尺码，例如 42、42.5',
  `colorway` VARCHAR(128) DEFAULT NULL COMMENT '配色，从库选择时默认带入，可修改',
  `purchase_date` DATE DEFAULT NULL COMMENT '购买日期',
  `purchase_price` DECIMAL(10,2) DEFAULT NULL COMMENT '购买价格',

  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间',

  CONSTRAINT `my_shoes_status_check` CHECK (`status` IN (1, 2, 3)),
  KEY `idx_my_shoes_user_deleted_status` (`user_id`, `deleted_at`, `status`),
  KEY `idx_my_shoes_user_deleted_created` (`user_id`, `deleted_at`, `created_at`),
  KEY `idx_my_shoes_user_library` (`user_id`, `library_id`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COMMENT='我的球鞋表';
