-- 打球记录支持关联我的球鞋（shoeId）
-- 展示名称 shoeName 优先按 shoe_id 关联 shoe.name 实时返回；无 shoe_id 时回退到快照字段 shoe_name
ALTER TABLE tennis_sessions ADD COLUMN shoe_id BIGINT NOT NULL DEFAULT 0 COMMENT '使用球鞋ID' AFTER racket_id;
