-- SQL-console alternative to 029_add_court_location.sql; DO NOT run both blindly.
-- Select the business database and back up tennis_sessions/user_settings first.
-- Step 1: run this ONE SELECT. It returns SQL text and does not modify any table.
-- Step 2: execute each sql_to_execute ONLY if all statuses are READY.
-- BLOCKED means stop and investigate duplicate settings. DONE means no changes needed.
-- Re-run this SELECT after execution to verify DONE; safe after partial execution.
SELECT
 CASE WHEN problems.n>0 THEN 'BLOCKED: duplicate active user settings'
      WHEN changes.sql_to_execute IS NULL THEN 'DONE'
      ELSE 'READY' END AS migration_status,
 CASE WHEN problems.n>0 THEN NULL ELSE changes.sql_to_execute END AS sql_to_execute
FROM (
 SELECT COUNT(*) AS n FROM (
  SELECT user_id,setting_key FROM user_settings WHERE deleted_at IS NULL
  GROUP BY user_id,setting_key HAVING COUNT(*)>1
 ) duplicates
) problems
LEFT JOIN (
 SELECT d.sequence_no,CONCAT('ALTER TABLE tennis_sessions ADD COLUMN ',d.column_ddl) AS sql_to_execute
 FROM (
  SELECT 1 AS sequence_no,'court_address' AS column_name,
   'court_address VARCHAR(512) NOT NULL DEFAULT ''''' AS column_ddl
  UNION ALL SELECT 2,'court_latitude','court_latitude DECIMAL(10,7) NULL DEFAULT NULL'
  UNION ALL SELECT 3,'court_longitude','court_longitude DECIMAL(10,7) NULL DEFAULT NULL'
  UNION ALL SELECT 4,'court_coordinate_system','court_coordinate_system VARCHAR(16) NOT NULL DEFAULT '''''
 ) d
 LEFT JOIN INFORMATION_SCHEMA.COLUMNS c ON c.TABLE_SCHEMA=DATABASE()
  AND c.TABLE_NAME='tennis_sessions' AND c.COLUMN_NAME=d.column_name
 WHERE c.COLUMN_NAME IS NULL
 UNION ALL
 SELECT 5,'ALTER TABLE user_settings MODIFY COLUMN setting_value TEXT NOT NULL'
 FROM INFORMATION_SCHEMA.COLUMNS c
 WHERE c.TABLE_SCHEMA=DATABASE() AND c.TABLE_NAME='user_settings' AND c.COLUMN_NAME='setting_value'
  AND (c.DATA_TYPE<>'text' OR c.IS_NULLABLE<>'NO' OR c.COLUMN_DEFAULT IS NOT NULL)
) changes ON 1=1
ORDER BY changes.sequence_no;
