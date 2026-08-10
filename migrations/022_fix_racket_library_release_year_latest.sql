-- TennisDaily：代际模糊的型号按最新款上市年份补齐 release_year
-- 依据（最新一代/最新在售款的上市年份）：
--   Babolat：Boost Aero=2024、Boost Strike=2024、Boost Drive=2025、Evo Drive=2025、
--            Evo Strike=2024、Evoke=2024；Evo Aero Gen2=2026（线上已填，保持一致）
--   Dunlop：CX 200 系=2024（CX 200 Limited=2025）、FX Team 100=2026、LX Team 107=2024、
--           SX 全系=2025（2025 新模）、Tristorm=2025
--   HEAD：Radical=2023（Auxetic 2.0）、Boom MP Orlinski=2025、Squared=2026
--   Prince：Beast=2025、Neon=2025、O3 Legacy=2025、Phantom=2024、Premier=2025、
--           Ripcord/TXTZ=2025、Skulls=2025、Tour（4 代）=2025、Tour 100P=2023（末代）、
--           Warrior=2021、Classic Graphite 100=2015（复刻）
--   Tecnifibre：Fire 全系=2026、T-Fight=2022（ISO 代）、TEMPO=2024（V2 代）、
--               TF-40=2024（V3 代）、TF-X1=2024（V2 代）
--   Yonex：Muse=2026（首发）、Astrel=2020（第二代）、VCORE Alpha L=2026、Vcore 98=2026（第八代 VCORE）
-- 脚本可重复执行：所有 UPDATE 均限定 release_year = 0，幂等。

USE tennis_diary;

-- Babolat
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Babolat' AND series IN ('Boost Aero', 'Boost Strike') AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Babolat' AND series IN ('Boost Drive', 'Evo Drive') AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Babolat' AND series = 'Evo Aero' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Babolat' AND series IN ('Evo Strike', 'Evoke') AND release_year = 0;

-- Dunlop
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Dunlop' AND series = 'CX' AND model <> 'CX 200 Limited' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Dunlop' AND model = 'CX 200 Limited' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Dunlop' AND model = 'FX Team 100' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Dunlop' AND series = 'LX Team' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Dunlop' AND series = 'SX' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Dunlop' AND series = 'Tristorm' AND release_year = 0;

-- HEAD
UPDATE racket_library SET release_year = 2023
WHERE brand = 'HEAD' AND series = 'Radical' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'HEAD' AND model = 'Boom MP Orlinski Limited Edition' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'HEAD' AND series = 'Squared' AND release_year = 0;

-- Prince
UPDATE racket_library SET release_year = 2025
WHERE brand = 'Prince' AND series IN ('Beast', 'Neon', 'O3', 'Premier', 'Ripcord', 'Skulls', 'TXTZ')
  AND release_year = 0;

UPDATE racket_library SET release_year = 2015
WHERE brand = 'Prince' AND model = 'Classic Graphite 100 (Special Edition)' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Prince' AND series = 'Phantom' AND release_year = 0;

UPDATE racket_library SET release_year = 2023
WHERE brand = 'Prince' AND model = 'Tour 100P (305g)' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Prince' AND series = 'Tour' AND release_year = 0;

UPDATE racket_library SET release_year = 2021
WHERE brand = 'Prince' AND series = 'Warrior' AND release_year = 0;

-- Tecnifibre
UPDATE racket_library SET release_year = 2026
WHERE brand = 'Tecnifibre' AND series IN ('FIRE', 'Fire') AND release_year = 0;

UPDATE racket_library SET release_year = 2022
WHERE brand = 'Tecnifibre' AND series = 'T-Fight' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Tecnifibre' AND series IN ('TEMPO', 'TF-40', 'TF-X1') AND release_year = 0;

-- Yonex
UPDATE racket_library SET release_year = 2020
WHERE brand = 'Yonex' AND series = 'Astrel' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Yonex' AND series = 'Muse' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Yonex' AND model IN ('VCORE Alpha L', 'Vcore 98') AND release_year = 0;

-- 复核：预期无剩余 release_year = 0 的行
SELECT brand, COUNT(*) AS cnt, SUM(release_year = 0) AS zero_cnt
FROM racket_library
GROUP BY brand ORDER BY brand;
