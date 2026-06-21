# AGENTS.md

> 本文档是 TennisDaily 后端仓库的 Agent 开发指南。任何 AI/Agent 在修改本仓库前，应先阅读本文，理解当前系统边界、代码结构、接口约定和开发规则。

---

## 1. 项目定位

TennisDaily 是一个网球日记小程序的 Go 后端服务。

当前后端主要负责：

- 微信/手机号登录与 JWT 鉴权
- 打球记录 CRUD
- 按日期查看打球记录
- 最近一次打球记录
- 日历月份标记
- 首页聚合统计
- 球拍库、我的球拍、球拍使用统计
- 穿线记录与装备消费统计

核心产品目标：

> 支撑小程序用户快速记录一次打球，并能长期查看自己的打球、装备和消费数据。

当前不是以下系统：

- 社交平台
- 约球平台
- AI 教练系统
- 管理后台
- 微服务架构
- 复杂 BI 报表系统

---

## 2. 技术栈

当前实际技术栈：

- Go `1.23`
- Gin
- GORM
- MySQL
- JWT
- YAML 配置
- Docker 可选

主要依赖见：

```text
go.mod
```

启动入口：

```text
cmd/api/main.go
```

配置文件：

```text
config.yaml
```

---

## 3. 当前目录结构

```text
cmd/
  api/
    main.go
internal/
  config/
    config.go
  handler/
    auth_handler.go
    home_handler.go
    racket_handler.go
    session_handler.go
    stats_handler.go
  service/
    auth_service.go
    errors.go
    home_service.go
    racket_service.go
    session_service.go
    stats_service.go
  repository/
    racket_repository.go
    session_repository.go
    user_repository.go
  model/
    auth.go
    enum.go
    home.go
    racket.go
    racket_library.go
    session.go
    stats.go
    user.go
  middleware/
    auth_middleware.go
  response/
    response.go
  logger/
    logger.go
migrations/
  *.sql
docs/
  api.md
  technical-overview.md
Dockerfile
```

职责约定：

| 目录 | 职责 |
|---|---|
| `cmd/api` | 服务启动、依赖装配、路由注册 |
| `internal/config` | 配置加载 |
| `internal/handler` | HTTP 参数绑定、鉴权用户获取、响应返回 |
| `internal/service` | 业务规则、参数默认值、统计聚合 |
| `internal/repository` | 数据库读写 |
| `internal/model` | GORM 模型、请求 DTO、响应 DTO、枚举 |
| `internal/middleware` | JWT 鉴权中间件 |
| `internal/response` | 统一响应格式 |
| `migrations` | MySQL 表结构迁移 |
| `docs` | 接口和技术文档 |

---

## 4. 分层规则

### 4.1 Handler

Handler 只做：

- 绑定 path/query/body 参数
- 获取当前登录用户 ID
- 调用 service
- 将错误映射为统一响应
- 写日志

Handler 不应：

- 直接写复杂 SQL
- 直接实现业务统计
- 操作多个 repository 拼装业务结果
- 保存业务状态

### 4.2 Service

Service 负责：

- 业务校验
- 默认值处理
- 日期解析和时区处理
- 业务聚合
- 调用 repository
- 返回 model response

Service 不应依赖：

- `gin.Context`
- HTTP 状态码
- 前端页面状态

### 4.3 Repository

Repository 负责：

- GORM 查询
- 数据库写入
- 事务
- 查询排序
- SQL 聚合

Repository 不应：

- 返回 HTTP 错误
- 处理中文业务文案
- 依赖 Handler/Service 的表现层逻辑

---

## 5. 统一响应格式

所有接口统一返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应：

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

错误码定义在：

```text
internal/response/response.go
```

当前错误码：

| code | message | 含义 |
|---:|---|---|
| `0` | `ok` | 成功 |
| `40001` | `invalid request` | 参数错误 |
| `40101` | `unauthorized` | 未登录或 token 无效 |
| `40301` | `forbidden` | 无权限 |
| `40401` | `not found` | 资源不存在 |
| `50001` | `internal error` | 服务内部错误 |

---

## 6. 鉴权规则

除以下接口外，所有 `/api/*` 业务接口均需要 JWT：

- `GET /health`
- `POST /api/auth/wechat-login`
- `POST /api/auth/phone-login`

业务接口请求头：

```http
Authorization: Bearer <token>
```

鉴权中间件：

```text
internal/middleware/auth_middleware.go
```

当前用户 ID 通过：

```go
middleware.CurrentUserID(c)
```

或 handler 内部 helper 获取。

禁止通过请求参数中的 `userId` 操作数据。所有用户数据必须来自 token 中解析出的当前用户 ID。

---

## 7. API 概览

当前主要接口：

### 7.1 健康检查

```http
GET /health
```

### 7.2 认证

```http
POST /api/auth/wechat-login
POST /api/auth/phone-login
GET  /api/auth/me
POST /api/auth/logout
PUT  /api/auth/profile
```

### 7.3 首页聚合

```http
GET /api/home/summary
GET /api/home/summary?year=2026&month=6
```

返回首页需要的聚合数据：

- 本月打球次数
- 本月累计时长
- 本月打球消费
- 本月平均时长
- 今年打球次数 `yearCount`
- 历史全部记录数 `totalCount`
- 本月球拍购买费用
- 本月穿线费用
- 本月总消费
- 最近一次打球
- 最近评分趋势

### 7.4 打球记录

```http
GET    /api/sessions
GET    /api/sessions?date=YYYY-MM-DD
POST   /api/sessions
GET    /api/sessions/latest
GET    /api/sessions/calendar?year=2026
GET    /api/sessions/calendar?year=2026&month=6
GET    /api/sessions/:id
PUT    /api/sessions/:id
DELETE /api/sessions/:id
```

说明：

- `GET /api/sessions?date=YYYY-MM-DD` 用于日历点击某一天查看当天记录。
- `GET /api/sessions/calendar` 用于日历页显示有记录的日期。
- 删除为逻辑删除。

### 7.5 统计

```http
GET /api/stats/month
```

旧首页统计接口仍保留，但首页推荐使用：

```http
GET /api/home/summary
```

### 7.6 球拍与穿线

```http
GET    /api/racket-library
GET    /api/my-rackets
GET    /api/rackets/stats
GET    /api/rackets
POST   /api/rackets
GET    /api/rackets/selectable
GET    /api/rackets/:id
PUT    /api/rackets/:id
DELETE /api/rackets/:id
POST   /api/rackets/:id/set-primary
POST   /api/rackets/:id/retire
POST   /api/rackets/:id/stringing-records
```

详细接口说明以：

```text
docs/api.md
```

为准。

---

## 8. 首页聚合接口约定

接口：

```http
GET /api/home/summary
```

支持 query：

| 参数 | 类型 | 说明 |
|---|---|---|
| `year` | number | 统计年份，不传则当前年 |
| `month` | number | 统计月份，不传则当前月 |

规则：

- `year` 和 `month` 都不传：默认当前自然月
- 只传其中一个：返回 `invalid request`
- `year` 范围：`2000-2100`
- `month` 范围：`1-12`
- 时区：`Asia/Shanghai`

核心响应结构：

```json
{
  "year": 2026,
  "month": 6,
  "session": {
    "monthCount": 8,
    "monthMinutes": 960,
    "monthCost": 320,
    "monthAverageMinutes": 120,
    "yearCount": 28,
    "totalCount": 42
  },
  "expense": {
    "sessionCost": 320,
    "racketCost": 1280,
    "stringingCost": 260,
    "totalCost": 1860
  },
  "latestSession": null,
  "ratingTrend": []
}
```

字段口径：

| 字段 | 口径 |
|---|---|
| `monthCount` | 查询月份打球记录数 |
| `monthMinutes` | 查询月份打球总分钟数 |
| `monthCost` | 查询月份 `tennis_sessions.cost` 合计 |
| `monthAverageMinutes` | 查询月份单次打球平均时长，单位分钟，按 `AVG(duration_minutes)` 计算，保留 1 位小数 |
| `yearCount` | 查询年份打球记录数 |
| `totalCount` | 历史累计所有未删除打球记录数 |
| `sessionCost` | 同 `monthCost` |
| `racketCost` | 查询月份球拍购买费用 |
| `stringingCost` | 查询月份穿线费用 |
| `totalCost` | `sessionCost + racketCost + stringingCost` |

---

## 9. 数据模型与枚举

### 9.1 打球类型 `SessionType`

定义在：

```text
internal/model/enum.go
```

当前枚举：

| 值 | 含义 |
|---:|---|
| `1` | 双打 |
| `2` | 单打 |
| `3` | 训练 |
| `4` | 单打比赛 |
| `5` | 双打比赛 |

### 9.2 比赛成绩 `MatchRank`

| 值 | 含义 |
|---:|---|
| `0` | 无 |
| `1` | 冠军 |
| `2` | 亚军 |
| `3` | 四强 |
| `4` | 八强 |
| `5` | 小组赛 |

业务规则：

- 只有单打比赛、双打比赛可以保留 `matchRank`
- 普通双打、单打、训练会强制 `matchRank = 0`

### 9.3 球拍状态

定义在：

```text
internal/model/racket.go
```

| 值 | 含义 |
|---:|---|
| `1` | 主力拍 |
| `2` | 在用 |
| `3` | 已退役 |

---

## 10. 日期与时区规则

### 10.1 API 日期格式

对外日期字段使用：

```text
YYYY-MM-DD
```

例如：

```json
{
  "date": "2026-06-13"
}
```

不要让前端传 RFC3339 作为打球日期。

### 10.2 Go 解析日期

Service 层使用：

```go
time.ParseInLocation("2006-01-02", date, loc)
```

### 10.3 统计时区

统计统一使用：

```text
Asia/Shanghai
```

自然月范围使用左闭右开：

```text
[start, end)
```

例如 2026 年 6 月：

```text
2026-06-01 <= date < 2026-07-01
```

---

## 11. 软删除规则

当前关键业务表使用逻辑删除。

打球记录：

```text
tennis_sessions.deleted_at
```

球拍：

```text
racket.deleted_at
```

穿线记录：

```text
racket_stringing_record.deleted_at
```

普通查询、统计、列表必须过滤未删除数据。

注意：GORM 的 `gorm.DeletedAt` 在 `Model(&Struct{})` 场景通常会自动过滤软删除，但如果使用 `Table(...)`、自定义 join 或 raw SQL，必须显式处理 `deleted_at IS NULL`。

---

## 12. 数据库表名约定

当前实际表名包括：

```text
users
tennis_sessions
racket
racket_stringing_record
racket_library
```

注意：球拍表名是单数：

```text
racket
```

穿线表名是：

```text
racket_stringing_record
```

不要误写为 `rackets` 或 `stringing_records`。

---

## 13. JSON 字段规则

API JSON 字段使用 camelCase：

```json
{
  "durationMinutes": 120,
  "matchRank": 1,
  "courtName": "奥森网球场",
  "racketId": 12
}
```

数据库字段使用 snake_case：

```text
duration_minutes
match_rank
court_name
racket_id
```

不要把数据库字段名直接暴露给前端。

---

## 14. 业务校验规则

### 14.1 打球记录

- `date` 必填，格式 `YYYY-MM-DD`
- `durationMinutes` 默认 `120`
- `durationMinutes` 必须 `> 0` 且不超过 `600`
- `rating` 默认 `3`
- `rating` 范围 `1-5`
- `type` 必须是合法 `SessionType`
- `matchRank` 必须是合法 `MatchRank`
- 非比赛类型强制 `matchRank = 0`

### 14.2 球拍

- `name` 必填
- `libraryId` 可选
- `purchaseDate` 可选，格式 `YYYY-MM-DD`
- `purchasePrice` 可选
- 设置主力拍时，同一用户其他主力拍应自动变为在用

### 14.3 穿线记录

- `stringName` 必填
- `stringDate` 必填，格式 `YYYY-MM-DD`
- `cost` 数字，允许 0
- 穿线记录归属当前用户和指定球拍

---

## 15. 日志规则

项目有简单 logger：

```text
internal/logger/logger.go
```

Handler 中建议记录：

- 请求开始
- 当前用户 ID
- 关键 query/path 参数
- 请求成功和关键结果数量

不要在日志里打印：

- JWT 完整 token
- 微信 session_key
- 手机号敏感信息
- 过长请求体

---

## 16. 开发命令

格式化：

```bash
gofmt -w <files>
```

运行测试/编译检查：

```bash
go test ./...
```

本地启动：

```bash
go run ./cmd/api
```

Docker 构建如需使用：

```bash
docker build -t tennisdaily-backend .
```

Agent 修改 Go 代码后必须至少运行：

```bash
go test ./...
```

---

## 17. 修改代码时的检查清单

修改或新增接口时，必须检查：

1. 是否需要鉴权
2. 是否只操作当前用户数据
3. 是否遵守统一响应格式
4. 是否正确映射 service 错误
5. 是否过滤软删除数据
6. 日期是否使用 `Asia/Shanghai`
7. JSON 字段是否 camelCase
8. 是否更新 `docs/api.md`、`docs/change-log.md` 或相关接口文档
9. 是否执行 `gofmt`
10. 是否执行 `go test ./...`

---

## 18. 禁止事项

禁止：

- 引入微服务、消息队列、Redis 缓存等过度架构
- Handler 中写复杂业务或 SQL
- 通过前端传入的 `userId` 决定数据归属
- 物理删除核心业务数据，除非明确是清理任务
- 随意改变现有 API 字段名
- 返回 snake_case JSON 字段
- 大量使用 `map[string]interface{}` 表达核心业务数据
- 在 Go 代码中散落枚举魔法数字
- 统计接口漏掉 `deleted_at` 规则
- 把中文文案作为数据库枚举主值

---

## 19. 推荐开发顺序

新增功能建议按以下顺序：

1. 明确接口文档和响应结构
2. 在 `internal/model` 增加 request/response DTO
3. 在 repository 增加最小必要查询
4. 在 service 实现业务规则和聚合
5. 在 handler 暴露接口
6. 在 `cmd/api/main.go` 注册路由
7. 更新 `docs/api.md`
8. 执行 `gofmt` 和 `go test ./...`

---

## 20. 当前系统重点注意

### 20.1 首页推荐使用聚合接口

首页应优先使用：

```http
GET /api/home/summary
```

而不是让前端分别请求：

```http
GET /api/sessions
GET /api/stats/month
GET /api/sessions/latest
```

### 20.2 日历点击日期

日历点击有打球记录的日期时，前端应请求：

```http
GET /api/sessions?date=YYYY-MM-DD
```

后端已支持 `date` query，不需要前端拉全量后过滤。

### 20.3 月份日历

日历页可按年或按月使用：

```http
GET /api/sessions/calendar?year=YYYY
GET /api/sessions/calendar?year=YYYY&month=M
```

只传 `year` 时返回全年有记录的日期；同时传 `year` 和 `month` 时返回指定月份有记录的日期。响应中 `activeDayCount` 表示查询范围内有打球记录的天数，不是打球记录条数。

### 20.4 费用统计

首页 `expense.totalCost` 口径：

```text
打球消费 + 球拍购买费用 + 穿线费用
```

其中：

- 打球消费：`tennis_sessions.cost`
- 球拍购买费用：`racket.purchase_price`
- 穿线费用：`racket_stringing_record.cost`

---

## 21. 文档维护

主要文档：

```text
docs/api.md
docs/change-log.md
AGENTS.md
```

当接口有新增、字段变化或统计口径变化时，必须更新文档。

每次修改代码后，必须在 `docs/change-log.md` 追加本次修改总结，建议包含：

- 修改时间
- 需求/变更内容
- 修改文件
- 接口变化
- 数据库变化
- 兼容性说明
- 已执行检查命令
- 测试结果

旧文档：

```text
AGENTS.go-backend.md
```

仅作为历史参考；新的 Agent 指南以本文件为准。
