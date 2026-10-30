-- Run once before deploying the API. No historical data conversion.
-- Replaces the unexecuted JSON-column proposal; no store_location JSON column is used.
ALTER TABLE racket_stringing_record
 ADD COLUMN store_address VARCHAR(512) NOT NULL DEFAULT '',
 ADD COLUMN store_latitude DECIMAL(10,7) NULL DEFAULT NULL,
 ADD COLUMN store_longitude DECIMAL(10,7) NULL DEFAULT NULL,
 ADD COLUMN store_coordinate_system VARCHAR(16) NOT NULL DEFAULT '';
