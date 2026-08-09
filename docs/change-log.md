# Change Log

## 2026-08-10 stats 品牌分组改为按 brand_id + JOIN shoe_brands 取规范品牌名

### 需求/变更内容

- `GET /api/shoe-library/stats` 的品牌/系列统计改为只按 `brand_id` 分组，品牌显示名 LEFT JOIN `shoe_brands` 取主表规范写法，避免 `shoe_library.brand` 历史大小写不一致导致同一品牌被拆成多条。

### 修改文件

- `internal/repository/shoe_repository.go`（`LibraryBrandStats` / `LibrarySeriesStats` 增加 `LEFT JOIN shoe_brands`，`GROUP BY shoe_library.brand_id`；品牌名用 `COALESCE(MAX(shoe_brands.name), MAX(shoe_library.brand))` 兼容 `only_full_group_by`）
- `docs/change-log.md`

### 接口变化

- 响应结构不变；品牌名统一来自 `shoe_brands.name`（如 `ASICS`）。

### 数据库变化

- 无（只读查询变更）。

### 兼容性说明

- `shoe_brands` 中不存在的孤儿 brand_id 回退使用 `shoe_library.brand`。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-10 球鞋品牌名数据规范化（shoe_brands / shoe_library.brand）

### 需求/变更内容

- 修正品牌名大小写不一致：`shoe_brands` 的 `Asics` 改为官方全大写 `ASICS`；K-Swiss 合并重复品牌行（删除未使用的 id=8，规范 id=37 为 `K-Swiss`）。
- `shoe_library.brand` 与主表对齐：`ASICS`、`HEAD`、`K-Swiss`。

### 修改文件

- `migrations/018_normalize_shoe_brand_names.sql`（新增，需在本地与云端数据库各执行一次）
- `docs/change-log.md`

### 接口变化

- 无。

### 数据库变化

- `shoe_brands`：id=3 名称 `Asics -> ASICS`；删除未使用的 id=8；id=37 名称 `KSwiss -> K-Swiss`、slug `k-swiss`。
- `shoe_library`：brand_id=3/17/37 的冗余品牌名修正为 `ASICS` / `HEAD` / `K-Swiss`；引用 brand_id=8 的行归并到 37。

### 兼容性说明

- 品牌 ID 保持稳定（除未使用的 8 删除外），不影响现有球鞋与系列关联。

### 已执行检查命令

- 本地数据库执行迁移并验证分组结果。

### 测试结果

- 通过。

## 2026-08-09 取消球鞋手动输入：POST /api/shoes 强制要求 libraryId

### 需求/变更内容

- 新增球鞋不再支持手动输入（前端移除选鞋页「手动输入」入口，后端创建接口强制从球鞋库选择）。

### 修改文件

- `internal/model/shoe.go`（`CreateShoeRequest.LibraryID` 增加 `binding:"required"`）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- `POST /api/shoes`：`libraryId` 由可选变为必填，缺失或传 `0` 返回 `40001 invalid request`。
- `PUT /api/shoes/:id` 不变：已存在的球鞋（含历史手动输入的球鞋）仍可编辑。

### 数据库变化

- 无。

### 兼容性说明

- 历史手动创建的球鞋数据不受影响，仍可查看与编辑。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-09 管理员新增鞋款：系列改为按名称查找/自动插入

### 需求/变更内容

- 新增鞋款时前端不再单独调用 `POST /api/admin/shoe-series`，也不传系列 ID；只传 `seriesName`，后端按 品牌+性别+系列名 查找，不存在则自动插入系列后写入鞋款。

### 修改文件

- `internal/model/shoe.go`（`CreateShoeLibraryRequest`：`SeriesID` 改为 `SeriesName`）
- `internal/service/shoe_service.go`（`CreateLibraryItems`：按系列名查找或创建，删除系列 ID 校验）
- `internal/handler/shoe_handler.go`（日志字段改为 `seriesName`）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- `POST /api/admin/shoe-library` 请求体由 `seriesId` 改为 `seriesName`；前端不再调用 `POST /api/admin/shoe-series`（该接口保留，仍可用于单独维护系列）。

### 数据库变化

- 无表结构变化；`shoe_series` 仍按 `(brand_id, gender, name)` 唯一键保证不重复。

### 兼容性说明

- 接口为本轮新增，尚无其他调用方；前端同步修改。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-07 管理员新增鞋款：配色多选合并为单条配色

### 需求/变更内容

- 管理员新增鞋款时，多选主流颜色不再每个配色插入一行，而是合并为一条组合配色（如 `Red/Black`，颜色首字母大写）后插入单条 `shoe_library` 记录。

### 修改文件

- `internal/model/shoe.go`（`CreateShoeLibraryRequest.Colorways []string` 改为 `Colorway string`）
- `internal/service/shoe_service.go`（`CreateLibraryItems` 改为单配色写入；移除 `normalizeColorways`）
- `internal/service/shoe_admin_test.go`（移除配色数组规范化测试）
- `internal/handler/shoe_handler.go`（日志字段改为 `colorway`）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- `POST /api/admin/shoe-library` 请求体由 `colorways: string[]` 改为 `colorway: string`，每次请求只插入一条。

### 数据库变化

- 无表结构变化。

### 兼容性说明

- 接口为本轮新增，尚无其他调用方；按新协议提交即可。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-07 球鞋库图片迁移到微信云托管对象存储

### 需求/变更内容

- 新增定时任务：将 `shoe_library.image_url` 的旧外链图片下载并上传到微信云托管对象存储，逻辑与球拍库图片迁移任务一致。
- 上传成功后仅更新 `shoe_library.file_id`（不动 `image_url`，保留旧链接兜底/回退），并写入上传记录表 `shoe_library_image_upload`，保证幂等与审计。
- 存储路径：`{folder}/{brand}/{series}/{id}{ext}`，默认 `shoe_library/{brand}/{series}/{id}.jpg`；brand/series 会清洗为 cloudPath 允许的字符。
- 下载采用内存方式，不落临时文件；同一时间只允许一轮迁移运行（防重入）；单条失败保留原样，下个周期自动重试。
- 与球拍库共用同一 `storage.enabled`/`envId`/`bucket`/`region` 配置与 gocron 调度装配；球鞋库有独立的 `shoeFolder`/`shoeSchedule`/`shoeRunOnStart` 配置项。

### 修改文件

- `config.yaml`（storage 段新增 `shoeFolder`/`shoeSchedule`/`shoeRunOnStart`）
- `internal/config/config.go`（加载球鞋库 storage 配置，支持 `SHOE_STORAGE_FOLDER`/`SHOE_STORAGE_SCHEDULE`/`SHOE_STORAGE_RUN_ON_START` 环境变量覆盖，默认 `shoe_library`/`0 3 * * *`/false）
- `internal/model/shoe.go`（新增 `ShoeLibraryImageUpload` 模型）
- `migrations/018_create_shoe_library_image_upload.sql`（新增上传记录表）
- `internal/repository/shoe_repository.go`（待迁移查询 + 事务更新 file_id 与记录表）
- `internal/service/image_migration.go`（新增，图片迁移共用下载/结果类型）
- `internal/service/racket_image_migration.go`（改为复用共用下载与结果类型，行为不变）
- `internal/service/shoe_image_migration.go`（新增，球鞋库迁移服务）
- `cmd/api/main.go`（gocron 定时任务装配，球拍/球鞋共用调度器公共函数）
- `docs/change-log.md`

### 接口变化

- 无 HTTP 接口变化。

### 数据库变化

- 新增表 `shoe_library_image_upload`（`shoe_library_id` 唯一、`file_id`、`object_key`、`source_url`、`created_at`）。

### 兼容性说明

- `storage.enabled` 默认 `false`，不开启不影响现有服务；球鞋库配置缺省时使用默认目录 `shoe_library` 与默认 cron `0 3 * * *`。
- `image_url` 不被修改，前端在 `fileId` 为空时才回退 `imageUrl`，兼容迁移前/迁移中状态。
- 迁移依赖云托管「开放接口服务」已开启且服务版本已重建；对象存储读权限需允许所有用户读取。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/config/config.go internal/model/shoe.go internal/repository/shoe_repository.go internal/service/image_migration.go internal/service/racket_image_migration.go internal/service/shoe_image_migration.go`
- `go test ./...`

### 测试结果

- 通过（`go test ./...` 全部包通过）。

## 2026-08-07 管理员维护球鞋库接口与小程序入口

### 需求/变更内容

- 新增“管理员新增球鞋”能力：管理员在小程序“选鞋”页品牌下拉里看到“新增鞋款”入口，可新增品牌、系列和鞋款（配色多选，一个配色一行入库）。
- 管理员身份采用配置白名单（`config.yaml` 的 `admin.userIds`），不新增数据库字段；新增 `RequireAdmin` 中间件保护管理员接口。

### 修改文件

- `config.yaml`（新增 `admin.userIds` 管理员白名单）
- `internal/config/config.go`（`Config` / `fileConfig` 增加 `AdminUserIDs`，`Load()` 装配）
- `internal/service/admin.go`（新增 `IsAdminUser`）
- `internal/middleware/admin_middleware.go`（新增 `RequireAdmin`）
- `internal/model/shoe.go`（新增 `CreateShoeBrandRequest` / `CreateShoeSeriesRequest` / `CreateShoeLibraryRequest` / `CreateShoeLibraryResponse`）
- `internal/repository/shoe_repository.go`（新增品牌/系列查询与创建、库条目查重与批量事务写入）
- `internal/service/shoe_service.go`（新增 `CreateBrand` / `CreateSeries` / `CreateLibraryItems` 与 slug、配色规范化纯函数）
- `internal/service/shoe_admin_test.go`（新增 slug、配色去重、管理员判断测试）
- `internal/handler/shoe_handler.go`（新增管理员三个接口 handler）
- `internal/handler/auth_handler.go`（新增 `AdminPermissions`，构造函数增加管理员白名单参数）
- `cmd/api/main.go`（注册 `/api/admin/permissions` 与 `/api/admin/*` 路由）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- 新增 `GET /api/admin/permissions`：返回 `{isAdmin}`，前端控制管理员入口显隐。
- 新增 `POST /api/admin/shoe-brands`：管理员新增球鞋品牌，`name` 必填且唯一，`slug` 可选自动生成。
- 新增 `POST /api/admin/shoe-series`：管理员新增球鞋系列，同品牌同性别下系列名唯一。
- 新增 `POST /api/admin/shoe-library`：管理员新增鞋款，`colorways` 配色数组多选，一个配色一行，逐配色查重后批量事务写入。

### 数据库变化

- 无表结构变化，仅写入 `shoe_brands` / `shoe_series` / `shoe_library` 数据。

### 兼容性说明

- 全部为新增接口，既有接口无改动；`GET /api/auth/me` 等响应结构保持不变。
- 管理员白名单为空时所有用户都是普通用户，管理员接口一律 `403`。
- 默认白名单为 `[1]`，实际账号 ID 不同时改 `config.yaml` 即可。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-07 调低球鞋自然老化系数（restWearPerDay 0.12 -> 0.06）

### 需求/变更内容

- 球鞋磨损度估算中自然老化过快：原默认 `restWearPerDay = 0.12` 意味着仅自然老化约 500 天（约 16.5 个月）即耗尽标准寿命。
- 调整为 `0.06`：仅自然老化约 1000 天（约 2.7 年）耗尽标准寿命，上场小时数仍是磨损主导因素。

### 修改文件

- `config.yaml`（`shoeWear.restWearPerDay` 0.12 -> 0.06）
- `internal/config/shoe_wear.go`（`applyShoeWearDefaults` 默认值同步 0.12 -> 0.06）
- `internal/config/shoe_wear_test.go`（默认值断言同步更新）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- 无。响应结构不变，`wear.score` / `wear.remainingHours` 数值随参数变化。

### 数据库变化

- 无。

### 兼容性说明

- 纯参数调整，同一双球鞋的磨损度得分会整体升高（更晚报废），前端无需改动；如需继续调整，改 `config.yaml` 的 `shoeWear.restWearPerDay` 即可。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-06 球鞋使用统计与磨损度估算

### 需求/变更内容

- 参考球拍使用统计，为球鞋接口补充使用数据：按 `shoeId` 聚合当前用户未删除打球记录的出场次数与时长。
- 参考球线健康度估算，新增球鞋磨损度估算（`wear`）：以购买日期和上场小时数为输入，按 `effectiveWear = 上场小时数 + 购买后天数 × restWearPerDay` 计算得分和剩余寿命，阈值、标准寿命、静置老化系数通过 `config.yaml` 的 `shoeWear` 段配置。

### 修改文件

- `internal/config/shoe_wear.go`（新增 `ShoeWearConfig` / `ShoeWearResolver`，默认 `standardLifeHours = 60`、`restWearPerDay = 0.12`、6 档状态）
- `internal/config/config.go`（`Config` / `fileConfig` 增加 `ShoeWear`，`Load()` 校验并装配）
- `internal/config/shoe_wear_test.go`（新增默认值、校验、config.yaml 加载测试）
- `config.yaml`（新增 `shoeWear` 配置段）
- `internal/model/shoe.go`（`Shoe` / `ShoeResponse` 增加 `UsageCount`、`UsageMinutes`、`UsageHours`、`TotalMinutes`、`TotalHours`、`Wear`；新增 `ShoeUsageStats`、`ShoeWearResponse`）
- `internal/repository/shoe_repository.go`（新增 `UsageStats`，按 `shoe_id` 聚合 `tennis_sessions`）
- `internal/service/shoe_service.go`（`enrichShoes` 补全使用统计与磨损度；新增 `calculateShoeWear` / `calculateShoeWearAt`；构造函数增加 `ShoeWearResolver` 参数）
- `internal/service/shoe_service_test.go`（新增磨损度估算测试）
- `cmd/api/main.go`（装配 `cfg.ShoeWear`）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- `GET /api/shoes`、`GET /api/shoes/selectable`、`GET /api/my-shoes`、`GET /api/my-shoes/primary`、`GET /api/shoes/:id` 及创建/编辑/设主力/退役接口返回的 `ShoeResponse` 新增字段：
  - `usageCount` / `usageMinutes` / `usageHours`：关联打球记录的上场次数与时长
  - `totalMinutes` / `totalHours`：累计使用时长（当前口径与 `usageMinutes` / `usageHours` 一致）
  - `wear`：磨损度估算对象 `{ state, display, score, remainingHours }`
- 均为新增可选字段，不改变既有字段语义。

### 数据库变化

- 无表结构变化。使用统计只读 `tennis_sessions`（`shoe_id`、`duration_minutes`、`deleted_at`），磨损度只读 `my_shoes.purchase_date`。

### 兼容性说明

- 响应为纯新增字段，旧客户端忽略新字段即可，无需同步升级。
- 磨损度默认参数（60 小时标准寿命、每日 0.12 静置老化）为经验默认值，后续可按需在 `config.yaml` 的 `shoeWear` 段调整。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-06 球鞋库 stats 系列统计改为按性别分组的哈希返回

### 需求/变更内容

- `GET /api/shoe-library/stats` 品牌总数保持不区分性别（先选品牌、再选性别的交互）。
- 品牌下的 `series` 改为按性别分组的哈希：key 为性别字符串（`"1"` 男 / `"2"` 女 / `"0"` 未知 / `"3"` 童），值为该性别下的系列统计数组。

### 修改文件

- `internal/model/shoe.go`（`ShoeLibrarySeriesStats` 增加 `Gender`；`ShoeLibraryBrandStatsResponse.Series` 改为 `map[string][]ShoeLibrarySeriesStatsResponse`）
- `internal/repository/shoe_repository.go`（`LibrarySeriesStats` 的 `SELECT` / `GROUP BY` 增加 `gender`）
- `internal/service/shoe_service.go`（合并响应时按 `gender` 分组到哈希）
- `docs/api.md`、`docs/change-log.md`

### 接口变化

- `GET /api/shoe-library/stats` 响应中 `series` 由数组改为按性别分组的哈希（key 为 `"1"` / `"2"` 等性别字符串）。
- 品牌统计 `count` 口径不变（仍为品牌全部行数，不区分性别）。

### 数据库变化

- 无（查询变更，无表结构变化）。

### 兼容性说明

- `series` 结构由数组变为哈希，属于响应结构调整，需前后端同步升级；旧客户端按数组解析会取不到数据。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-05 移除球鞋库发布状态（shoe_library.status）

### 需求/变更内容

- 取消 `shoe_library` 的发布状态概念（`draft / published / archived`）：库内数据入库后默认全部可被 C 端读取，不再需要发布流程。
- 删除数据库 `shoe_library.status` 字段。

### 修改文件

- `internal/model/shoe.go`（`ShoeLibrary` 移除 `Status` 字段）
- `internal/repository/shoe_repository.go`（`LibraryList` / `LibraryBrandStats` / `LibrarySeriesStats` / `LibraryFindByID` / `LibraryFindByIDs` 移除 `status = 'published'` 过滤）
- `internal/service/shoe_service.go`（`applyLibraryDefaults` 注释更新）
- `migrations/017_drop_shoe_library_status.sql`（新增，`ALTER TABLE shoe_library DROP COLUMN status`）
- `migrations/init.sql`（全新环境建表语句同步移除 `shoe_library.status` 列）
- `docs/api.md`、`docs/shoe-add-api-design.md`、`docs/change-log.md`

### 接口变化

- `GET /api/shoe-library`：不再只返回 `published`，返回球鞋库全部数据。
- `GET /api/shoe-library/stats`：不再只统计 `published`，统计全部库数据。
- `POST /api/shoes`：`libraryId` 只需对应库表存在的记录，不再要求 `published`。
- API JSON 响应不变（`status` 本就是内部字段，未返回给前端）。

### 数据库变化

- `shoe_library` 删除 `status` 列（执行 `migrations/017_drop_shoe_library_status.sql`）。
- 已存在的 78 条数据不受影响，删除列后默认全部可读。

### 兼容性说明

- 前端无感知：库列表/统计/创建接口的请求与响应结构不变，只是过滤条件放宽。
- `my_shoes.status`（主力鞋/在用/退役）与本次改动无关，保持不变。

### 已执行检查命令

- `gofmt -w` 相关 Go 文件
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-05 我的球鞋表定名 my_shoes；修复 015 迁移并同步建表

### 需求/变更内容

- "我的球鞋"表名由草案的 `shoe` 定名为 `my_shoes`（用户确认），与前端 `/my-shoes`、`/shoes` 接口对齐。
- 修复 `migrations/015_create_shoe_tables.sql` 中 `my_shoes`（原 `shoe`）表缺少 `PRIMARY KEY (id)`，导致 MySQL ERROR 1075 无法建表的问题。
- `migrations/init.sql` 同步补齐 `my_shoes` / `shoe_brands` / `shoe_series` / `shoe_library` 四张表，避免全新环境缺表。
- 本地 `tennis_diary` 数据库已执行修复后的迁移，创建 `my_shoes` 表。

### 修改文件

- `migrations/015_create_shoe_tables.sql`
- `migrations/init.sql`
- `internal/model/shoe.go`
- `AGENTS.md`
- `docs/api.md`
- `docs/shoe-add-api-design.md`
- `docs/change-log.md`

### 接口变化

- 无。API 路径（`/api/my-shoes`、`/api/shoes` 等）与 JSON 字段（`shoeId`、`shoeName`、`shoeCost`）不变。

### 数据库变化

- 新增 `my_shoes` 表（原草案表名 `shoe` 未在任何库中实际创建，无需数据迁移）。
- 约束/索引命名随表名调整：`my_shoes_status_check`、`idx_my_shoes_user_deleted_status`、`idx_my_shoes_user_deleted_created`、`idx_my_shoes_user_library`。
- 执行方式：`mysql tennis_diary < migrations/015_create_shoe_tables.sql`（幂等，`CREATE TABLE IF NOT EXISTS`）。

### 兼容性说明

- 表名仅影响后端 GORM 映射与 DDL；前端不直接接触数据库表名，无需改动。

### 已执行检查命令

- `go build ./...`
- `go test ./...`
- 建表后执行原报错 SQL（`SELECT COALESCE(SUM(purchase_price), 0) FROM my_shoes ...`）验证通过。

### 测试结果

- 通过。

## 2026-08-05 打球记录支持 shoeId 关联；球鞋购买费用计入首页/统计 expense

### 需求/变更内容

- 打球记录新增/编辑支持传 `shoeId` 关联我的球鞋；响应 `shoeName` 按 `shoeId` 实时返回，`shoeId = 0` 时回退快照 `shoe_name`（保留快照列，兼容旧客户端）。
- 首页 `expense`、日历汇总、统计图表 `summary` 新增 `shoeCost`（按月份/范围内 `shoe.purchase_price` 合计），`totalCost` 口径变为 `sessionCost + racketCost + stringingCost + shoeCost`。
- 统计图表与日历 `expenseBreakdown` 新增"球鞋"分项。

### 修改文件

- `migrations/016_add_session_shoe_id.sql`（新增，`tennis_sessions` 增加 `shoe_id` 列）
- `migrations/init.sql`（建表语句同步增加 `shoe_id` 列）
- `internal/model/session.go`（`TennisSession` / `SessionResponse` / 请求 DTO 增加 `shoeId`；`NewSessionResponse` 增加 `shoeName` 参数）
- `internal/model/home.go`（`HomeExpenseSummaryResponse` 增加 `shoeCost`）
- `internal/model/stats.go`（`StatsChartsSummaryResponse` 增加 `shoeCost`）
- `internal/repository/shoe_repository.go`（新增 `NamesByIDs`、`SumPurchaseCostByMonth` / `SumPurchaseCostByRange`）
- `internal/service/session_service.go`（注入 `shoeRepo`；写入/解析 `shoeId`；日历汇总与 expenseBreakdown 增加球鞋费用）
- `internal/service/home_service.go`（首页 expense 增加 `shoeCost`；最近一次打球关联球鞋名称）
- `internal/service/stats_service.go`（注入 `shoeRepo`；summary 与 expenseBreakdown 增加球鞋费用）
- `internal/service/session_service_test.go`（构造参数同步）
- `cmd/api/main.go`（依赖装配同步）
- `AGENTS.md`、`docs/api.md`、`docs/shoe-add-api-design.md`
- `docs/change-log.md`

### 接口变化

- `POST /api/sessions`、`PUT /api/sessions/:id` 请求体新增可选字段 `shoeId`（我的球鞋 ID，`0` 表示不关联）。
- `SessionResponse` 新增 `shoeId`；`shoeName` 语义变为优先按 `shoeId` 实时关联，无关联时回退快照。
- 首页 `expense` 新增 `shoeCost`，`totalCost` 口径更新。
- 日历 `summary`、`/api/stats/charts` 的 `summary` 新增 `shoeCost`，`totalCost` 口径更新；`expenseBreakdown` 新增"球鞋"分项。

### 数据库变化

- `tennis_sessions` 新增 `shoe_id BIGINT NOT NULL DEFAULT 0`（执行 `migrations/016_add_session_shoe_id.sql`）。

### 兼容性说明

- `shoeName` 请求字段与快照列保留：旧客户端继续传 `shoeName` 可正常保存与展示；新客户端传 `shoeId` 时展示以实时关联为准。
- `totalCost` 口径变化会使历史统计数据变大（新增球鞋购买费用），属预期行为。

### 已执行检查命令

- `gofmt -w internal/model/session.go internal/model/home.go internal/model/stats.go internal/repository/shoe_repository.go internal/service/session_service.go internal/service/home_service.go internal/service/stats_service.go internal/service/session_service_test.go cmd/api/main.go`
- `go build ./...`
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-05 新增球鞋管理功能（球鞋库 + 我的球鞋）

### 需求/变更内容

- 新增球鞋功能，整体流程与球拍一致：球鞋库（品牌/系列/性别）选择 + 手动输入，我的球鞋列表/详情/编辑/设置主力/退役/删除。
- 球鞋库只返回 `published` 数据；同一型号不同配色各占一行，`product_code` 仅内部回源用。
- 打球记录仍使用 `shoeName` 自由文本，本次不接入 `shoeId` 关联；球鞋费用暂不计入首页/统计 `expense`（见设计文档开放问题）。

### 修改文件

- `migrations/015_create_shoe_tables.sql`（新增，`shoe_brands` / `shoe_series` / `shoe_library` / `shoe` 四张表）
- `internal/model/shoe.go`（新增，模型与请求/响应 DTO）
- `internal/repository/shoe_repository.go`（新增）
- `internal/service/shoe_service.go`（新增）
- `internal/handler/shoe_handler.go`（新增）
- `cmd/api/main.go`（装配 `ShoeService` / `ShoeHandler`，注册路由）
- `docs/shoe-add-api-design.md`
- `docs/api.md`（新增 23. 球鞋管理接口）
- `docs/change-log.md`

### 接口变化

新增接口（均需鉴权）：

- `GET /api/shoe-brands`
- `GET /api/shoe-series`（支持 `brandId`、`gender`）
- `GET /api/shoe-library`（支持 `brandId`、`seriesId`、`gender`）
- `GET /api/shoe-library/stats`
- `GET /api/shoes`、`POST /api/shoes`
- `GET /api/shoes/:id`、`PUT /api/shoes/:id`、`DELETE /api/shoes/:id`
- `GET /api/shoes/selectable`、`GET /api/shoes/stats`
- `GET /api/my-shoes`、`GET /api/my-shoes/primary`
- `POST /api/shoes/:id/set-primary`、`POST /api/shoes/:id/retire`

### 数据库变化

- 新增 `shoe_brands`、`shoe_series`、`shoe_library`、`shoe` 四张表（执行 `migrations/015_create_shoe_tables.sql`）。

### 兼容性说明

- 全部为新增接口，不影响现有接口与字段。
- `tennis_sessions.shoe_name` 保持不变，打球记录相关接口无变化。

### 已执行检查命令

- `gofmt -w internal/model/shoe.go internal/repository/shoe_repository.go internal/service/shoe_service.go internal/handler/shoe_handler.go cmd/api/main.go`
- `go build ./...`
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-05 打球记录取消球拍名称快照字段，展示改为关联我的球拍

### 需求/变更内容

- 取消打球记录自带的球拍名称快照字段 `racket_name`，新增/编辑打球记录不再接收 `racketName`。
- 展示名称改为按 `racketId` 关联我的球拍（`racket.name`）实时返回：读取记录时通过 `racket` 表按 ID 关联当前球拍名称，不再依赖历史快照。
- 球拍已被删除的记录，响应 `racketName` 返回空字符串。

### 修改文件

- `internal/model/session.go`（删除 `TennisSession.RacketName`、请求 DTO 的 `RacketName`；`SessionResponse` 保留 `racketName`，由构造函数入参填充）
- `internal/repository/racket_repository.go`（新增 `NamesByIDs`，批量查询当前用户未删除球拍的名称）
- `internal/service/session_service.go`（列表/详情/新增/更新/最近一次读取时关联球拍名称；请求处理移除 `racketName`）
- `internal/service/home_service.go`（首页最近一次打球同样关联球拍名称）
- `internal/service/session_service_test.go`（请求体移除 `racketName`）
- `migrations/014_drop_session_racket_name.sql`（新增，删除 `tennis_sessions.racket_name` 列）
- `migrations/init.sql`（建表语句移除 `racket_name` 列）
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `POST /api/sessions`、`PUT /api/sessions/:id` 请求体不再接受 `racketName`。
- `SessionResponse` 仍返回 `racketName`，语义变为按 `racketId` 关联我的球拍当前名称。

### 数据库变化

- `tennis_sessions` 删除 `racket_name` 列（需要执行 `migrations/014_drop_session_racket_name.sql`）。

### 兼容性说明

- 前端仍可继续使用响应中的 `racketName` 做展示，展示内容变为球拍当前名称而非创建时的名称快照。
- 旧客户端/旧请求若继续传 `racketName`，后端忽略该字段，不报错。

### 已执行检查命令

- `gofmt -w internal/model/session.go internal/repository/racket_repository.go internal/service/session_service.go internal/service/home_service.go internal/service/session_service_test.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-05 项目 logger 支持可配置日志级别

### 需求/变更内容

- 项目自带 logger（`internal/logger`）新增级别过滤，支持 `debug / info / warn / error` 四档；低于当前级别的日志不输出。
- 日志级别可通过 `config.yaml` 的 `logger.level` 配置，也支持环境变量 `LOG_LEVEL` 覆盖（优先级高于配置文件）。
- 未配置时默认 `debug`，保持原有行为不变。
- 将错误路径日志从 `logger.Debug` 提升为 `logger.Warn` / `logger.Error`，避免生产环境关闭 debug 后丢失关键错误信息（球拍图片迁移失败、service 内部错误、认证内部错误）。

### 修改文件

- `internal/logger/logger.go`（新增 `Level` / `SetLevel` / `Enabled` / `Info` / `Warn` / `Error`，按级别过滤输出）
- `internal/logger/logger_test.go`（新增，级别解析与过滤单元测试）
- `internal/config/config.go`（加载 `logger.level`，支持 `LOG_LEVEL` 环境变量覆盖）
- `config.yaml`（新增 `logger` 段，默认 `debug`）
- `cmd/api/main.go`（图片迁移运行错误改为 `Error` 级别）
- `internal/service/racket_image_migration.go`（单条迁移失败改为 `Warn` 级别）
- `internal/handler/session_handler.go`（service 内部错误改为 `Error` 级别）
- `internal/handler/auth_handler.go`（认证内部错误补充 `Error` 级别日志）
- `docs/change-log.md`

### 接口变化

- 无 HTTP 接口变化。

### 数据库变化

- 无。

### 兼容性说明

- 默认级别仍为 `debug`，不配置时日志行为与之前一致。
- 生产环境建议设置 `LOG_LEVEL=info`（或 `warn`），可配合 `GIN_MODE=release` 减少线上日志量。
- `LOG_LEVEL` 非法值会导致服务启动失败（panic），与现有配置校验行为一致。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/config/config.go internal/handler/auth_handler.go internal/handler/session_handler.go internal/logger/logger.go internal/logger/logger_test.go internal/service/racket_image_migration.go`
- `go test ./...`

### 测试结果

- `go test ./...` 全部通过。

## 2026-08-04 球拍库图片迁移到微信云托管对象存储

### 需求/变更内容

- 新增定时任务：将 `racket_library.image_url` 的旧外链图片下载并上传到微信云托管对象存储。
- 上传成功后仅更新 `racket_library.file_id`（不动 `image_url`，保留旧链接兜底/回退），并写入上传记录表 `racket_library_image_upload`，保证幂等与审计。
- 上传使用微信云托管官方服务端路径：开放接口服务 `/_/cos/getauth` 获取临时密钥 + `/_/cos/metaid/encode` 生成文件元数据（管理端 openid 为空），写入 `x-cos-meta-fileid` 头后经 COS SDK `putObject` 上传。
- 存储路径：`{folder}/{brand}/{series}/{id}{ext}`，默认 `racket_library/{brand}/{series}/{id}.jpg`；brand/series 会清洗为 cloudPath 允许的字符（中文保留，空格等特殊字符转 `-`）。
- 下载采用内存方式，不落临时文件，因此无残留文件需要删除。
- 同一时间只允许一轮迁移运行（防重入）；单条失败保留原样，下个周期自动重试。
- 前端展示改为优先读 `fileId`、`imageUrl` 兜底（由前端仓库另行修改）。

### 修改文件

- `config.yaml`（新增 `storage` 段：`enabled`/`envId`/`bucket`/`region`/`folder`/`schedule`/`runOnStart`）
- `internal/config/config.go`（加载 storage 配置，支持环境变量覆盖）
- `internal/model/racket_library.go`（新增 `RacketLibraryImageUpload` 模型）
- `migrations/013_create_racket_library_image_upload.sql`（新增上传记录表）
- `internal/storage/cloud_storage.go`（新增，云存储客户端）
- `internal/storage/cloud_storage_test.go`（新增，单元测试）
- `internal/repository/racket_repository.go`（待迁移查询 + 事务更新 file_id 与记录表）
- `internal/service/racket_image_migration.go`（新增，迁移服务）
- `cmd/api/main.go`（gocron 定时任务装配 + runOnStart）
- `docs/change-log.md`

### 接口变化

- 无 HTTP 接口变化。

### 数据库变化

- 新增表 `racket_library_image_upload`（`racket_library_id` 唯一、`file_id`、`object_key`、`source_url`、`created_at`）。

### 兼容性说明

- `storage.enabled` 默认 `false`，不开启不影响现有服务。
- `image_url` 不被修改，前端在 `fileId` 为空时才回退 `imageUrl`，兼容迁移前/迁移中状态。
- 迁移依赖云托管「开放接口服务」已开启且服务版本已重建；对象存储读权限需允许所有用户读取。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/config/config.go internal/model/racket_library.go internal/repository/racket_repository.go internal/service/racket_image_migration.go internal/storage/cloud_storage.go internal/storage/cloud_storage_test.go`
- `go test ./...`

### 测试结果

- 通过（`go test ./...` 全部包通过，含新增 storage 单元测试）。

## 2026-08-04 聚酯线健康度配置化

### 需求/变更内容

- 将 `stringHealthState` 中的状态分档（key、minScore、label、display、color）从代码硬编码迁移到 `config.yaml`。
- 同步将公式常量 `standardLifeHours`（16）和 `restWearPerDay`（0.18）迁移到配置，业务层通过 resolver 读取。
- display 模板支持 `{remainingHours}` 占位符，由 resolver 渲染为实际剩余小时数。
- 配置段缺失时使用代码内置默认值，配置存在时严格校验（key 唯一、minScore 在 `[0, 100]` 且唯一、必须存在 `minScore: 0` 的兜底状态、display 必填）。

### 修改文件

- `internal/config/string_health.go`（新增）
- `internal/config/string_health_test.go`（新增）
- `internal/config/config.go`
- `config.yaml`
- `internal/service/racket_service.go`
- `internal/service/racket_service_test.go`
- `cmd/api/main.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 无。`stringHealth` 响应结构不变（`state` / `display` / `score` / `remainingHours`），仅阈值和文案来源改为配置。

### 数据库变化

- 无。

### 兼容性说明

- 默认配置与原有硬编码行为完全一致，现有客户端无需改动。
- 修改 `config.yaml` 的 `polyesterStringHealth` 段即可调整分档阈值和展示文案，无需重新编译。
- 配置段缺失时回退到内置默认值，不影响服务启动。

### 已执行检查命令

- `gofmt -w internal/config/config.go internal/config/string_health.go internal/config/string_health_test.go internal/service/racket_service.go internal/service/racket_service_test.go cmd/api/main.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-04 球拍新增聚酯线健康度

### 需求/变更内容

- 依据 `polyester-string-health.md` 实现聚酯线衰减健康度计算。
- `RacketResponse` 新增 `stringHealth` 字段，包含 `state`、`display`、`score`、`remainingHours`。
- 健康度基于最近一次穿线日期和穿线后累计打球分钟数计算，公式与产品文档一致：
  `effective_wear = hours_played + days_since_stringing * 0.18`，标准寿命 `16` 小时。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `internal/service/racket_service_test.go`（新增）
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `RacketResponse` 新增响应字段：`stringHealth`。
- 生效接口：`GET /api/rackets`、`GET /api/rackets/selectable`、`GET /api/my-rackets`、`GET /api/my-rackets/primary`、`GET /api/rackets/:id`、创建/编辑/设为主力/退役接口的返回。
- 存在最近穿线记录时返回健康度对象；无穿线记录时返回 `null`。
- 随健康度一并明确 `afterStringingUsageCount` / `afterStringingUsageMinutes` / `afterStringingUsageHours` 的口径：最近一次穿线之后的累计使用统计。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 当前版本未区分线型，统一按聚酯线公式计算，不区分线径、打法、气温、湿度、击球强度和穿线师差异。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go internal/service/racket_service_test.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-08-03 新增球拍库品牌系列数量统计接口

### 需求/变更内容

- 新增 `GET /api/racket-library/stats`，返回球拍库中每个品牌的球拍数量。
- 同一响应中返回每个品牌下每个系列的球拍数量。
- 统计口径来自 `racket_library`，按数量降序返回。

### 修改文件

- `cmd/api/main.go`
- `internal/model/racket_library.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `internal/handler/racket_handler.go`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- 新增接口：
  - `GET /api/racket-library/stats`

### 数据库变化

- 无。

### 兼容性说明

- 仅新增接口，不修改现有接口字段和行为。
- 接口需要 JWT 鉴权，和现有球拍库接口保持一致。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/model/racket_library.go internal/repository/racket_repository.go internal/service/racket_service.go internal/handler/racket_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-23 球拍列表返回最新穿线记录

### 需求/变更内容

- `GET /api/rackets` 每把球拍新增 `latestStringingRecord`，返回该球拍最近一次完整穿线记录。
- `GET /api/rackets/:id` 的 `data.racket` 同步返回 `latestStringingRecord`，保持 `RacketResponse` 结构一致。
- 保留原有 `stringName`、`verticalTension`、`horizontalTension`、`lastStringDate`、`lastStringCost` 摊平字段，兼容旧客户端。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/rackets` 每个球拍新增响应字段：`latestStringingRecord`。
- `GET /api/rackets/:id` 的 `data.racket` 同步新增响应字段：`latestStringingRecord`。
- 无穿线记录时，`latestStringingRecord` 返回 `null`。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 最新穿线记录复用现有批量查询逻辑，不新增 N+1 查询。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-19 新增球拍品牌和系列接口

### 需求/变更内容

- 初始化脚本新增 `racket_brands` 和 `racket_series` 表，用于维护球拍品牌和系列基础数据。
- 新增 `GET /api/racket-brands` 返回球拍品牌列表。
- 新增 `GET /api/racket-series` 返回球拍系列列表，并支持通过 `brandId` 查询指定品牌下的系列。

### 修改文件

- `cmd/api/main.go`
- `internal/model/racket_library.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `internal/handler/racket_handler.go`
- `migrations/init.sql`
- `migrations/010_create_racket_brand_series.sql`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- 新增接口：
  - `GET /api/racket-brands`
  - `GET /api/racket-series`
  - `GET /api/racket-series?brandId=1`

### 数据库变化

- 新增表 `racket_brands`：
  - `id BIGINT PRIMARY KEY AUTO_INCREMENT`
  - `name VARCHAR(128) NOT NULL`
  - `file_id VARCHAR(255) DEFAULT NULL`
  - `created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`
- 新增表 `racket_series`：
  - `id BIGINT PRIMARY KEY AUTO_INCREMENT`
  - `brand_id BIGINT NOT NULL`
  - `name VARCHAR(128) NOT NULL`
  - `created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP`

### 兼容性说明

- 新增接口和新表，不影响现有球拍库、我的球拍和打球记录接口。
- `racket_series.brand_id` 仅逻辑关联 `racket_brands.id`，未添加数据库外键。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/model/racket_library.go internal/repository/racket_repository.go internal/service/racket_service.go internal/handler/racket_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-19 球拍库支持品牌和系列过滤

### 需求/变更内容

- `GET /api/racket-library` 新增可选 query 参数 `brandId` 和 `seriesId`。
- 支持按品牌 ID、系列 ID 或品牌 ID + 系列 ID 组合过滤球拍库。
- 不传过滤参数时保持原有全量返回和按品牌分组结构不变。

### 修改文件

- `internal/model/racket_library.go`
- `internal/handler/racket_handler.go`
- `internal/service/racket_service.go`
- `internal/repository/racket_repository.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/racket-library` 新增 query 参数：
  - `brandId`：品牌 ID，可选，必须为正整数。
  - `seriesId`：系列 ID，可选，必须为正整数。

### 数据库变化

- 无。

### 兼容性说明

- 不传参数时接口响应与原行为一致。
- 响应仍按品牌分组，过滤后只返回匹配球拍项所在的品牌分组。

### 已执行检查命令

- `gofmt -w internal/model/racket_library.go internal/handler/racket_handler.go internal/service/racket_service.go internal/repository/racket_repository.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-19 球拍列表返回 fileId

### 需求/变更内容

- `GET /api/rackets` 返回增加球拍 `fileId`，从关联的球拍库 `racket_library.file_id` 填充。
- `GET /api/rackets/:id` 详情中的 `data.racket` 同步返回 `fileId`，保持 `RacketResponse` 响应结构一致。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/rackets` 每个球拍新增响应字段：`fileId`。
- `GET /api/rackets/:id` 的 `data.racket` 同步新增响应字段：`fileId`。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 未关联球拍库或球拍库记录不存在时，`fileId` 返回空字符串。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-19 球拍库新增 fileId 字段

### 需求/变更内容

- `racket_library` 新增 `file_id` 字段，用于保存球拍图片文件 ID。
- `GET /api/racket-library` 每个球拍项新增 `fileId` 响应字段。
- 初始化脚本同步补充建表字段和旧表幂等补列逻辑。

### 修改文件

- `internal/model/racket_library.go`
- `migrations/init.sql`
- `migrations/009_add_racket_library_file_id.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/racket-library` 的每个球拍项新增 `fileId`。

### 数据库变化

- `racket_library` 新增字段：
  - `file_id VARCHAR(255) NULL DEFAULT NULL COMMENT '球拍图片文件ID'`

### 兼容性说明

- 仅新增字段，不删除或改名已有字段。
- 旧数据 `fileId` 默认为空字符串响应。

### 已执行检查命令

- `gofmt -w internal/model/racket_library.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-19 球拍列表返回球拍库参数

### 需求/变更内容

- `GET /api/rackets` 返回增加球拍参数，包含关联球拍库的发布年份、重量、拍面大小和穿线模式。
- `GET /api/rackets/:id` 详情中的 `racket` 同步返回相同球拍参数，保持 `RacketResponse` 响应结构一致。
- 批量查询球拍库参数并在 service enrich 阶段补充，避免列表接口 N+1 查询。

### 修改文件

- `internal/model/racket.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/rackets` 每个球拍新增响应字段：
  - `releaseYear`
  - `weight`
  - `headSize`
  - `stringPattern`
- `GET /api/rackets/:id` 的 `data.racket` 同步新增上述字段。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 未关联球拍库或球拍库记录不存在时，数字字段返回 `0`，`stringPattern` 返回空字符串。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/repository/racket_repository.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-18 球拍库新增穿线模式字段

### 需求/变更内容

- `racket_library` 新增 `string_pattern` 字段，用于保存穿线模式，例如 `16x19`、`18x20`、`16/19`。
- 确认 `head_size` 字段已存在，用于保存拍面大小数值，例如 `98`，展示单位为 `sq.in.`。
- `GET /api/racket-library` 每个球拍项新增 `stringPattern` 响应字段。

### 修改文件

- `internal/model/racket_library.go`
- `migrations/init.sql`
- `migrations/008_add_racket_library_string_pattern.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/racket-library` 的每个球拍项新增 `stringPattern`。

### 数据库变化

- `racket_library` 新增字段：
  - `string_pattern VARCHAR(64) NULL DEFAULT NULL COMMENT '穿线模式，例如 16x19、18x20、16/19'`
- `migrations/init.sql` 同步补充建表和幂等补列逻辑。

### 兼容性说明

- 旧数据 `stringPattern` 默认为空字符串响应。
- `headSize` 保持数值字段，不在后端拼接单位。

### 已执行检查命令

- `gofmt -w internal/model/racket_library.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-18 球拍库新增品牌与系列字段并清空旧数据

### 需求/变更内容

- 为 `racket_library` 新增 `brand_id`、`series_id`、`series` 字段，便于前端按品牌/系列做展示和筛选。
- 本次迁移同步清空 `racket_library` 旧数据，作为球拍库重建的起点。
- 球拍库响应同步返回新增字段，品牌分组顶层也补充 `brandId`。
- 新增初始化脚本补列逻辑和独立迁移，保证新旧库都能平滑升级。

### 修改文件

- `internal/model/racket_library.go`
- `internal/service/racket_service.go`
- `migrations/init.sql`
- `migrations/007_add_racket_library_series_fields.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/racket-library` 的分组项新增 `brandId`。
- `GET /api/racket-library` 的每个球拍项新增 `brandId`、`seriesId`、`series`。
- 其他球拍相关接口响应不变。

### 数据库变化

- `racket_library` 新增字段：
  - `brand_id BIGINT NOT NULL DEFAULT 0`
  - `series_id BIGINT NOT NULL DEFAULT 0`
  - `series VARCHAR(100) NOT NULL DEFAULT ''`
- 迁移执行时先清空 `racket_library` 旧数据，再添加上述字段。
- 新增迁移：`migrations/007_add_racket_library_series_fields.sql`。
- `migrations/init.sql` 同步补充建表和幂等补列逻辑。

### 兼容性说明

- `racket_library` 旧数据会被清空，请先确认不再需要历史球拍库数据。
- 现有品牌分组逻辑仍按 `brand` 分组，不依赖新字段。

### 已执行检查命令

- `gofmt -w internal/model/racket_library.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-15 用户输入文本接入微信内容安全检测

### 需求/变更内容

- 新增服务端微信内容安全检测能力，调用 `wxa/msg_sec_check` 检查用户可输入文本是否违法违规。
- 在用户昵称、打球记录文本字段、球拍文本字段、穿线球线名称保存前执行检测。
- 检测到敏感内容时阻止保存，并返回友好提示：`输入内容包含敏感信息，请修改后重试`。
- 本地未配置 `wechat.appId` / `wechat.appSecret` 时跳过真实检测，避免影响开发联调。

### 修改文件

- `cmd/api/main.go`
- `internal/service/content_security_service.go`
- `internal/service/errors.go`
- `internal/service/auth_service.go`
- `internal/service/session_service.go`
- `internal/service/racket_service.go`
- `internal/handler/auth_handler.go`
- `internal/handler/session_handler.go`
- `docs/change-log.md`

### 接口变化

- 无新增接口。
- 以下接口在保存用户输入文本前会进行内容安全检测：
  - `PUT /api/auth/profile`
  - `POST /api/sessions`
  - `PUT /api/sessions/:id`
  - `POST /api/rackets`
  - `PUT /api/rackets/:id`
  - `POST /api/rackets/:id/stringing-records`

### 数据库变化

- 无。

### 兼容性说明

- API 字段和响应结构不变。
- 需要线上配置有效 `WECHAT_APP_ID` / `WECHAT_APP_SECRET` 或 `config.yaml` 中的微信配置后，才会调用微信内容安全接口。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/service/content_security_service.go internal/service/errors.go internal/service/auth_service.go internal/service/session_service.go internal/service/racket_service.go internal/handler/auth_handler.go internal/handler/session_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 修复 matchRank 校验范围

### 需求/变更内容

- 修复比赛成绩重排后 `matchRank=4` 等有效值提交打球记录返回 `invalid request` 的问题。
- `MatchRank.IsValid()` 校验上限从 `MatchRankThirdPlace` 修正为 `MatchRankGroupStage`。

### 修改文件

- `internal/model/enum.go`
- `docs/change-log.md`

### 接口变化

- `POST /api/sessions` 和 `PUT /api/sessions/:id` 支持提交当前全部有效比赛成绩：`0-7`。

### 数据库变化

- 无。

### 兼容性说明

- 无。

### 已执行检查命令

- `gofmt -w internal/model/enum.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 stats/charts 类型占比改为配置驱动

### 需求/变更内容

- 优化 `GET /api/stats/charts` 的类型占比返回，一级/二级类型改为结合 `sessionEnums` 配置生成。
- 配置新增二级类型后，`sessionSubCategory*Breakdown` 自动返回对应项，无需再改代码硬编码。
- 类型占比项新增 `category` / `subCategory` 字段，方便前端按枚举值识别。
- 类型占比项 `key` 改为稳定数值格式：`category_{category}`、`category_{category}_sub_{subCategory}`。

### 修改文件

- `cmd/api/main.go`
- `internal/model/stats.go`
- `internal/service/stats_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/stats/charts` 的以下字段改为按配置返回：
  - `charts.sessionCategoryCountBreakdown`
  - `charts.sessionSubCategoryCountBreakdown`
  - `charts.sessionCategoryDurationBreakdown`
  - `charts.sessionSubCategoryDurationBreakdown`
  - `charts.sessionCategoryCostBreakdown`
  - `charts.sessionSubCategoryCostBreakdown`
- 上述 breakdown item 新增 `category`，二级类型 item 额外新增 `subCategory`。
- 上述 breakdown item 的 `key` 从旧英文固定值调整为数值组合 key。

### 数据库变化

- 无。

### 兼容性说明

- 前端不应再依赖旧 `daily`、`training_serve` 等硬编码 key，应优先使用 `category` / `subCategory`。
- 数据库中已不在当前配置内的历史类型不会出现在 breakdown 明细项中。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/model/stats.go internal/service/stats_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 session-config 删除 defaultValues 字段

### 需求/变更内容

- 删除 `GET /api/session-config` 响应中的 `defaultValues` 字段。
- 删除不再使用的 `SessionDefaultEnumsResponse` DTO。

### 修改文件

- `internal/model/enums.go`
- `internal/service/enum_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/session-config` 的 `data.defaultValues` 不再返回。

### 数据库变化

- 无。

### 兼容性说明

- 前端如需默认值，应使用本地表单默认逻辑或独立配置，不再从该接口读取。

### 已执行检查命令

- `gofmt -w internal/model/enums.go internal/service/enum_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 session-config 删除 legacyTypes 字段

### 需求/变更内容

- 删除 `GET /api/session-config` 响应中的顶层 `legacyTypes` 字段。
- 保留 `subCategories[].legacyType`，用于现有旧 `type` 字段兼容映射。

### 修改文件

- `internal/model/enums.go`
- `internal/service/enum_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/session-config` 的 `data.legacyTypes` 不再返回。

### 数据库变化

- 无。

### 兼容性说明

- 新客户端应使用 `categories` 和 `subCategories` 构建打球类型选择。
- 旧 `type` 兼容值仍可从每个 `subCategories[].legacyType` 获取。

### 已执行检查命令

- `gofmt -w internal/model/enums.go internal/service/enum_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 matchRanks 新增季军并重排枚举值

### 需求/变更内容

- 比赛成绩 `matchRanks` 新增 `季军`。
- 按“不追加枚举值、不考虑历史数据”的要求，将 `季军` 插入为 `matchRank=3`。
- 后续名次顺延为：`四强=4`、`八强=5`、`16强=6`、`小组赛=7`。

### 修改文件

- `internal/model/enum.go`
- `internal/service/enum_service.go`
- `migrations/init.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/session-config` 的 `matchRanks` 新增 `季军`，并重排后续枚举值。
- `POST /api/sessions` 和 `PUT /api/sessions/:id` 的比赛类型记录支持 `matchRank=7`。

### 数据库变化

- 无表结构变化。
- `migrations/init.sql` 更新 `match_rank` 字段注释。

### 兼容性说明

- 本次按需求不考虑历史数据兼容，既有 `match_rank` 数值会按新枚举含义展示。

### 已执行检查命令

- `gofmt -w internal/model/enum.go internal/service/enum_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 matchRanks 新增 16 强

### 需求/变更内容

- 比赛成绩 `matchRanks` 新增 `16强`。
- `GET /api/session-config` 的 `matchRanks` 返回新增 `{ "value": 5, "label": "16强" }`。
- 原 `小组赛` 顺延为 `matchRank=6`。

### 修改文件

- `internal/model/enum.go`
- `internal/service/enum_service.go`
- `migrations/init.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/session-config` 的 `matchRanks` 新增 `16强`，并将 `小组赛` 调整为 `value=6`。
- `POST /api/sessions` 和 `PUT /api/sessions/:id` 的比赛类型记录支持 `matchRank=6`。

### 数据库变化

- 无表结构变化。
- `migrations/init.sql` 更新 `match_rank` 字段注释。

### 兼容性说明

- 若已有线上数据使用 `match_rank=5` 表示小组赛，需要执行数据迁移将历史小组赛更新为 `6`，否则会按新枚举显示为 `16强`。

### 已执行检查命令

- `gofmt -w internal/model/enum.go internal/service/enum_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-04 打球记录分类枚举改为配置驱动

### 需求/变更内容

- 将打球记录一级/二级分类枚举从纯代码硬编码改为 `config.yaml` 配置驱动。
- `GET /api/enums` 的 `session.categories` 从配置读取，便于新增“训练-截击”等二级类型。
- 新增/编辑打球记录时，`category` / `subCategory` 合法组合改为按配置校验。
- 打球记录响应中的 `categoryText`、`subCategoryText`、`typeText` 优先使用配置文案。
- `config.yaml` 未配置 `sessionEnums.categories` 时启动直接报错，避免线上枚举来源不明确。

### 修改文件

- `cmd/api/main.go`
- `internal/config/config.go`
- `internal/config/session_enums.go`
- `internal/model/session.go`
- `internal/service/enum_service.go`
- `internal/service/home_service.go`
- `internal/service/session_service.go`
- `docs/change-log.md`

### 接口变化

- `GET /api/session-config` 新增公开接口，专用于获取打球记录配置，无需 JWT 鉴权。
- 原 `GET /api/enums` 不再注册，打球记录配置改由 `GET /api/session-config` 返回。
- `GET /api/session-config` 响应 `data` 直接返回 `categories`、`matchRanks`，不再嵌套 `session` / `racket`。

### 数据库变化

- 无。

### 兼容性说明

- `GET /api/session-config` 为公开配置接口，不依赖当前登录用户。
- 未配置 `sessionEnums.categories` 时服务启动失败，需显式维护分类配置。
- 配置新增的训练类二级分类会同步映射到旧 `type=3`，兼容旧字段。
- 比赛类仍按旧规则映射到旧单打/双打比赛类型，用于保留 `matchRank` 规则。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/enum_handler.go internal/model/enums.go internal/service/enum_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-01 优化穿线后使用统计口径

### 需求/变更内容

- 修复球拍穿线后使用次数和使用时间统计只按穿线日期计算的问题。
- 穿线后使用统计改为按最近一次穿线记录的完整 `string_date` 计算，精确到分钟。
- 优化球拍列表/详情 enrich 逻辑，将每把球拍一次穿线后统计查询改为批量聚合查询，避免 N+1 查询。

### 修改文件

- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 响应字段不变。
- `afterStringingUsageCount`、`afterStringingUsageMinutes`、`afterStringingUsageHours` 的统计口径从“穿线当天及之后”调整为“最近一次穿线时间及之后”。

### 数据库变化

- 无。

### 兼容性说明

- JSON 字段名和结构不变。
- 对于同一天穿线前的打球记录，新口径不再计入穿线后使用统计。

### 已执行检查命令

- `gofmt -w internal/repository/racket_repository.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-01 打球记录支持开始时间到分钟

### 需求/变更内容

- 新增打球记录 `date` 支持 `YYYY-MM-DD HH:mm`，用于记录开始时间到分钟。
- 保留旧格式 `YYYY-MM-DD` 兼容，按当天 `00:00` 处理。
- 打球记录响应中的 `date` 统一返回 `YYYY-MM-DD HH:mm`。
- 日历和按日期查询仍按自然日聚合/筛选。
- 修复日历天数聚合 SQL，避免 `SELECT DATE_FORMAT(date, ...)` 与 `GROUP BY DATE(date)` 表达式不一致在 MySQL `ONLY_FULL_GROUP_BY` 下触发 500。

### 修改文件

- `internal/model/session.go`
- `internal/service/session_service.go`
- `internal/repository/session_repository.go`
- `migrations/init.sql`
- `migrations/006_change_session_date_to_datetime.sql`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- `POST /api/sessions` 和 `PUT /api/sessions/:id` 的请求字段 `date` 支持：
  - `YYYY-MM-DD HH:mm`
  - `YYYY-MM-DD`（兼容旧客户端）
- `SessionResponse.date` 返回格式调整为 `YYYY-MM-DD HH:mm`。
- `GET /api/sessions?date=YYYY-MM-DD` 仍按自然日查询当天记录。

### 数据库变化

- `tennis_sessions.date` 修改为 `DATETIME NOT NULL COMMENT '打球开始时间'`。
- 新增迁移：`migrations/006_change_session_date_to_datetime.sql`。

### 兼容性说明

- 历史日期数据迁移为 `DATETIME` 后时间默认为 `00:00:00`。
- 日期范围统计、首页统计和日历接口继续使用自然日/月/年范围。

### 已执行检查命令

- `gofmt -w internal/model/session.go internal/service/session_service.go internal/repository/session_repository.go`
- `go test ./...`
- 检查所有仓储日期聚合 SQL，确认仅日历天数聚合存在同类表达式不一致问题。

### 测试结果

- 通过。

## 2026-07-01 球拍穿线记录支持横竖线磅数

### 需求/变更内容

- 新增球拍穿线记录支持分别填写竖线磅数 `verticalTension` 和横线磅数 `horizontalTension`。
- 接口不再接收或返回旧字段 `tension`。
- 球拍列表、详情和穿线记录响应新增 `verticalTension` / `horizontalTension`。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `migrations/init.sql`
- `migrations/005_add_stringing_split_tension.sql`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- `POST /api/rackets/:id/stringing-records` 新增请求字段：
  - `verticalTension`
  - `horizontalTension`
- `RacketResponse` 和 `StringingRecordResponse` 新增响应字段：
  - `verticalTension`
  - `horizontalTension`

### 数据库变化

- `racket_stringing_record` 新增字段：
  - `vertical_tension DECIMAL(4,1) DEFAULT NULL`
  - `horizontal_tension DECIMAL(4,1) DEFAULT NULL`
- 新增迁移：`migrations/005_add_stringing_split_tension.sql`。

### 兼容性说明

- API 不再兼容旧字段 `tension`；客户端需提交 `verticalTension` / `horizontalTension`。
- 本次迁移只新增字段，不再对历史数据执行回填兜底。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-01 球拍穿线记录支持小时分钟

### 需求/变更内容

- 新增球拍穿线记录时，`stringDate` 支持 `YYYY-MM-DD HH:mm`，可精确到小时分钟。
- 兼容旧格式 `YYYY-MM-DD`，按当天 `00:00` 处理。
- 穿线记录响应中的 `stringDate` 统一返回 `YYYY-MM-DD HH:mm`。
- 数据库字段 `racket_stringing_record.string_date` 从 `DATE` 调整为 `DATETIME`，避免丢失时分。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `migrations/init.sql`
- `migrations/004_change_stringing_date_to_datetime.sql`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- `POST /api/rackets/:id/stringing-records` 的请求字段 `stringDate` 支持：
  - `YYYY-MM-DD HH:mm`
  - `YYYY-MM-DD`（兼容旧客户端）
- `StringingRecordResponse.stringDate` 返回格式调整为 `YYYY-MM-DD HH:mm`。

### 数据库变化

- `racket_stringing_record.string_date` 修改为 `DATETIME NOT NULL COMMENT '穿线时间'`。
- 新增迁移：`migrations/004_change_stringing_date_to_datetime.sql`。

### 兼容性说明

- 旧客户端继续传 `YYYY-MM-DD` 可正常新增穿线记录。
- 历史日期数据迁移为 `DATETIME` 后时间默认为 `00:00:00`。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-27 统计图表新增打球次数和时长类型占比

### 需求/变更内容

- 修改 `GET /api/stats/charts`，新增打球次数按一级类型和二级类型维度的占比统计。
- 修改 `GET /api/stats/charts`，新增打球时长按一级类型和二级类型维度的占比统计。
- 复用类型聚合查询，同时返回次数、分钟数和费用，避免重复 SQL。
- 一级类型固定返回日常球局、训练、比赛；二级类型固定返回 6 个合法类型组合。
- 保留原有费用占比和消费占比字段。
- 按需求取消 `GET /api/stats/charts` 的 `charts.sessionTypeBreakdown` 返回，前端应改用新增的 `sessionCategoryCountBreakdown` / `sessionSubCategoryCountBreakdown`。

### 修改文件

- `internal/model/stats.go`
- `internal/repository/session_repository.go`
- `internal/service/stats_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/stats/charts` 的 `charts` 新增字段：
  - `sessionCategoryCountBreakdown`
  - `sessionSubCategoryCountBreakdown`
  - `sessionCategoryDurationBreakdown`
  - `sessionSubCategoryDurationBreakdown`
- `GET /api/stats/charts` 的 `charts` 不再返回字段：
  - `sessionTypeBreakdown`

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 无记录时新增占比字段仍固定返回类型项，`value` 和 `percent` 为 `0`。

### 已执行检查命令

- `gofmt -w internal/model/stats.go internal/repository/session_repository.go internal/service/stats_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-27 统计图表新增打球费用类型占比

### 需求/变更内容

- 修改 `GET /api/stats/charts`，新增打球费用按一级类型和二级类型维度的占比统计。
- 一级类型按 `category` 聚合，固定返回日常球局、训练、比赛。
- 二级类型按 `category + sub_category` 聚合，固定返回 6 个合法类型组合。
- 占比分母为查询范围内 `tennis_sessions.cost` 合计，即 `summary.sessionCost`。
- 保留原有 `expenseBreakdown` 和 `sessionTypeBreakdown` 字段，避免影响现有前端展示。

### 修改文件

- `internal/model/stats.go`
- `internal/repository/session_repository.go`
- `internal/service/stats_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/stats/charts` 的 `charts` 新增字段：
  - `sessionCategoryCostBreakdown`
  - `sessionSubCategoryCostBreakdown`

### 数据库变化

- 无。

### 兼容性说明

- 仅新增 JSON 字段，不删除或改名已有字段。
- 无打球费用时新增占比字段仍固定返回类型项，`value` 和 `percent` 为 `0`。

### 已执行检查命令

- `gofmt -w internal/model/stats.go internal/repository/session_repository.go internal/service/stats_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-27 打球记录新增两级分类字段

### 需求/变更内容

- 在打球记录中新增 `category` / `subCategory` 两级分类字段。
- 新客户端可提交一级类型和二级类型，后端校验合法组合并同步旧 `type` 字段。
- 旧客户端仍可只提交 `type`，后端自动映射为新两级分类。
- 响应新增 `category`、`categoryText`、`subCategory`、`subCategoryText`、`typeText`，方便列表、详情、首页最近记录展示。
- 新增 `GET /api/enums`，为前端提供打球类型、旧类型、比赛成绩和球拍状态枚举值及展示文案。
- 历史数据通过迁移按旧 `type` 回填新字段，旧训练统一归为 `2 / 4`。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/enum_handler.go`
- `internal/model/enum.go`
- `internal/model/enums.go`
- `internal/model/racket.go`
- `internal/model/session.go`
- `internal/service/enum_service.go`
- `internal/service/session_service.go`
- `internal/repository/session_repository.go`
- `migrations/init.sql`
- `migrations/002_add_session_category.sql`
- `migrations/003_convert_session_category_to_numeric.sql`
- `docs/api.md`
- `docs/session-type-redesign.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- `GET /api/enums` 新增响应，返回前端枚举值：
  - `session.categories`
  - `session.legacyTypes`
  - `session.matchRanks`
  - `session.defaultValues`
  - `racket.statuses`
- `POST /api/sessions` 支持新增请求字段：
  - `category`
  - `subCategory`
- `PUT /api/sessions/:id` 支持新增请求字段：
  - `category`
  - `subCategory`
- `GET /api/sessions`、`GET /api/sessions/:id`、`GET /api/sessions/latest`、`GET /api/home/summary` 中的 `SessionResponse` 新增响应字段：
  - `category`
  - `categoryText`
  - `subCategory`
  - `subCategoryText`
  - `typeText`
- `type` 和 `typeLabel` 保留，用于兼容旧客户端。

### 数据库变化

- `tennis_sessions` 新增字段：
  - `category SMALLINT NOT NULL DEFAULT 1`
  - `sub_category SMALLINT NOT NULL DEFAULT 2`
- 新增迁移：`migrations/002_add_session_category.sql`。
- 新增兼容转换迁移：`migrations/003_convert_session_category_to_numeric.sql`。
- `migrations/init.sql` 同步新增字段、旧数据回填逻辑和字符串版迁移兼容转换逻辑；`category` / `sub_category` 使用 `SMALLINT`，避免 MySQL `ENUM` 扫描到 Go 数字枚举时触发 500。

### 兼容性说明

- 新增接口不涉及数据库读写，仅返回后端静态枚举定义。
- 旧客户端继续提交 `type` 可正常新增/编辑。
- 新客户端优先提交 `category` / `subCategory`。
- 新字段优先级高于旧 `type`；后端会根据新字段同步生成旧 `type`。
- 旧训练数据无法识别是否为发球训练，统一回填为 `2 / 4`。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/enum_handler.go internal/handler/session_handler.go internal/model/enum.go internal/model/enums.go internal/model/racket.go internal/model/session.go internal/service/enum_service.go internal/service/session_service.go internal/repository/session_repository.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-27 打球记录类型两级分类设计文档

### 需求/变更内容

- 新增打球记录类型两级分类改造设计文档，供设计、产品和前端评审使用。
- 明确一级类型为“日常球局、训练、比赛”。
- 明确二级类型联动规则：日常球局包含“打单、双打”，训练包含“发球、其他”，比赛包含“单打、双打”。
- 补充新增/编辑页交互、展示文案、接口影响、旧数据兼容映射和验收点。

### 修改文件

- `docs/session-type-redesign.md`
- `docs/change-log.md`

### 接口变化

- 本次仅新增设计文档，未实际修改接口实现。
- 文档建议后续调整 `POST /api/sessions`、`PUT /api/sessions/:id`、`GET /api/sessions/:id`、`GET /api/sessions`、`GET /api/sessions/latest` 和 `GET /api/home/summary` 的类型字段。

### 数据库变化

- 无。

### 兼容性说明

- 无代码变更，不影响现有接口。
- 文档中提供旧 `type` 到新两级分类的兼容映射建议。

### 已执行检查命令

- 未执行，原因：本次仅新增 Markdown 文档。

### 测试结果

- 未执行，原因：无代码变更。

## 2026-06-21 修复日历接口训练比赛占比数据

### 需求/变更内容

- 修复 `GET /api/sessions/calendar` 中训练/单打/双打/比赛占比数据缺失问题。
- 日历接口只补充统计字段，不改变原有查询范围、日期标记、费用统计和趋势逻辑。
- 恢复/补充 `charts.sessionTypeBreakdown`，用于前端展示训练比赛占比。
- 在 `summary` 中同步返回 `trainingCount`、`singlesCount`、`doublesCount`、`matchCount`，方便前端直接读取计数。

### 修改文件

- `internal/model/session.go`
- `internal/service/session_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `GET /api/sessions/calendar` 的 `summary` 增加字段：
  - `trainingCount`
  - `singlesCount`
  - `doublesCount`
  - `matchCount`
- `GET /api/sessions/calendar` 的 `charts` 增加/恢复字段：
  - `sessionTypeBreakdown`

### 数据库变化

- 无。

### 兼容性说明

- 仅新增/恢复 JSON 字段，不删除或改名已有字段。
- 日历接口原有 `weeklySessions`、`ratingTrend`、`expenseBreakdown` 保持不变。

### 已执行检查命令

- `gofmt -w internal/model/session.go internal/service/session_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-21 日历统计概览确认花费和平均时长字段

### 需求/变更内容

- 确认需求为日历页面统计概览展示花费和平均时长，不是首页统计概览。
- `GET /api/sessions/calendar` 的 `summary` 已支持：
  - `averageMinutes`：查询范围内平均单次打球时长。
  - `sessionCost`：查询范围内打球消费。
  - `racketCost`：查询范围内球拍购买费用。
  - `stringingCost`：查询范围内穿线费用。
  - `totalCost`：查询范围内总消费。
- 撤回对首页聚合接口和旧月统计接口新增平均时长字段的改动，避免扩大接口变更范围。

### 修改文件

- `internal/model/home.go`
- `internal/model/stats.go`
- `internal/repository/session_repository.go`
- `internal/service/home_service.go`
- `docs/change-log.md`

### 接口变化

- 无新增接口字段。
- 日历页面继续使用：`GET /api/sessions/calendar?year=YYYY&month=M` 或 `GET /api/sessions/calendar?year=YYYY`。
- 前端展示字段应从响应中的 `summary` 读取：
  - 平均时长：`summary.averageMinutes`
  - 总花费：`summary.totalCost`
  - 打球花费：`summary.sessionCost`

### 数据库变化

- 无。

### 兼容性说明

- 首页聚合接口未新增字段。
- 旧统计接口 `/api/stats/month` 未新增字段。
- 日历接口已有字段保持兼容。

### 已执行检查命令

- `gofmt -w internal/model/home.go internal/model/stats.go internal/repository/session_repository.go internal/service/home_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-06-16 日历接口支持按年查询

### 需求/变更内容

- 修改 `GET /api/sessions/calendar`，支持只传 `year` 参数查询全年有记录的日期。
- 保留原有 `year + month` 按月查询能力。

### 修改文件

- `internal/handler/session_handler.go`
- `internal/service/session_service.go`
- `docs/api.md`
- `AGENTS.md`

### 接口变化

- 新增支持：`GET /api/sessions/calendar?year=2026`
- 保持兼容：`GET /api/sessions/calendar?year=2026&month=6`
- 只传 `year` 时，响应 `month` 为 `0`，查询范围为该年自然年 `[YYYY-01-01, YYYY+1-01-01)`。
- 同时传 `year` 和 `month` 时，查询范围为指定自然月。

### 数据库变化

- 无。

### 兼容性说明

- 原有按月查询请求不受影响。
- 缺少 `year`、`year` 非法、`month` 非法或 `month` 不在 `1-12` 时仍返回 `invalid request`。

### 已执行检查命令

- `gofmt -w internal/handler/session_handler.go internal/service/session_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-25 打球记录时长校验前置

### 需求/变更内容

- 针对 `POST /api/sessions` 提交 `durationMinutes=1800` 的异常请求，确保服务层先执行业务参数校验。
- `Create` / `Update` 在内容安全检测前先构建并校验打球记录，避免无效请求触发外部内容安全检查或进入写入流程。
- 新增服务层单测覆盖超过 600 分钟的打球时长返回 `invalid request`。

### 修改文件

- `internal/service/session_service.go`
- `internal/service/session_service_test.go`
- `docs/change-log.md`

### 接口变化

- 无新增接口。
- `POST /api/sessions` 和 `PUT /api/sessions/:id` 对 `durationMinutes > 600` 继续返回 `invalid request`。

### 数据库变化

- 无。

### 兼容性说明

- 合法请求不受影响。
- 超过 600 分钟的打球时长请求会在内容安全检查前被拒绝。

### 已执行检查命令

- `gofmt -w internal/service/session_service.go internal/service/session_service_test.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-25 打球记录参数错误提示明确化

### 需求/变更内容

- `durationMinutes` 超出范围时，错误响应从通用 `invalid request` 调整为明确字段原因。
- 新增可携带具体原因的参数错误封装，保持 `ErrInvalidRequest` 可被 `errors.Is` 识别。
- Handler 对参数错误继续返回错误码 `40001`，但 `message` 使用 service 返回的具体原因。

### 修改文件

- `internal/service/errors.go`
- `internal/service/session_service.go`
- `internal/service/session_service_test.go`
- `internal/handler/session_handler.go`
- `docs/change-log.md`

### 接口变化

- `POST /api/sessions` 和 `PUT /api/sessions/:id` 当 `durationMinutes` 不在 `1-600` 范围时返回：`durationMinutes must be between 1 and 600`。
- 错误码仍为 `40001`，HTTP 状态码仍为 `400`。

### 数据库变化

- 无。

### 兼容性说明

- 响应结构不变。
- 仅参数错误 `message` 更明确，便于前端直接提示用户。

### 已执行检查命令

- `gofmt -w internal/service/errors.go internal/service/session_service.go internal/service/session_service_test.go internal/handler/session_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-25 打球记录参数错误改为中文提示

### 需求/变更内容

- 将打球记录参数校验失败提示从英文改为中文，便于前端直接展示给用户。
- `durationMinutes` 超出范围时提示：`打球时长必须在 1-600 分钟之间`。
- `rating` 超出范围时提示：`评分必须在 1-5 之间`。

### 修改文件

- `internal/service/session_service.go`
- `internal/service/session_service_test.go`
- `docs/change-log.md`

### 接口变化

- `POST /api/sessions` 和 `PUT /api/sessions/:id` 的参数错误 `message` 改为中文提示。
- 错误码仍为 `40001`，HTTP 状态码仍为 `400`，响应结构不变。

### 数据库变化

- 无。

### 兼容性说明

- 仅错误提示文案变化。
- 合法请求不受影响。

### 已执行检查命令

- `gofmt -w internal/service/session_service.go internal/service/session_service_test.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-25 穿线记录新增门店字段

### 需求/变更内容

- 新增穿线记录时支持填写穿线门店 `storeName`。
- 穿线记录详情、球拍列表和球拍详情中的最近穿线记录同步返回 `storeName`。
- `storeName` 纳入用户输入内容安全检查，未填写时默认为空字符串。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `migrations/init.sql`
- `migrations/011_add_stringing_store_name.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `POST /api/rackets/:id/stringing-records` 新增可选请求字段 `storeName`。
- `RacketResponse` 和 `StringingRecordResponse` 新增响应字段 `storeName`。

### 数据库变化

- `racket_stringing_record` 新增字段：`store_name VARCHAR(100) NOT NULL DEFAULT '' COMMENT '穿线门店'`。

### 兼容性说明

- `storeName` 为可选字段，旧客户端不传不受影响。
- 旧数据迁移后 `store_name` 默认为空字符串。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-26 删除球拍实例图片字段

### 需求/变更内容

- 删除 `racket.image_url` 字段，不再在用户球拍实例上保存图片地址。
- 创建/更新球拍不再接收或保存实例级 `imageUrl`。
- 球拍列表、详情、可选球拍及状态变更响应不再返回实例级 `imageUrl`。
- 保留 `racket_library.image_url`，球拍库接口图片字段不受影响。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `migrations/init.sql`
- `migrations/012_drop_racket_image_url.sql`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `POST /api/rackets` 请求字段移除 `imageUrl`。
- `PUT /api/rackets/:id` 请求字段移除 `imageUrl`。
- `RacketResponse` 响应字段移除 `imageUrl`。

### 数据库变化

- `racket` 表删除字段：`image_url`。
- 新增迁移：`migrations/012_drop_racket_image_url.sql`。

### 兼容性说明

- 旧客户端继续传 `imageUrl` 时，后端会忽略该字段。
- 如需球拍库图片，继续使用 `GET /api/racket-library` 返回的 `imageUrl` 或 `fileId`。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-26 支持编辑穿线记录

### 需求/变更内容

- 新增编辑穿线记录接口，支持修改球线名称、穿线门店、竖/横线磅数、费用和穿线时间。
- 编辑时校验当前用户、球拍归属和穿线记录归属，不能跨用户或跨球拍修改。
- 编辑后球拍列表和详情中的最近穿线信息会基于更新后的记录重新计算。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/racket_handler.go`
- `internal/model/racket.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 新增 `PUT /api/rackets/:id/stringing-records/:recordId`。
- 请求体字段同新增穿线记录：`stringName`、`storeName`、`verticalTension`、`horizontalTension`、`cost`、`stringDate`。
- 响应返回更新后的 `StringingRecordResponse`。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增接口，不影响现有新增穿线记录和球拍详情接口。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/racket_handler.go internal/model/racket.go internal/repository/racket_repository.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-26 支持删除穿线记录

### 需求/变更内容

- 新增删除穿线记录接口，支持逻辑删除指定球拍下的穿线记录。
- 删除时校验当前用户、球拍归属和穿线记录归属，不能跨用户或跨球拍删除。
- 删除后球拍列表和详情中的最近穿线信息会自动排除已删除记录。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/racket_handler.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 新增 `DELETE /api/rackets/:id/stringing-records/:recordId`。
- 响应 `data` 为 `{ "deleted": true }`。

### 数据库变化

- 无新增字段；删除穿线记录通过 `racket_stringing_record.deleted_at` 逻辑删除。

### 兼容性说明

- 仅新增接口，不影响现有新增/编辑穿线记录接口。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/racket_handler.go internal/repository/racket_repository.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-26 新增获取主力球拍接口

### 需求/变更内容

- 新增获取当前用户主力球拍接口，返回主力球拍基础信息。
- 响应复用球拍聚合信息，包含最近一次穿线记录和累计使用次数/时间。
- 无主力球拍时返回 `data: null`。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/racket_handler.go`
- `internal/repository/racket_repository.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 新增 `GET /api/rackets/primary`。
- 响应为 `RacketResponse` 或 `null`。

### 数据库变化

- 无。

### 兼容性说明

- 仅新增接口，不影响现有球拍列表、详情和主力拍设置接口。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/racket_handler.go internal/repository/racket_repository.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-28 新增我的主力球拍接口

### 需求/变更内容

- 新增 `GET /api/my-rackets/primary`，用于获取当前用户主力球拍。
- 响应复用 `RacketResponse`，包含球拍信息、最新一条穿线信息、球拍累计使用次数和累计使用小时数。
- 不再保留 `GET /api/rackets/primary` 兼容接口。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/racket_handler.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- 新增 `GET /api/my-rackets/primary`。
- 响应 `data` 为 `RacketResponse` 或 `null`。

### 数据库变化

- 无。

### 兼容性说明

- 仅保留新增接口 `GET /api/my-rackets/primary`。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/racket_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-28 移除主力球拍兼容接口

### 需求/变更内容

- 移除 `GET /api/rackets/primary` 兼容接口，仅保留 `GET /api/my-rackets/primary`。
- 将“新增、保留或扩展兼容接口/字段前必须先询问用户确认”写入 `AGENTS.md`。

### 修改文件

- `cmd/api/main.go`
- `internal/handler/racket_handler.go`
- `docs/api.md`
- `docs/change-log.md`
- `AGENTS.md`

### 接口变化

- 移除 `GET /api/rackets/primary`。
- 保留 `GET /api/my-rackets/primary`。

### 数据库变化

- 无。

### 兼容性说明

- 本次明确不保留兼容接口。
- 之后是否兼容必须先询问用户确认。

### 已执行检查命令

- `gofmt -w cmd/api/main.go internal/handler/racket_handler.go`
- `go test ./...`

### 测试结果

- 通过。

## 2026-07-28 新增球拍支持直接设为主力

### 需求/变更内容

- `POST /api/rackets` 支持请求字段 `status`。
- `status = 1` 时，新建球拍直接设为主力拍，并自动将当前用户其他主力拍改为在用。
- `status` 不传或传 `0` 时仍默认 `2` 在用；创建接口仅允许 `1` 或 `2`。

### 修改文件

- `internal/model/racket.go`
- `internal/service/racket_service.go`
- `docs/api.md`
- `docs/change-log.md`

### 接口变化

- `POST /api/rackets` 新增请求字段 `status`。

### 数据库变化

- 无。

### 兼容性说明

- 不传 `status` 的旧请求仍默认创建为在用球拍。

### 已执行检查命令

- `gofmt -w internal/model/racket.go internal/service/racket_service.go`
- `go test ./...`

### 测试结果

- 通过。
