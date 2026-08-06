CREATE TABLE IF NOT EXISTS shoe_library_image_upload (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
  shoe_library_id BIGINT NOT NULL COMMENT '球鞋库ID',
  file_id VARCHAR(255) NOT NULL COMMENT '对象存储文件ID',
  object_key VARCHAR(255) NOT NULL COMMENT '对象存储路径',
  source_url VARCHAR(500) NOT NULL DEFAULT '' COMMENT '迁移前的图片URL',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',

  UNIQUE KEY uk_shoe_library_id (shoe_library_id)
) COMMENT='球鞋库图片上传记录';
