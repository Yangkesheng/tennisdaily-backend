-- 取消打球记录的球拍名称快照字段
-- 展示球拍名称改为按 racket_id 关联我的球拍（racket.name）实时查询
ALTER TABLE tennis_sessions DROP COLUMN racket_name;
