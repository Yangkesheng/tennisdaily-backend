# Change Log

## 2026-06-27 打球记录新增两级分类字段

### 需求/变更内容

- 在打球记录中新增 `category` / `subCategory` 两级分类字段。
- 新客户端可提交一级类型和二级类型，后端校验合法组合并同步旧 `type` 字段。
- 旧客户端仍可只提交 `type`，后端自动映射为新两级分类。
- 响应新增 `category`、`categoryText`、`subCategory`、`subCategoryText`、`typeText`，方便列表、详情、首页最近记录展示。
- 历史数据通过迁移按旧 `type` 回填新字段，旧训练统一归为 `2 / 4`。

### 修改文件

- `internal/model/enum.go`
- `internal/model/session.go`
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
  - `category ENUM('1','2','3') NOT NULL DEFAULT '1'`
  - `sub_category ENUM('1','2','3','4') NOT NULL DEFAULT '2'`
- 新增迁移：`migrations/002_add_session_category.sql`。
- 新增兼容转换迁移：`migrations/003_convert_session_category_to_numeric.sql`。
- `migrations/init.sql` 同步新增字段、旧数据回填逻辑和字符串版迁移兼容转换逻辑。

### 兼容性说明

- 旧客户端继续提交 `type` 可正常新增/编辑。
- 新客户端优先提交 `category` / `subCategory`。
- 新字段优先级高于旧 `type`；后端会根据新字段同步生成旧 `type`。
- 旧训练数据无法识别是否为发球训练，统一回填为 `2 / 4`。

### 已执行检查命令

- `gofmt -w internal/model/enum.go internal/model/session.go internal/service/session_service.go internal/repository/session_repository.go`
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
