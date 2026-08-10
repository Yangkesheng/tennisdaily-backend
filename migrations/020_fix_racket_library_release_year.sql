-- TennisDaily：修复 racket_library.release_year 缺失 + string_pattern 格式统一
-- 年份依据：Wilson 各系列/版本的上市年份（V2=2022、V3=2025、V4=2020、V5=2023、
--           V6=2026、V8=2023、V9=2024、V10=2026、V14=2023、RF=2024、Shift=2023、
--           Defyer=2026、Pro Open=2025、Blade Feel=2025、Classic=2026）
-- 脚本可重复执行：所有 UPDATE 均限定 release_year = 0，幂等。

USE tennis_diary;

-- Blade
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Wilson' AND series = 'Blade' AND model LIKE '%V9%'
  AND release_year = 0
  AND model NOT LIKE '%US Open%'
  AND model NOT IN ('Blade 100L V9');  -- 517 与 679 为不同配色，保留待人工确认

UPDATE racket_library SET release_year = 2025
WHERE model = 'Blade 98 16X19 V9 US Open' AND release_year = 0;

UPDATE racket_library SET release_year = 2024
WHERE model = 'Noir Blade 98 V9' AND release_year = 0;

UPDATE racket_library SET release_year = 2023
WHERE model LIKE 'Blade 100L V8%' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'Blade' AND model LIKE 'Blade Feel%' AND release_year = 0;

-- Clash
UPDATE racket_library SET release_year = 2022
WHERE brand = 'Wilson' AND series = 'Clash' AND model LIKE '%V2%' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'Clash' AND model LIKE '%V3%' AND release_year = 0;

-- Ultra
UPDATE racket_library SET release_year = 2018
WHERE brand = 'Wilson' AND series = 'Ultra' AND model LIKE '%V3%' AND release_year = 0;

UPDATE racket_library SET release_year = 2020
WHERE brand = 'Wilson' AND series = 'Ultra' AND model LIKE '%V4%' AND release_year = 0;

UPDATE racket_library SET release_year = 2023
WHERE brand = 'Wilson' AND series = 'Ultra'
  AND (model LIKE '%V5%' OR model = 'Ultra Power 100') AND release_year = 0;

-- Burn（注意 V5 存在 "V5" 与 "V 5" 两种写法）
UPDATE racket_library SET release_year = 2023
WHERE brand = 'Wilson' AND series = 'Burn'
  AND (model LIKE '%V5%' OR model LIKE '%V 5%') AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Wilson' AND series = 'Burn' AND model LIKE '%V6%' AND release_year = 0;

-- Pro Staff
UPDATE racket_library SET release_year = 2023
WHERE brand = 'Wilson' AND series = 'Pro' AND model LIKE '%V14%' AND release_year = 0;

UPDATE racket_library SET release_year = 2026
WHERE brand = 'Wilson' AND series = 'Pro' AND model LIKE '%Classic%' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'Pro'
  AND (model LIKE '%Precision%' OR model LIKE '%RXT%' OR model LIKE '%Legend%') AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'Pro' AND model LIKE '%Pro Open%' AND release_year = 0;

-- RF
UPDATE racket_library SET release_year = 2024
WHERE brand = 'Wilson' AND series = 'RF' AND release_year = 0 AND model NOT LIKE '%Future Lite%';

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'RF' AND model LIKE '%Future Lite%' AND release_year = 0;

-- Shift
UPDATE racket_library SET release_year = 2023
WHERE brand = 'Wilson' AND series = 'Shift' AND model LIKE '%V1%' AND release_year = 0
  AND model NOT LIKE '%2025%' AND model NOT LIKE '%US Open%' AND model NOT LIKE '%Session Soire%';

UPDATE racket_library SET release_year = 2024
WHERE brand = 'Wilson' AND series = 'Shift'
  AND (model LIKE '%US Open%' OR model LIKE '%Session Soire%') AND release_year = 0;

UPDATE racket_library SET release_year = 2023
WHERE model = 'Noir Shift 99 V1' AND release_year = 0;

-- Defyer / Tour
UPDATE racket_library SET release_year = 2026
WHERE brand = 'Wilson' AND series = 'Defyer' AND release_year = 0;

UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND series = 'Tour' AND release_year = 0;

-- 平价/休闲系列（Wilson 当前产品线，年份为估算值，可按需调整）
UPDATE racket_library SET release_year = 2025
WHERE brand = 'Wilson' AND release_year = 0
  AND series IN ('Envy', 'Allure', 'Hyper', 'Triad', 'Six', 'Roland', 'US');

-- string_pattern 统一为 16x19 风格（/、×、X、空格 统一为小写 x）
UPDATE racket_library
SET string_pattern = REPLACE(REPLACE(REPLACE(REPLACE(string_pattern, '/', 'x'), '×', 'x'), 'X', 'x'), ' ', '')
WHERE string_pattern IS NOT NULL AND string_pattern <> '';

-- 复核：剩余 release_year = 0 的行（预期仅 Blade 100L V9 配色行）
SELECT id, brand, series, model, release_year
FROM racket_library
WHERE release_year = 0;
