# Change Log

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
- `GET /api/session-config` 响应 `data` 直接返回 `categories`、`legacyTypes`、`matchRanks`、`defaultValues`，不再嵌套 `session` / `racket`。

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
