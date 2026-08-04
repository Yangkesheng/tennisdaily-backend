-- TennisDaily MySQL 初始化脚本
-- 可重复执行：创建数据库、创建表，并为旧表补齐新增字段和索引。

CREATE DATABASE IF NOT EXISTS tennis_diary
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE tennis_diary;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
  openid VARCHAR(128) DEFAULT NULL COMMENT '微信OpenID，微信登录用户唯一标识',
  phone VARCHAR(32) DEFAULT NULL COMMENT '手机号',
  nickname VARCHAR(128) NOT NULL DEFAULT '' COMMENT '昵称',
  avatar_url TEXT NOT NULL COMMENT '头像地址',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  UNIQUE KEY uk_users_openid (openid),
  UNIQUE KEY idx_users_phone (phone),
  INDEX idx_users_deleted_at (deleted_at)
) COMMENT='用户表';

CREATE TABLE IF NOT EXISTS user_wechat_identities (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '微信身份绑定ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  appid VARCHAR(64) NOT NULL COMMENT '微信小程序AppID',
  openid VARCHAR(128) NOT NULL COMMENT '微信OpenID',
  unionid VARCHAR(128) DEFAULT NULL COMMENT '微信UnionID',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  UNIQUE KEY idx_user_wechat_appid_openid (appid, openid),
  INDEX idx_user_wechat_user_id (user_id),
  INDEX idx_user_wechat_unionid (unionid),
  INDEX idx_user_wechat_deleted_at (deleted_at)
) COMMENT='用户微信身份绑定表';



CREATE TABLE IF NOT EXISTS racket_brands (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '品牌ID',
  name VARCHAR(128) NOT NULL COMMENT '品牌名称，例如 Wilson、Babolat、Yonex',
  file_id VARCHAR(255) DEFAULT NULL COMMENT '品牌图片文件ID',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  UNIQUE KEY uk_racket_brands_name (name)
) COMMENT='球拍品牌表';

CREATE TABLE IF NOT EXISTS racket_series (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '系列ID',
  brand_id BIGINT NOT NULL COMMENT '品牌ID，仅逻辑关联 racket_brands.id，不使用数据库外键',
  name VARCHAR(128) NOT NULL COMMENT '系列名称，例如 Blade、Pure Drive、EZONE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

  UNIQUE KEY uk_racket_series_brand_name (brand_id, name),
  INDEX idx_racket_series_brand_id (brand_id)
) COMMENT='球拍系列表';


CREATE TABLE IF NOT EXISTS tennis_sessions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '打球记录ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',

  date DATETIME NOT NULL COMMENT '打球开始时间',
  duration_minutes INT NOT NULL DEFAULT 120 COMMENT '打球时长，单位分钟',
  rating SMALLINT NOT NULL DEFAULT 3 COMMENT '手感评分，范围1-5',

  type SMALLINT NOT NULL DEFAULT 1 COMMENT '打球类型:1双打 2单打 3训练 4单打比赛 5双打比赛',
  category SMALLINT NOT NULL DEFAULT 1 COMMENT '一级类型:1日常球局 2训练 3比赛',
  sub_category SMALLINT NOT NULL DEFAULT 2 COMMENT '二级类型:1单打/打单 2双打 3发球 4其他',
  match_rank SMALLINT NOT NULL DEFAULT 0 COMMENT '比赛成绩:0无 1冠军 2亚军 3季军 4四强 5八强 6十六强 7小组赛',

  court_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '球场名称',
  partner VARCHAR(128) NOT NULL DEFAULT '' COMMENT '搭档',
  cost DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '本次打球费用',
  racket_id BIGINT NOT NULL DEFAULT 0 COMMENT '使用球拍ID',
  shoe_name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '球鞋名称',
  note TEXT NOT NULL COMMENT '备注',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  CONSTRAINT tennis_sessions_type_check
    CHECK (type IN (1, 2, 3, 4, 5)),
  CONSTRAINT tennis_sessions_category_check
    CHECK (category IN (1, 2, 3)),
  CONSTRAINT tennis_sessions_sub_category_check
    CHECK (sub_category IN (1, 2, 3, 4)),
  CONSTRAINT tennis_sessions_match_rank_check
    CHECK (match_rank IN (0, 1, 2, 3, 4, 5)),
  CONSTRAINT tennis_sessions_rating_check
    CHECK (rating >= 1 AND rating <= 5),

  INDEX idx_tennis_sessions_user_deleted_date (user_id, deleted_at, date DESC),
  INDEX idx_tennis_sessions_user_deleted_created (user_id, deleted_at, created_at DESC),
  INDEX idx_tennis_sessions_user_racket_deleted (user_id, racket_id, deleted_at)
) COMMENT='打球记录表';

CREATE TABLE IF NOT EXISTS racket (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '球拍ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  library_id BIGINT NOT NULL DEFAULT 0 COMMENT '球拍库ID',

  name VARCHAR(100) NOT NULL COMMENT '球拍名称',
  brand VARCHAR(50) DEFAULT NULL COMMENT '品牌',
  model VARCHAR(100) DEFAULT NULL COMMENT '型号',
  status TINYINT NOT NULL DEFAULT 2 COMMENT '状态:1主力拍 2在用 3已退役',
  purchase_date DATE DEFAULT NULL COMMENT '购买日期',
  purchase_price DECIMAL(10,2) DEFAULT NULL COMMENT '购买价格',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  CONSTRAINT racket_status_check
    CHECK (status IN (1, 2, 3)),

  INDEX idx_racket_user_deleted_status (user_id, deleted_at, status),
  INDEX idx_racket_user_deleted_created (user_id, deleted_at, created_at),
  INDEX idx_racket_user_library (user_id, library_id)
) COMMENT='球拍表';

CREATE TABLE IF NOT EXISTS racket_stringing_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '穿线记录ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  racket_id BIGINT NOT NULL COMMENT '球拍ID',

  string_name VARCHAR(100) NOT NULL COMMENT '球线名称',
  store_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '穿线门店',
  vertical_tension DECIMAL(4,1) DEFAULT NULL COMMENT '竖线磅数',
  horizontal_tension DECIMAL(4,1) DEFAULT NULL COMMENT '横线磅数',
  cost DECIMAL(10,2) NOT NULL COMMENT '穿线费用',
  string_date DATETIME NOT NULL COMMENT '穿线时间',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  INDEX idx_stringing_user_deleted_racket (user_id, deleted_at, racket_id),
  INDEX idx_stringing_user_deleted_date (user_id, deleted_at, string_date),
  INDEX idx_stringing_racket_deleted_date (racket_id, deleted_at, string_date)
) COMMENT='球拍穿线历史表';

CREATE TABLE IF NOT EXISTS racket_library (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '球拍库ID',
  brand_id BIGINT NOT NULL DEFAULT 0 COMMENT '品牌ID',
  brand VARCHAR(50) NOT NULL COMMENT '品牌',
  series_id BIGINT NOT NULL DEFAULT 0 COMMENT '系列ID',
  series VARCHAR(100) NOT NULL DEFAULT '' COMMENT '系列',
  model VARCHAR(100) NOT NULL COMMENT '型号',
  release_year SMALLINT NOT NULL COMMENT '版本年份',
  weight SMALLINT DEFAULT NULL COMMENT '裸拍重量(g)',
  head_size SMALLINT DEFAULT NULL COMMENT '拍面大小(sq in)',
  string_pattern VARCHAR(64) DEFAULT NULL COMMENT '穿线模式，例如 16x19、18x20、16/19',
  file_id VARCHAR(255) DEFAULT NULL COMMENT '球拍图片文件ID',
  image_url VARCHAR(500) DEFAULT NULL COMMENT '球拍图片',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  UNIQUE KEY uk_brand_model_year (brand, model, release_year)
) COMMENT='球拍库';

DELIMITER $$

DROP PROCEDURE IF EXISTS ensure_column$$
CREATE PROCEDURE ensure_column(
  IN p_table_name VARCHAR(64),
  IN p_column_name VARCHAR(64),
  IN p_column_definition TEXT,
  IN p_after_column VARCHAR(64)
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = p_table_name
      AND COLUMN_NAME = p_column_name
  ) THEN
    SET @ddl = CONCAT('ALTER TABLE `', p_table_name, '` ADD COLUMN ', p_column_definition);
    IF p_after_column IS NOT NULL AND p_after_column <> '' THEN
      SET @ddl = CONCAT(@ddl, ' AFTER `', p_after_column, '`');
    END IF;
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END$$

DROP PROCEDURE IF EXISTS ensure_index$$
CREATE PROCEDURE ensure_index(
  IN p_table_name VARCHAR(64),
  IN p_index_name VARCHAR(64),
  IN p_index_definition TEXT
)
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = p_table_name
      AND INDEX_NAME = p_index_name
  ) THEN
    SET @ddl = CONCAT('ALTER TABLE `', p_table_name, '` ADD ', p_index_definition);
    PREPARE stmt FROM @ddl;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END$$

DELIMITER ;

-- 兼容已存在的旧 users 表。
CALL ensure_column('users', 'phone', '`phone` VARCHAR(32) NULL COMMENT ''手机号''', 'openid');
CALL ensure_column('users', 'deleted_at', '`deleted_at` DATETIME NULL COMMENT ''删除时间''', 'updated_at');
CALL ensure_index('users', 'idx_users_phone', 'UNIQUE INDEX `idx_users_phone` (`phone`)');
CALL ensure_index('users', 'idx_users_deleted_at', 'INDEX `idx_users_deleted_at` (`deleted_at`)');

-- 兼容旧版本 users.openid NOT NULL；手机号登录允许手机号用户先存在，openid 可为空。
ALTER TABLE users MODIFY COLUMN openid VARCHAR(128) DEFAULT NULL COMMENT '微信OpenID，微信登录用户唯一标识';

-- 兼容已存在的旧 tennis_sessions 表。
CALL ensure_column('tennis_sessions', 'partner', '`partner` VARCHAR(128) NOT NULL DEFAULT '''' COMMENT ''搭档''', 'court_name');
CALL ensure_column('tennis_sessions', 'category', '`category` SMALLINT NOT NULL DEFAULT 1 COMMENT ''一级类型:1日常球局 2训练 3比赛''', 'type');
CALL ensure_column('tennis_sessions', 'sub_category', '`sub_category` SMALLINT NOT NULL DEFAULT 2 COMMENT ''二级类型:1单打/打单 2双打 3发球 4其他''', 'category');
UPDATE tennis_sessions
SET category = CASE category
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
ALTER TABLE tennis_sessions MODIFY COLUMN category SMALLINT NOT NULL DEFAULT 1 COMMENT '一级类型:1日常球局 2训练 3比赛';
ALTER TABLE tennis_sessions MODIFY COLUMN sub_category SMALLINT NOT NULL DEFAULT 2 COMMENT '二级类型:1单打/打单 2双打 3发球 4其他';
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
  END
WHERE category NOT IN (1, 2, 3)
  OR sub_category NOT IN (1, 2, 3, 4)
  OR (category = 1 AND sub_category = 2);
ALTER TABLE tennis_sessions MODIFY COLUMN date DATETIME NOT NULL COMMENT '打球开始时间';
CALL ensure_column('tennis_sessions', 'racket_id', '`racket_id` BIGINT NOT NULL DEFAULT 0 COMMENT ''使用球拍ID''', 'cost');
CALL ensure_index('tennis_sessions', 'idx_tennis_sessions_user_deleted_date', 'INDEX `idx_tennis_sessions_user_deleted_date` (`user_id`, `deleted_at`, `date` DESC)');
CALL ensure_index('tennis_sessions', 'idx_tennis_sessions_user_deleted_created', 'INDEX `idx_tennis_sessions_user_deleted_created` (`user_id`, `deleted_at`, `created_at` DESC)');
CALL ensure_index('tennis_sessions', 'idx_tennis_sessions_user_racket_deleted', 'INDEX `idx_tennis_sessions_user_racket_deleted` (`user_id`, `racket_id`, `deleted_at`)');

-- 兼容已存在的旧 racket 表。
CALL ensure_column('racket', 'library_id', '`library_id` BIGINT NOT NULL DEFAULT 0 COMMENT ''球拍库ID''', 'user_id');
CALL ensure_index('racket', 'idx_racket_user_deleted_status', 'INDEX `idx_racket_user_deleted_status` (`user_id`, `deleted_at`, `status`)');
CALL ensure_index('racket', 'idx_racket_user_deleted_created', 'INDEX `idx_racket_user_deleted_created` (`user_id`, `deleted_at`, `created_at`)');
CALL ensure_index('racket', 'idx_racket_user_library', 'INDEX `idx_racket_user_library` (`user_id`, `library_id`)');

-- 兼容旧版本 racket_stringing_record.string_date 只存日期。
ALTER TABLE racket_stringing_record MODIFY COLUMN string_date DATETIME NOT NULL COMMENT '穿线时间';
CALL ensure_column('racket_stringing_record', 'store_name', '`store_name` VARCHAR(100) NOT NULL DEFAULT '''' COMMENT ''穿线门店''', 'string_name');
CALL ensure_column('racket_stringing_record', 'vertical_tension', '`vertical_tension` DECIMAL(4,1) DEFAULT NULL COMMENT ''竖线磅数''', 'store_name');
CALL ensure_column('racket_stringing_record', 'horizontal_tension', '`horizontal_tension` DECIMAL(4,1) DEFAULT NULL COMMENT ''横线磅数''', 'vertical_tension');

-- 兼容已存在的旧 racket_library 表。
CALL ensure_column('racket_library', 'brand_id', '`brand_id` BIGINT NOT NULL DEFAULT 0 COMMENT ''品牌ID''', 'id');
CALL ensure_column('racket_library', 'series_id', '`series_id` BIGINT NOT NULL DEFAULT 0 COMMENT ''系列ID''', 'brand');
CALL ensure_column('racket_library', 'series', '`series` VARCHAR(100) NOT NULL DEFAULT '''' COMMENT ''系列''', 'series_id');
CALL ensure_column('racket_library', 'string_pattern', '`string_pattern` VARCHAR(64) NULL DEFAULT NULL COMMENT ''穿线模式，例如 16x19、18x20、16/19''', 'head_size');
CALL ensure_column('racket_library', 'file_id', '`file_id` VARCHAR(255) NULL DEFAULT NULL COMMENT ''球拍图片文件ID''', 'string_pattern');

-- 兼容已存在但索引不完整的表。
CALL ensure_index('user_wechat_identities', 'idx_user_wechat_appid_openid', 'UNIQUE INDEX `idx_user_wechat_appid_openid` (`appid`, `openid`)');
CALL ensure_index('user_wechat_identities', 'idx_user_wechat_user_id', 'INDEX `idx_user_wechat_user_id` (`user_id`)');
CALL ensure_index('user_wechat_identities', 'idx_user_wechat_unionid', 'INDEX `idx_user_wechat_unionid` (`unionid`)');
CALL ensure_index('racket_stringing_record', 'idx_stringing_user_deleted_racket', 'INDEX `idx_stringing_user_deleted_racket` (`user_id`, `deleted_at`, `racket_id`)');
CALL ensure_index('racket_stringing_record', 'idx_stringing_user_deleted_date', 'INDEX `idx_stringing_user_deleted_date` (`user_id`, `deleted_at`, `string_date`)');
CALL ensure_index('racket_stringing_record', 'idx_stringing_racket_deleted_date', 'INDEX `idx_stringing_racket_deleted_date` (`racket_id`, `deleted_at`, `string_date`)');
CALL ensure_index('racket_library', 'uk_brand_model_year', 'UNIQUE INDEX `uk_brand_model_year` (`brand`, `model`, `release_year`)');
CALL ensure_index('racket_brands', 'uk_racket_brands_name', 'UNIQUE INDEX `uk_racket_brands_name` (`name`)');
CALL ensure_index('racket_series', 'uk_racket_series_brand_name', 'UNIQUE INDEX `uk_racket_series_brand_name` (`brand_id`, `name`)');
CALL ensure_index('racket_series', 'idx_racket_series_brand_id', 'INDEX `idx_racket_series_brand_id` (`brand_id`)');

DROP PROCEDURE IF EXISTS ensure_column;
DROP PROCEDURE IF EXISTS ensure_index;
