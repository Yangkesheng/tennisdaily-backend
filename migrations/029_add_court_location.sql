-- Back up tennis_sessions/user_settings before applying. DDL is not transactional.
-- No TEMPORARY TABLEs: web SQL consoles may use different connections per statement.
-- Send CREATE PROCEDURE ... END as ONE statement. With mysql CLI, wrap that
-- statement in DELIMITER $$ / DELIMITER ;. Go dbclient with multiStatements=true
-- can execute this whole file directly. The routine is database-scoped.
-- Re-running is supported, including after a partially completed ALTER.
DROP PROCEDURE IF EXISTS migrate_court_location_029;
CREATE PROCEDURE migrate_court_location_029()
BEGIN
 IF EXISTS (
  SELECT 1 FROM user_settings WHERE deleted_at IS NULL
  GROUP BY user_id,setting_key HAVING COUNT(*) > 1
 ) THEN
  SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Duplicate active user settings; migration stopped before schema changes';
 END IF;

 IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='tennis_sessions' AND COLUMN_NAME='court_address') THEN
  ALTER TABLE tennis_sessions ADD COLUMN court_address VARCHAR(512) NOT NULL DEFAULT '';
 END IF;
 IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='tennis_sessions' AND COLUMN_NAME='court_latitude') THEN
  ALTER TABLE tennis_sessions ADD COLUMN court_latitude DECIMAL(10,7) NULL DEFAULT NULL;
 END IF;
 IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='tennis_sessions' AND COLUMN_NAME='court_longitude') THEN
  ALTER TABLE tennis_sessions ADD COLUMN court_longitude DECIMAL(10,7) NULL DEFAULT NULL;
 END IF;
 IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='tennis_sessions' AND COLUMN_NAME='court_coordinate_system') THEN
  ALTER TABLE tennis_sessions ADD COLUMN court_coordinate_system VARCHAR(16) NOT NULL DEFAULT '';
 END IF;
 IF NOT EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='user_settings' AND COLUMN_NAME='setting_value' AND DATA_TYPE='text' AND IS_NULLABLE='NO' AND COLUMN_DEFAULT IS NULL) THEN
  ALTER TABLE user_settings MODIFY COLUMN setting_value TEXT NOT NULL;
 END IF;
END;
CALL migrate_court_location_029();
DROP PROCEDURE IF EXISTS migrate_court_location_029;
-- key=4 stores default CourtLocation JSON (or JSON null).
-- Keep additive columns and expanded capacity on rollback; never drop saved locations.
