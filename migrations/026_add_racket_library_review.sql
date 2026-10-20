-- Pause external imports before applying. Run once before deploying the new API.
-- DEFAULT 1 initializes existing rows without an UPDATE that could re-approve future imports.
ALTER TABLE racket_library
  ADD COLUMN is_reviewed TINYINT NOT NULL DEFAULT 1 COMMENT '0待审核 1已审核',
  ADD COLUMN reviewed_by BIGINT NULL COMMENT '审核管理员ID；历史初始化为空',
  ADD COLUMN reviewed_at DATETIME NULL COMMENT '审核时间',
  ADD INDEX idx_review_brand_series (is_reviewed, brand_id, series_id);
-- Future inserts, including external imports that omit this column, are pending.
ALTER TABLE racket_library ALTER COLUMN is_reviewed SET DEFAULT 0;
