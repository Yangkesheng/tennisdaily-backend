ALTER TABLE users
  ADD COLUMN phone VARCHAR(32) NULL AFTER openid,
  ADD COLUMN deleted_at DATETIME NULL AFTER updated_at;

CREATE UNIQUE INDEX idx_users_phone ON users(phone);

CREATE TABLE IF NOT EXISTS user_wechat_identities (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  appid VARCHAR(64) NOT NULL,
  openid VARCHAR(128) NOT NULL,
  unionid VARCHAR(128) DEFAULT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,

  UNIQUE KEY idx_user_wechat_appid_openid (appid, openid),
  INDEX idx_user_wechat_user_id (user_id),
  INDEX idx_user_wechat_unionid (unionid)
) COMMENT='用户微信身份绑定表';
