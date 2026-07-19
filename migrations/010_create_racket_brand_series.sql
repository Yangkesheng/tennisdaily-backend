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
