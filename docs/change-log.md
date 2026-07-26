# Change Log

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
