CREATE TABLE IF NOT EXISTS racket (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '球拍ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',

  name VARCHAR(100) NOT NULL COMMENT '球拍名称',
  brand VARCHAR(50) DEFAULT NULL COMMENT '品牌',
  model VARCHAR(100) DEFAULT NULL COMMENT '型号',
  status TINYINT NOT NULL DEFAULT 2 COMMENT '状态:1主力拍 2在用 3已退役',
  image_url VARCHAR(500) DEFAULT NULL COMMENT '图片',
  purchase_date DATE DEFAULT NULL COMMENT '购买日期',
  purchase_price DECIMAL(10,2) DEFAULT NULL COMMENT '购买价格',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  CONSTRAINT racket_status_check
    CHECK (status IN (1, 2, 3)),

  INDEX idx_racket_user_deleted_status (user_id, deleted_at, status),
  INDEX idx_racket_user_deleted_created (user_id, deleted_at, created_at)
) COMMENT='球拍表';

CREATE TABLE IF NOT EXISTS racket_stringing_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '穿线记录ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  racket_id BIGINT NOT NULL COMMENT '球拍ID',

  string_name VARCHAR(100) NOT NULL COMMENT '球线名称',
  tension DECIMAL(4,1) DEFAULT NULL COMMENT '磅数',
  cost DECIMAL(10,2) NOT NULL COMMENT '穿线费用',
  string_date DATE NOT NULL COMMENT '穿线日期',

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME DEFAULT NULL COMMENT '删除时间',

  INDEX idx_stringing_user_deleted_racket (user_id, deleted_at, racket_id),
  INDEX idx_stringing_user_deleted_date (user_id, deleted_at, string_date),
  INDEX idx_stringing_racket_deleted_date (racket_id, deleted_at, string_date)
) COMMENT='球拍穿线历史表';



DROP TABLE IF EXISTS racket_library;

CREATE TABLE IF NOT EXISTS racket_library (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '球拍库ID',
  brand VARCHAR(50) NOT NULL COMMENT '品牌',
  model VARCHAR(100) NOT NULL COMMENT '型号',
  release_year SMALLINT NOT NULL COMMENT '版本年份',
  weight SMALLINT DEFAULT NULL COMMENT '裸拍重量(g)',
  head_size SMALLINT DEFAULT NULL COMMENT '拍面大小(sq in)',
  image_url VARCHAR(500) DEFAULT NULL COMMENT '球拍图片',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  UNIQUE KEY uk_brand_model_year (brand, model, release_year)
) COMMENT='球拍库';
