-- Back up tennis_sessions/user_settings and check active setting duplicates first:
-- SELECT user_id,setting_key,COUNT(*) FROM user_settings WHERE deleted_at IS NULL
-- GROUP BY user_id,setting_key HAVING COUNT(*)>1;
-- If any rows are returned, stop; do not guess which historical setting to retain.
-- Apply ONCE, before coordinated backend/mini-program release. DDL is not transactional.
-- This migration targets MySQL 8.0.16+ (CHECK enforcement; tested on 8.4).
-- Abort BEFORE schema changes when active settings are ambiguous.
DROP TEMPORARY TABLE IF EXISTS court_location_preflight;
CREATE TEMPORARY TABLE court_location_preflight (
 active_count BIGINT NOT NULL,
 CONSTRAINT court_location_no_duplicate_settings CHECK (active_count <= 1)
);
INSERT INTO court_location_preflight(active_count)
 SELECT COUNT(*) FROM user_settings WHERE deleted_at IS NULL
 GROUP BY user_id,setting_key HAVING COUNT(*) > 1;
DROP TEMPORARY TABLE court_location_preflight;

ALTER TABLE tennis_sessions
 ADD COLUMN court_address VARCHAR(512) NOT NULL DEFAULT '',
 ADD COLUMN court_latitude DECIMAL(10,7) NULL DEFAULT NULL,
 ADD COLUMN court_longitude DECIMAL(10,7) NULL DEFAULT NULL,
 ADD COLUMN court_coordinate_system VARCHAR(16) NOT NULL DEFAULT '';
ALTER TABLE user_settings MODIFY COLUMN setting_value TEXT NOT NULL;
-- key=4 stores default CourtLocation JSON (or JSON null); existing keys stay unchanged.
-- Keep additive columns and expanded capacity on rollback; never drop saved locations.
