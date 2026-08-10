-- TennisDaily：补齐 Babolat / Dunlop / HEAD / Prince / Tecnifibre / Yonex 的 release_year
-- 年份依据（仅填有把握的系列/版本，代际模糊的型号保持 0）：
--   HEAD Graphene XT = 2015、Boom Pro = 2021
--   Yonex Percept = 2024
--   Tecnifibre TF-40 V3 = 2024、T-Fight ISO = 2022、TF-X1 V2 = 2024、TEMPO V2 = 2024
--   Dunlop CX 200 Tour 16x19 = 2024、FX 500 系列（2026 代）= 2026
--   Prince Synergy 98 = 2021、Beast 100 (265g) 及 Smiley 限定 = 2025
--   Babolat Pure Drive Gen11 = 2025、Pure Drive Spectra = 2026、Pure Aero 2026 系列 = 2026
-- 脚本可重复执行：所有 UPDATE 均限定 release_year = 0，幂等。

USE tennis_diary;

-- Blade 100L V9 Rouge Limited Edition（美网限定红配色，2025 年上市；id=517 与 id=679 为不同配色，均保留）
UPDATE racket_library SET model = 'Blade 100L V9 Rouge', release_year = 2025
WHERE id = 517 AND release_year = 0;

-- HEAD
UPDATE racket_library SET release_year = 2015
WHERE brand = 'HEAD' AND model LIKE '%Graphene XT%' AND release_year = 0;

UPDATE racket_library SET release_year = 2021
WHERE brand = 'HEAD' AND model = 'Boom Pro' AND release_year = 0;

-- Yonex
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Yonex' AND series = 'Percept' AND release_year = 0;

-- Tecnifibre
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Tecnifibre' AND series = 'TF-40' AND model LIKE '%V 3%' AND release_year = 0;

UPDATE racket_library SET release_year = 2022
WHERE brand = 'Tecnifibre' AND series = 'T-Fight' AND model LIKE '%ISO%' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Tecnifibre' AND series = 'TF-X1' AND model LIKE '%V2%' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Tecnifibre' AND series = 'TEMPO' AND model LIKE '%V2%' AND release_year = 0;

-- Dunlop
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Dunlop' AND model = 'CX 200 Tour 16x19' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Dunlop' AND series = 'FX'
  AND model IN ('FX 500', 'FX 500 Lite', 'FX 500 LS', 'FX 500 Tour', 'FX 500 Tour Tour racket Testracket')
  AND release_year = 0;

-- Prince
UPDATE racket_library SET release_year = 2021
WHERE brand = 'Prince' AND model = 'Synergy 98' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Prince' AND series = 'Beast'
  AND (model LIKE '%Smiley%' OR model = 'Beast 100 (265g)') AND release_year = 0;

-- Babolat
UPDATE racket_library SET release_year = 2025
WHERE brand = 'Babolat' AND series = 'Pure Drive'
  AND model IN ('Pure Drive', 'Pure Drive +', 'Pure Drive 98', 'Pure Drive Lite',
                'Pure Drive Team', 'Pure Drive 107', 'Pure Drive Wimbledon')
  AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Babolat' AND series = 'Pure Drive' AND model LIKE '%Spectra%' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Babolat' AND series = 'Pure Aero'
  AND model IN ('Pure Aero', 'Pure Aero +', 'Pure Aero Lite', 'Pure Aero Team')
  AND release_year = 0;

-- 复核：剩余 release_year = 0 的行（应为代际模糊/无法确认的型号）
SELECT brand, series, model, id
FROM racket_library
WHERE release_year = 0
ORDER BY brand, series, model;
