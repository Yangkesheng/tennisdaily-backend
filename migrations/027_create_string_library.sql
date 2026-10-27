-- Independent catalog only. No changes to racket_stringing_record.
-- Apply once after backup. No unverified product data is inserted.
CREATE TABLE IF NOT EXISTS string_library (
 id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
 brand VARCHAR(100) NOT NULL,
 model VARCHAR(150) NOT NULL,
 name_zh VARCHAR(150) NOT NULL DEFAULT '',
 aliases VARCHAR(1000) NOT NULL DEFAULT '',
 material VARCHAR(50) NOT NULL DEFAULT '',
 source_url VARCHAR(1000) NOT NULL,
 published BOOLEAN NOT NULL DEFAULT FALSE,
 created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
 UNIQUE KEY uk_string_brand_model (brand, model),
 KEY idx_string_published_brand (published, brand)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS string_variants (
 id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
 string_id BIGINT NOT NULL,
 diameter_mm DECIMAL(4,2) NOT NULL,
 UNIQUE KEY uk_string_diameter (string_id, diameter_mm),
 CONSTRAINT fk_string_variant FOREIGN KEY (string_id) REFERENCES string_library(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
