# Go 后端 Agent 开发指南

> 本文档用于指导 AI/Agent 开发网球日记项目的自建 Go 后端服务。目标是在不破坏前端“30 秒快速记录”体验的前提下，逐步接入登录、数据同步、记录 CRUD 和统计能力。

---

# 一、后端定位

Go 后端的核心目标不是做复杂平台，而是为网球日记小程序提供稳定的数据服务。

第一阶段目标：

- 微信小程序登录
- 打球记录云端保存
- 打球记录 CRUD
- 最近一次打球记录
- 本月统计
- 为后续装备管理、费用统计、数据同步做基础

当前不做：

- 社交系统
- 约球系统
- AI 教练
- 管理后台
- 微服务
- 复杂权限系统
- 复杂报表系统

核心原则：

> 后端服务必须服务于“快速记录”，不能让网络请求、复杂表单或同步逻辑拖慢用户记录体验。

---

# 二、推荐技术栈

## 1. 语言与框架

- Go
- Gin 作为 HTTP 框架
- PostgreSQL 作为主数据库
- GORM 作为 MVP 阶段 ORM
- JWT 用于后端登录态
- Docker 可选

## 2. 为什么这样选

Gin：

- 上手快
- 生态成熟
- 适合 REST API
- AI 生成代码稳定

PostgreSQL：

- 适合长期数据存储
- 统计能力强
- 后续月度、年度、装备统计更方便

GORM：

- MVP 开发效率高
- 类型结构清晰
- 后续如果复杂度上升，可以逐步迁移到 sqlc

---

# 三、后端开发原则

## 1. 单体服务优先

当前阶段只允许单体服务。

禁止：

- 微服务拆分
- 事件总线
- 消息队列优先
- 复杂 DDD 分层
- 过度抽象 repository/service

允许：

- 简单分层
- 清晰职责
- 可测试的业务函数

---

## 2. 类型安全

所有核心业务对象必须定义 Go struct。

禁止在业务核心逻辑里大量使用：

```go
map[string]interface{}
```

核心对象必须有明确类型，例如：

- User
- TennisSession
- SessionStats
- CreateSessionRequest
- UpdateSessionRequest
- LoginRequest
- LoginResponse

---

## 3. Handler 不写复杂业务

Handler 只负责：

- 绑定参数
- 获取当前用户
- 调用 service
- 返回响应

禁止在 Handler 中直接写复杂业务逻辑或数据库查询。

---

## 4. Service 负责业务规则

Service 负责：

- 默认值处理
- 业务校验
- 类型/成绩合法性判断
- 统计聚合
- 调用 repository

Service 层不应该依赖 Gin Context。

---

## 5. Repository 只负责数据读写

Repository 负责：

- 数据库查询
- 数据库写入
- 事务处理

Repository 不负责 HTTP，不负责业务文案。

---

## 6. API 字段使用 camelCase

接口 JSON 字段统一使用 camelCase。

示例：

```json
{
  "durationMinutes": 120,
  "matchRank": 1,
  "courtName": "奥森网球场"
}
```

数据库字段统一使用 snake_case。

示例：

```sql
duration_minutes
match_rank
court_name
created_at
```

---

# 四、目录结构约定

推荐后端目录：

```text
cmd/
  api/
    main.go
internal/
  config/
    config.go
  handler/
    auth_handler.go
    session_handler.go
    stats_handler.go
  service/
    auth_service.go
    session_service.go
    stats_service.go
  repository/
    user_repository.go
    session_repository.go
  model/
    user.go
    session.go
    stats.go
    enum.go
  middleware/
    auth_middleware.go
  response/
    response.go
migrations/
  001_init.sql
go.mod
go.sum
```

## 目录职责

| 目录 | 职责 |
|---|---|
| `cmd/api` | 服务启动入口 |
| `internal/config` | 配置加载 |
| `internal/handler` | HTTP Handler |
| `internal/service` | 业务逻辑 |
| `internal/repository` | 数据库读写 |
| `internal/model` | 数据模型、枚举、DTO |
| `internal/middleware` | 鉴权中间件 |
| `internal/response` | 统一响应结构 |
| `migrations` | SQL 迁移文件 |

---

# 五、MVP API 范围

第一阶段只实现以下 API。

## 1. 健康检查

```text
GET /health
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "ok"
  }
}
```

---

## 2. 微信登录

```text
POST /api/auth/wechat-login
```

请求：

```json
{
  "code": "wx_login_code"
}
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "jwt-token",
    "userId": 1,
    "openid": "openid"
  }
}
```

说明：

- 小程序调用 `wx.login()` 获取 code
- 后端用 code 调微信接口换 openid
- 后端创建或查询用户
- 返回自有 JWT

---

## 3. 打球记录 CRUD

```text
GET    /api/sessions
POST   /api/sessions
GET    /api/sessions/:id
PUT    /api/sessions/:id
DELETE /api/sessions/:id
```

---

## 4. 最近一次打球

```text
GET /api/sessions/latest
```

---

## 5. 本月统计

```text
GET /api/stats/month
```

---

# 六、统一响应格式

所有接口统一响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

错误响应示例：

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

## 常用错误码建议

| code | 含义 |
|---|---|
| 0 | 成功 |
| 40001 | 请求参数错误 |
| 40101 | 未登录或 token 无效 |
| 40301 | 无权限 |
| 40401 | 资源不存在 |
| 50001 | 服务内部错误 |

---

# 七、数据库设计约定

## 1. 用户表

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  openid VARCHAR(128) NOT NULL UNIQUE,
  nickname VARCHAR(128) NOT NULL DEFAULT '',
  avatar_url TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 2. 打球记录表

```sql
CREATE TABLE tennis_sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),

  date DATE NOT NULL,
  duration_minutes INT NOT NULL DEFAULT 120,
  rating SMALLINT NOT NULL DEFAULT 3,

  type SMALLINT NOT NULL DEFAULT 1,
  match_rank SMALLINT NOT NULL DEFAULT 0,

  court_name VARCHAR(128) NOT NULL DEFAULT '',
  cost NUMERIC(10,2) NOT NULL DEFAULT 0,
  racket_name VARCHAR(128) NOT NULL DEFAULT '',
  shoe_name VARCHAR(128) NOT NULL DEFAULT '',
  note TEXT NOT NULL DEFAULT '',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
```

## 3. 索引

```sql
CREATE INDEX idx_tennis_sessions_user_date ON tennis_sessions(user_id, date DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_tennis_sessions_user_created ON tennis_sessions(user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_tennis_sessions_user_deleted ON tennis_sessions(user_id, deleted_at);
```

## 4. 约束

```sql
ALTER TABLE tennis_sessions
ADD CONSTRAINT tennis_sessions_type_check
CHECK (type IN (1, 2, 3, 4, 5));

ALTER TABLE tennis_sessions
ADD CONSTRAINT tennis_sessions_match_rank_check
CHECK (match_rank IN (0, 1, 2, 3, 4, 5));

ALTER TABLE tennis_sessions
ADD CONSTRAINT tennis_sessions_rating_check
CHECK (rating >= 1 AND rating <= 5);
```

## 5. 逻辑删除

打球记录必须优先使用逻辑删除（Soft Delete），不建议物理删除。

原因：

- 未来需要支持多设备同步
- 如果 A 设备删除了一条记录，B 设备离线期间仍保留旧数据
- 后端如果物理删除，B 设备增量同步时很难知道该记录已经被删除
- 使用 `deleted_at` 可以让客户端通过增量同步感知删除事件

实现建议：

- `tennis_sessions` 表增加 `deleted_at TIMESTAMPTZ`
- GORM 可使用 `gorm.DeletedAt` 或显式 `*time.Time`
- 默认查询必须过滤 `deleted_at IS NULL`
- 删除接口执行软删除，只更新 `deleted_at` 和 `updated_at`
- 只有未来做数据清理任务时，才考虑定期物理清除很久以前的已删除数据

DELETE 接口语义仍然保持：

```text
DELETE /api/sessions/:id
```

但数据库实现不使用真正的 `DELETE FROM tennis_sessions`，而是：

```sql
UPDATE tennis_sessions
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = ? AND user_id = ? AND deleted_at IS NULL;
```

---

# 八、枚举约定

## 1. 打球类型 SessionType

数据库使用 SMALLINT。

| 数字 | 含义 | 前端当前字符串 |
|---|---|---|
| 1 | 双打 | `doubles` |
| 2 | 单打 | `singles` |
| 3 | 训练 | `training` |
| 4 | 单打比赛 | `singlesMatch` |
| 5 | 双打比赛 | `doublesMatch` |

Go 代码必须定义常量，禁止散落魔法数字。

```go
type SessionType int16

const (
	SessionTypeDoubles      SessionType = 1
	SessionTypeSingles      SessionType = 2
	SessionTypeTraining     SessionType = 3
	SessionTypeSinglesMatch SessionType = 4
	SessionTypeDoublesMatch SessionType = 5
)
```

---

## 2. 比赛成绩 MatchRank

| 数字 | 含义 | 前端当前字符串 |
|---|---|---|
| 0 | 无 | `''` |
| 1 | 冠军 | `champion` |
| 2 | 亚军 | `runnerUp` |
| 3 | 四强 | `semiFinal` |
| 4 | 八强 | `quarterFinal` |
| 5 | 小组赛 | `groupStage` |

```go
type MatchRank int16

const (
	MatchRankNone         MatchRank = 0
	MatchRankChampion     MatchRank = 1
	MatchRankRunnerUp     MatchRank = 2
	MatchRankSemiFinal    MatchRank = 3
	MatchRankQuarterFinal MatchRank = 4
	MatchRankGroupStage   MatchRank = 5
)
```

---

# 九、核心 Go Model 建议

## 1. TennisSession

```go
type TennisSession struct {
	ID              int64       `json:"id" gorm:"primaryKey"`
	UserID          int64       `json:"userId"`
	Date            time.Time   `json:"date"`
	DurationMinutes int         `json:"durationMinutes"`
	Rating          int16       `json:"rating"`
	Type            SessionType `json:"type"`
	MatchRank       MatchRank   `json:"matchRank"`
	CourtName       string      `json:"courtName"`
	Cost            float64     `json:"cost"`
	RacketName      string      `json:"racketName"`
	ShoeName        string      `json:"shoeName"`
	Note            string      `json:"note"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}
```

注意：

- 数据库 date 用 DATE
- Go 内部可使用 time.Time
- API 层可格式化为 `YYYY-MM-DD`
- 如果使用 `gorm.DeletedAt`，对外 JSON 默认不返回 `deletedAt`
- 普通列表、详情、统计接口默认只查询未删除记录

---

## 2. CreateSessionRequest

```go
type CreateSessionRequest struct {
	Date            string      `json:"date" binding:"required"`
	DurationMinutes int         `json:"durationMinutes"`
	Rating          int16       `json:"rating"`
	Type            SessionType `json:"type" binding:"required"`
	MatchRank       MatchRank   `json:"matchRank"`
	CourtName       string      `json:"courtName"`
	Cost            float64     `json:"cost"`
	RacketName      string      `json:"racketName"`
	ShoeName        string      `json:"shoeName"`
	Note            string      `json:"note"`
}
```

---

# 十、业务校验规则

## 1. 日期

- 必须是 `YYYY-MM-DD`
- 后端保存为 DATE
- 不要依赖 Go 对 `time.Time` 的默认 JSON 反序列化来解析该字段，因为默认格式是 RFC3339，不是 `YYYY-MM-DD`
- `CreateSessionRequest.Date` 建议保持 string 类型
- Service 层必须使用标准日期解析：

```go
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil {
	return err
}

sessionDate, err := time.ParseInLocation("2006-01-02", req.Date, loc)
if err != nil {
	return err
}
```

- 统计自然月时也应基于同一时区，避免跨时区导致日期归属错误

## 2. 时长

- 默认 120 分钟
- 必须大于 0
- 建议最大不超过 600 分钟

## 3. 手感评分

- 范围 1 到 5
- 默认 3

## 4. 类型

允许值：

```text
1, 2, 3, 4, 5
```

## 5. 比赛成绩

允许值：

```text
0, 1, 2, 3, 4, 5
```

规则：

- 当 type 为单打比赛或双打比赛时，可以填写 matchRank
- 当 type 为双打、单打、训练时，matchRank 应强制为 0

---

# 十一、前后端数据兼容策略

当前小程序前端 MVP 使用字符串类型：

```ts
type: 'doublesMatch'
matchRank: 'champion'
```

Go 后端推荐使用数字：

```json
{
  "type": 5,
  "matchRank": 1
}
```

因此接入后端时必须有 mapper 层。

## 1. 前端原则

- 页面层不直接散落后端数字枚举
- WXML 中不要出现业务魔法数字，例如 `draft.type === 5`
- 应通过 constants 或 mapper 统一转换
- 页面继续调用 `services/session-service.ts`
- 页面不直接调用 `wx.request`

## 2. 后端原则

API 可以返回：

```json
{
  "type": 5,
  "typeLabel": "双打比赛",
  "matchRank": 1,
  "matchRankLabel": "冠军"
}
```

这样前端展示更简单，调试也更直观。

---

# 十二、鉴权设计

## 1. 小程序登录流程

```text
小程序 wx.login()
  -> 获取 code
  -> POST /api/auth/wechat-login
  -> Go 后端请求微信 code2session
  -> 获取 openid
  -> 查询或创建 user
  -> 签发 JWT
  -> 小程序本地保存 token
```

## 2. 请求鉴权

后续业务接口请求头：

```text
Authorization: Bearer <token>
```

## 3. Handler 获取当前用户

鉴权中间件解析 JWT 后，将 userID 放入 context。

业务接口必须基于 userID 查询自己的数据。

禁止用户通过传 userId 操作他人数据。

---

# 十三、统计接口规则

## 1. 本月统计

接口：

```text
GET /api/stats/month
```

返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "monthCount": 8,
    "monthMinutes": 960,
    "monthCost": 320,
    "totalCount": 42
  }
}
```

统计口径：

- 当前自然月
- 按 session.date 计算，不按 created_at
- 只统计当前登录用户

---

# 十四、性能和同步原则

当前数据量很小，不需要复杂优化。

必须避免：

- 为小数据引入 Redis
- 过早做缓存
- 过早做消息队列
- 复杂异步同步架构

允许后续做：

- 本地优先保存
- 后台同步远程
- 失败重试
- 冲突时以后端 updated_at 为准或提示用户

核心体验要求：

> 用户保存一条打球记录时，不应该因为网络慢而明显卡顿。

---

# 十五、测试建议

MVP 阶段至少覆盖：

- 创建记录
- 查询列表
- 查询单条
- 更新记录
- 删除记录，必须验证软删除后普通查询不再返回，且增量同步场景仍可识别 deleted_at
- 比赛类型成绩校验
- 非比赛类型清空成绩
- 本月统计
- 用户只能访问自己的记录

建议优先写 service 层测试。

---

# 十六、禁止事项

禁止：

- 一开始做微服务
- 一开始引入复杂权限模型
- 页面直接调用后端 API 绕过 service
- 后端业务逻辑散落在 handler
- Go 代码里到处写业务魔法数字
- 数据库直接保存中文业务文案作为枚举主值
- 对打球记录执行物理删除，除非是明确的数据清理任务
- 为 MVP 引入过多中间件
- 强制用户联网才能记录

---

# 十七、开发顺序建议

推荐按以下顺序开发：

1. 初始化 Go 项目
2. 加入 Gin
3. 加入配置加载
4. 加入 PostgreSQL / GORM
5. 建 users、tennis_sessions 表
6. 实现统一响应
7. 实现健康检查
8. 实现微信登录
9. 实现 JWT 中间件
10. 实现打球记录 CRUD
11. 实现最近一次记录
12. 实现本月统计
13. 小程序 service 层接入 API
14. 增加本地缓存/离线同步策略

---

# 十八、小结

Go 后端的第一目标是让当前小程序从本地记录平滑升级到云端记录。

架构保持：

```text
小程序 pages
  -> miniprogram/services
  -> Go REST API
  -> PostgreSQL
```

不要牺牲当前 MVP 最重要的能力：

> 打完球后 30 秒内完成一次记录。
