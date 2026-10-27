-- Backup before applying. Historical string_name is intentionally untouched.
-- New clients must send lines.main/cross; deploy with the new mini-program.
CREATE TABLE IF NOT EXISTS racket_stringing_lines (
 id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
 stringing_record_id BIGINT NOT NULL,
 axis VARCHAR(10) NOT NULL,
 string_id BIGINT NULL,
 variant_id BIGINT NULL,
 name_snapshot VARCHAR(200) NOT NULL,
 diameter_mm_snapshot DECIMAL(4,2) NULL,
 color_snapshot VARCHAR(50) NOT NULL DEFAULT '',
 UNIQUE KEY uk_stringing_axis (stringing_record_id, axis),
 KEY idx_stringing_string (string_id, stringing_record_id),
 CONSTRAINT fk_line_record FOREIGN KEY (stringing_record_id) REFERENCES racket_stringing_record(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- Do not drop these tables on rollback after new records have been written.
