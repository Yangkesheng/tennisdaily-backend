# TennisDaily Go 服务端接口文档

> 本文档描述 TennisDaily 自建 Go 后端当前已实现的 HTTP API，供小程序前端联调和后续维护使用。

---

## 1. 基础信息

### 1.1 服务说明

服务端用于为网球日记小程序提供：

- 微信登录
- JWT 鉴权
- 打球记录新增、查询、编辑、删除
- 最近一次打球记录
- 本月统计
- 球拍管理和球拍消费统计

### 1.2 Base URL

本地开发环境默认端口来自当前 `config.yaml`：

```text
server.port = 8081
```

因此当前本地 Base URL 为：

```text
http://localhost:8081
```

真机调试时不要使用 `localhost`，需要使用电脑局域网 IP，例如：

```text
http://192.168.x.x:8081
```

正式环境建议使用 HTTPS 域名。

### 1.3 当前本地配置摘要

当前本地配置如下：

| 配置项 | 当前值 | 说明 |
|---|---|---|
| server.port | `8081` | Go API 服务端口 |
| database.host | `localhost` | MySQL 主机 |
| database.port | `3306` | MySQL 端口 |
| database.name | `tennis_diary` | 数据库名 |
| database.charset | `utf8mb4` | 字符集 |
| database.loc | `Asia/Shanghai` | 数据库时区 |
| jwt.expireHours | `720` | JWT 有效期，约 30 天 |
| wechat.appId | 空 | 本地未配置微信正式 AppID |
| wechat.appSecret | 空 | 本地未配置微信正式密钥 |

由于当前 `wechat.appId` 和 `wechat.appSecret` 为空，本地登录会走开发模式：

```text
openid = "dev_openid_" + code
```

也就是说，同一个测试 code 会得到稳定的开发 openid，方便本地联调。

### 1.4 Content-Type

除健康检查外，请求体均使用 JSON：

```http
Content-Type: application/json
```

### 1.5 鉴权方式

除以下接口外，其余 `/api/*` 接口均需要登录鉴权：

- `GET /health`
- `POST /api/auth/wechat-login`
- `POST /api/auth/phone-login`

登录后请求头携带：

```http
Authorization: Bearer <token>
```

---

## 2. 统一响应格式

所有接口统一返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 2.1 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 2.2 错误响应

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

### 2.3 错误码

| HTTP 状态码 | code | message | 含义 |
|---|---:|---|---|
| 200 | 0 | ok | 成功 |
| 400 | 40001 | invalid request | 请求参数错误 |
| 401 | 40101 | unauthorized | 未登录或 token 无效 |
| 403 | 40301 | forbidden | 无权限 |
| 404 | 40401 | not found | 资源不存在 |
| 500 | 50001 | internal error | 服务内部错误 |

---

## 3. 枚举定义

### 3.1 打球类型 category / subCategory

新增和编辑打球记录推荐使用两级分类字段：

| category | categoryText | subCategory | subCategoryText | typeText |
|---:|---|---:|---|---|
| 1 | 日常球局 | 1 | 打单 | 日常球局 · 打单 |
| 1 | 日常球局 | 2 | 双打 | 日常球局 · 双打 |
| 2 | 训练 | 3 | 发球 | 训练 · 发球 |
| 2 | 训练 | 4 | 其他 | 训练 · 其他 |
| 3 | 比赛 | 1 | 单打 | 比赛 · 单打 |
| 3 | 比赛 | 2 | 双打 | 比赛 · 双打 |

后端会校验 `category` 和 `subCategory` 组合是否合法。

### 3.2 旧打球类型 type

`type` 为兼容旧客户端和旧统计逻辑保留。新客户端优先提交和展示 `category` / `subCategory`。

| 值 | 含义 | typeLabel | 映射 category | 映射 subCategory |
|---:|---|---|---:|---:|
| 1 | 双打 | 双打 | 1 | 2 |
| 2 | 单打 | 单打 | 1 | 1 |
| 3 | 训练 | 训练 | 2 | 4 |
| 4 | 单打比赛 | 单打比赛 | 3 | 1 |
| 5 | 双打比赛 | 双打比赛 | 3 | 2 |

说明：旧训练数据统一映射为 `2 / 4`（训练 / 其他），无法自动识别为发球训练。

### 3.3 比赛成绩 matchRank

| 值 | 含义 | matchRankLabel |
|---:|---|---|
| 0 | 无 | 空字符串 |
| 1 | 冠军 | 冠军 |
| 2 | 亚军 | 亚军 |
| 3 | 四强 | 四强 |
| 4 | 八强 | 八强 |
| 5 | 小组赛 | 小组赛 |

说明：

- 当 `category` 为 `3` 时，`matchRank` 可以为 `1-5`
- 当 `category` 为 `1` 或 `2` 时，后端会强制将 `matchRank` 处理为 `0`
- 旧客户端仍可通过 `type = 4` 或 `type = 5` 提交比赛成绩

---

## 4. 数据模型

### 4.1 SessionResponse

打球记录响应结构：

```json
{
  "id": 1,
  "date": "2026-01-15",
  "durationMinutes": 120,
  "rating": 3,
  "type": 5,
  "typeLabel": "双打比赛",
  "category": 3,
  "categoryText": "比赛",
  "subCategory": 2,
  "subCategoryText": "双打",
  "typeText": "比赛 · 双打",
  "matchRank": 1,
  "matchRankLabel": "冠军",
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 80,
  "racketName": "Wilson Blade",
  "shoeName": "Asics Gel Resolution",
  "note": "今天状态不错",
  "createdAt": "2026-01-15T12:00:00+08:00",
  "updatedAt": "2026-01-15T12:00:00+08:00"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 记录 ID |
| date | string | 打球日期，格式 `YYYY-MM-DD` |
| durationMinutes | number | 打球时长，单位分钟 |
| rating | number | 今日手感，1-5 |
| type | number | 旧打球类型枚举，兼容保留 |
| typeLabel | string | 旧打球类型中文文案 |
| category | number | 一级类型枚举：`1` 日常球局、`2` 训练、`3` 比赛 |
| categoryText | string | 一级类型中文文案 |
| subCategory | number | 二级类型枚举：`1` 单打/打单、`2` 双打、`3` 发球、`4` 其他 |
| subCategoryText | string | 二级类型中文文案 |
| typeText | string | 两级类型组合展示文案 |
| matchRank | number | 比赛成绩枚举 |
| matchRankLabel | string | 比赛成绩中文文案 |
| courtName | string | 场地名称 |
| partner | string | 搭档 |
| cost | number | 花费 |
| racketName | string | 球拍 |
| shoeName | string | 球鞋 |
| note | string | 备注 |
| createdAt | string | 创建时间，RFC3339 |
| updatedAt | string | 更新时间，RFC3339 |

### 4.2 CreateSessionRequest / UpdateSessionRequest

新增和更新记录请求结构一致：

```json
{
  "date": "2026-01-15",
  "durationMinutes": 120,
  "rating": 3,
  "category": 3,
  "subCategory": 2,
  "matchRank": 1,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 80,
  "racketName": "Wilson Blade",
  "shoeName": "Asics Gel Resolution",
  "note": "今天状态不错"
}
```

字段规则：

| 字段 | 是否必填 | 规则 |
|---|---|---|
| date | 是 | 格式 `YYYY-MM-DD` |
| durationMinutes | 否 | 默认 120；必须大于 0，最大 600 |
| rating | 否 | 默认 3；范围 1-5 |
| category | 新客户端必填 | 一级类型，允许 `1`、`2`、`3` |
| subCategory | 新客户端必填 | 二级类型，允许 `1`、`2`、`3`、`4`；必须属于当前 `category` |
| type | 旧客户端必填 | 旧类型枚举，允许 1-5；未传新字段时后端按旧类型映射 |
| matchRank | 否 | 允许 0-5；非比赛类型会被强制改为 0 |
| courtName | 否 | 字符串 |
| partner | 否 | 字符串，搭档名称 |
| cost | 否 | 数字 |
| racketName | 否 | 字符串 |
| shoeName | 否 | 字符串 |
| note | 否 | 字符串 |

注意：

- 新客户端应提交 `category` 和 `subCategory`，后端会同步生成兼容旧字段 `type`
- 旧客户端仍可只提交 `type`，后端会自动映射出 `category` 和 `subCategory`
- `date` 必须是 `YYYY-MM-DD`，不是完整 ISO 时间
- `durationMinutes` 传 `0` 时后端会使用默认值 `120`
- `rating` 传 `0` 时后端会使用默认值 `3`

---

## 5. 接口列表总览

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| GET | `/health` | 否 | 健康检查 |
| POST | `/api/auth/wechat-login` | 否 | 微信登录 |
| POST | `/api/auth/phone-login` | 否 | 手机号登录 |
| GET | `/api/enums` | 是 | 获取前端枚举值和展示文案 |
| GET | `/api/sessions` | 是 | 获取打球记录列表 |
| POST | `/api/sessions` | 是 | 新增打球记录 |
| GET | `/api/sessions/latest` | 是 | 获取最近一次打球记录 |
| GET | `/api/sessions/calendar` | 是 | 获取日历记录标记，支持按月或按年 |
| GET | `/api/sessions/:id` | 是 | 获取单条打球记录 |
| PUT | `/api/sessions/:id` | 是 | 更新打球记录 |
| DELETE | `/api/sessions/:id` | 是 | 删除打球记录，软删除 |
| GET | `/api/stats/month` | 是 | 获取本月统计 |
| GET | `/api/rackets/stats` | 是 | 获取球拍统计 |

---

## 6. 枚举值

### 6.1 GET /api/enums

获取前端展示和表单选择所需的枚举值及含义。

#### 请求

```http
GET /api/enums
Authorization: Bearer <token>
```

#### 响应 data

```json
{
  "session": {
    "categories": [
      {
        "value": 1,
        "label": "日常球局",
        "subCategories": [
          {
            "value": 1,
            "label": "打单",
            "category": 1,
            "typeText": "日常球局 · 打单",
            "legacyType": 2
          },
          {
            "value": 2,
            "label": "双打",
            "category": 1,
            "typeText": "日常球局 · 双打",
            "legacyType": 1
          }
        ]
      },
      {
        "value": 2,
        "label": "训练",
        "subCategories": [
          {
            "value": 3,
            "label": "发球",
            "category": 2,
            "typeText": "训练 · 发球",
            "legacyType": 3
          },
          {
            "value": 4,
            "label": "其他",
            "category": 2,
            "typeText": "训练 · 其他",
            "legacyType": 3
          }
        ]
      },
      {
        "value": 3,
        "label": "比赛",
        "subCategories": [
          {
            "value": 1,
            "label": "单打",
            "category": 3,
            "typeText": "比赛 · 单打",
            "legacyType": 4
          },
          {
            "value": 2,
            "label": "双打",
            "category": 3,
            "typeText": "比赛 · 双打",
            "legacyType": 5
          }
        ]
      }
    ],
    "legacyTypes": [
      { "value": 1, "label": "双打" },
      { "value": 2, "label": "单打" },
      { "value": 3, "label": "训练" },
      { "value": 4, "label": "单打比赛" },
      { "value": 5, "label": "双打比赛" }
    ],
    "matchRanks": [
      { "value": 0, "label": "" },
      { "value": 1, "label": "冠军" },
      { "value": 2, "label": "亚军" },
      { "value": 3, "label": "四强" },
      { "value": 4, "label": "八强" },
      { "value": 5, "label": "小组赛" }
    ],
    "defaultValues": {
      "category": 1,
      "subCategory": 2,
      "durationMinutes": 120,
      "rating": 3,
      "matchRank": 0
    }
  },
  "racket": {
    "statuses": [
      { "value": 1, "label": "主力拍" },
      { "value": 2, "label": "在用" },
      { "value": 3, "label": "已退役" }
    ]
  }
}
```

#### 字段说明

| 字段 | 说明 |
|---|---|
| `session.categories` | 打球记录两级分类，前端可直接用于一级/二级联动选择 |
| `subCategories[].legacyType` | 当前二级分类对应的旧 `type`，用于旧客户端兼容展示 |
| `session.legacyTypes` | 旧打球类型枚举，兼容保留 |
| `session.matchRanks` | 比赛成绩枚举 |
| `session.defaultValues` | 新增打球记录推荐默认值 |
| `racket.statuses` | 球拍状态枚举 |

---

## 7. 健康检查

### 7.1 GET /health

用于确认服务是否可访问。

#### 请求

```http
GET /health
```

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "status": "ok"
  }
}
```

#### curl 示例

```bash
curl http://localhost:8081/health
```

---

## 8. 微信登录

### 8.1 POST /api/auth/wechat-login

使用微信小程序 `wx.login()` 获取的 code 换取后端 JWT。

本地 MVP 联调时，如果服务端未配置微信 AppID 和 AppSecret，后端会使用 `dev_openid_ + code` 生成稳定开发 openid。

#### 请求

```http
POST /api/auth/wechat-login
Content-Type: application/json
```

#### 请求体

```json
{
  "code": "wx_login_code"
}
```

#### 响应 data

```json
{
  "token": "jwt-token",
  "userId": 1,
  "openid": "openid"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "userId": 1,
    "openid": "dev_openid_test_code"
  }
}
```

#### curl 示例

```bash
curl -X POST http://localhost:8081/api/auth/wechat-login \
  -H 'Content-Type: application/json' \
  -d '{"code":"test_code"}'
```

---

## 9. 获取打球记录列表

### 9.1 GET /api/sessions

获取当前登录用户的未删除打球记录。支持通过 `date` 查询某一天的打球记录。

排序规则：

1. 未传 `date` 时：`date` 倒序，`created_at` 倒序
2. 传 `date` 时：`created_at` 倒序

#### 请求

```http
GET /api/sessions
Authorization: Bearer <token>
```

#### Query Parameters

| 参数 | 类型 | 必填 | 示例 | 说明 |
|---|---|---:|---|---|
| `date` | string | 否 | `2026-01-15` | 打球日期，格式 `YYYY-MM-DD` |

#### 按日期查询示例

```http
GET /api/sessions?date=2026-01-15
Authorization: Bearer <token>
```

#### 响应 data

```json
[
  {
    "id": 1,
    "date": "2026-01-15",
    "durationMinutes": 120,
    "rating": 3,
    "type": 5,
    "typeLabel": "双打比赛",
    "category": 3,
    "categoryText": "比赛",
    "subCategory": 2,
    "subCategoryText": "双打",
    "typeText": "比赛 · 双打",
    "matchRank": 1,
    "matchRankLabel": "冠军",
    "courtName": "奥森网球场",
    "cost": 80,
    "racketName": "Wilson Blade",
    "shoeName": "Asics Gel Resolution",
    "note": "今天状态不错",
    "createdAt": "2026-01-15T12:00:00+08:00",
    "updatedAt": "2026-01-15T12:00:00+08:00"
  }
]
```

#### curl 示例

```bash
curl http://localhost:8081/api/sessions \
  -H 'Authorization: Bearer <token>'
```

```bash
curl 'http://localhost:8081/api/sessions?date=2026-01-15' \
  -H 'Authorization: Bearer <token>'
```

---

## 10. 新增打球记录

### 10.1 POST /api/sessions

新增一条打球记录。

#### 请求

```http
POST /api/sessions
Authorization: Bearer <token>
Content-Type: application/json
```

#### 请求体示例：日常球局双打

```json
{
  "date": "2026-01-15",
  "durationMinutes": 120,
  "rating": 4,
  "category": 1,
  "subCategory": 2,
  "matchRank": 0,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 80,
  "racketName": "Wilson Blade",
  "shoeName": "Asics Gel Resolution",
  "note": "今天状态不错"
}
```

#### 请求体示例：双打比赛冠军

```json
{
  "date": "2026-01-15",
  "durationMinutes": 120,
  "rating": 5,
  "category": 3,
  "subCategory": 2,
  "matchRank": 1,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 120,
  "racketName": "Wilson Blade",
  "shoeName": "Asics Gel Resolution",
  "note": "双打比赛冠军"
}
```

#### 响应 data

返回新创建的 `SessionResponse`。

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "date": "2026-01-15",
    "durationMinutes": 120,
    "rating": 5,
    "type": 5,
    "typeLabel": "双打比赛",
    "category": 3,
    "categoryText": "比赛",
    "subCategory": 2,
    "subCategoryText": "双打",
    "typeText": "比赛 · 双打",
    "matchRank": 1,
    "matchRankLabel": "冠军",
    "courtName": "奥森网球场",
    "cost": 120,
    "racketName": "Wilson Blade",
    "shoeName": "Asics Gel Resolution",
    "note": "双打比赛冠军",
    "createdAt": "2026-01-15T12:00:00+08:00",
    "updatedAt": "2026-01-15T12:00:00+08:00"
  }
}
```

#### curl 示例

```bash
curl -X POST http://localhost:8081/api/sessions \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"date":"2026-01-15","durationMinutes":120,"rating":5,"category":3,"subCategory":2,"matchRank":1,"courtName":"奥森网球场","partner":"张三","cost":120,"racketName":"Wilson Blade","shoeName":"Asics Gel Resolution","note":"双打比赛冠军"}'
```

---

## 11. 获取单条打球记录

### 11.1 GET /api/sessions/:id

获取当前登录用户的一条未删除打球记录。

#### 请求

```http
GET /api/sessions/1
Authorization: Bearer <token>
```

#### 响应 data

返回 `SessionResponse`。

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "date": "2026-01-15",
    "durationMinutes": 120,
    "rating": 5,
    "type": 5,
    "typeLabel": "双打比赛",
    "category": 3,
    "categoryText": "比赛",
    "subCategory": 2,
    "subCategoryText": "双打",
    "typeText": "比赛 · 双打",
    "matchRank": 1,
    "matchRankLabel": "冠军",
    "courtName": "奥森网球场",
    "cost": 120,
    "racketName": "Wilson Blade",
    "shoeName": "Asics Gel Resolution",
    "note": "双打比赛冠军",
    "createdAt": "2026-01-15T12:00:00+08:00",
    "updatedAt": "2026-01-15T12:00:00+08:00"
  }
}
```

#### 可能错误

| 场景 | HTTP | code | message |
|---|---:|---:|---|
| id 非法 | 400 | 40001 | invalid request |
| 记录不存在 | 404 | 40401 | not found |
| 未登录 | 401 | 40101 | unauthorized |

---

## 12. 更新打球记录

### 12.1 PUT /api/sessions/:id

更新当前登录用户的一条未删除打球记录。

#### 请求

```http
PUT /api/sessions/1
Authorization: Bearer <token>
Content-Type: application/json
```

#### 请求体

与新增接口一致。

```json
{
  "date": "2026-01-16",
  "durationMinutes": 90,
  "rating": 4,
  "category": 3,
  "subCategory": 1,
  "matchRank": 2,
  "courtName": "国家网球中心",
  "partner": "李四",
  "cost": 100,
  "racketName": "Babolat Pure Drive",
  "shoeName": "Nike Vapor",
  "note": "更新后的记录"
}
```

#### 响应 data

返回更新后的 `SessionResponse`。

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "date": "2026-01-16",
    "durationMinutes": 90,
    "rating": 4,
    "type": 4,
    "typeLabel": "单打比赛",
    "matchRank": 2,
    "matchRankLabel": "亚军",
    "courtName": "国家网球中心",
    "cost": 100,
    "racketName": "Babolat Pure Drive",
    "shoeName": "Nike Vapor",
    "note": "更新后的记录",
    "createdAt": "2026-01-15T12:00:00+08:00",
    "updatedAt": "2026-01-16T10:00:00+08:00"
  }
}
```

#### curl 示例

```bash
curl -X PUT http://localhost:8081/api/sessions/1 \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"date":"2026-01-16","durationMinutes":90,"rating":4,"category":3,"subCategory":1,"matchRank":2,"courtName":"国家网球中心","partner":"李四","cost":100,"racketName":"Babolat Pure Drive","shoeName":"Nike Vapor","note":"更新后的记录"}'
```

---

## 13. 删除打球记录

### 13.1 DELETE /api/sessions/:id

删除当前登录用户的一条打球记录。

当前实现为逻辑删除，不是物理删除。服务端会更新 `deleted_at` 和 `updated_at`。

#### 请求

```http
DELETE /api/sessions/1
Authorization: Bearer <token>
```

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "deleted": true
  }
}
```

#### 可能错误

| 场景 | HTTP | code | message |
|---|---:|---:|---|
| id 非法 | 400 | 40001 | invalid request |
| 记录不存在或已删除 | 404 | 40401 | not found |
| 未登录 | 401 | 40101 | unauthorized |

#### curl 示例

```bash
curl -X DELETE http://localhost:8081/api/sessions/1 \
  -H 'Authorization: Bearer <token>'
```

---

## 14. 获取最近一次打球记录

### 14.1 GET /api/sessions/latest

获取当前登录用户最近一次未删除打球记录。

排序规则：

1. `date` 倒序
2. `created_at` 倒序

#### 请求

```http
GET /api/sessions/latest
Authorization: Bearer <token>
```

#### 有记录时响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "date": "2026-01-16",
    "durationMinutes": 90,
    "rating": 4,
    "type": 4,
    "typeLabel": "单打比赛",
    "matchRank": 2,
    "matchRankLabel": "亚军",
    "courtName": "国家网球中心",
    "cost": 100,
    "racketName": "Babolat Pure Drive",
    "shoeName": "Nike Vapor",
    "note": "更新后的记录",
    "createdAt": "2026-01-15T12:00:00+08:00",
    "updatedAt": "2026-01-16T10:00:00+08:00"
  }
}
```

#### 无记录时响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": null
}
```

---

## 15. 获取日历记录标记

### 15.1 GET /api/sessions/calendar

获取当前登录用户在指定月份或指定年份内有打球记录的日期、统计摘要和图表数据。

#### 请求

按月查询：

```http
GET /api/sessions/calendar?year=2026&month=6
Authorization: Bearer <token>
```

按年查询：

```http
GET /api/sessions/calendar?year=2026
Authorization: Bearer <token>
```

#### Query Parameters

| 参数 | 类型 | 必填 | 示例 | 说明 |
|---|---|---:|---|---|
| `year` | number | 是 | `2026` | 查询年份，范围 `2000-2100` |
| `month` | number | 否 | `6` | 查询月份，范围 `1-12`；不传则查询全年 |

#### 响应说明

- 传 `year + month`：按指定自然月统计，响应中 `month` 为对应月份。
- 只传 `year`：按指定自然年统计，响应中 `month` 为 `0`。
- `days` 返回查询范围内有记录的日期列表。
- `activeDayCount` 表示查询范围内有打球记录的天数，不是打球记录条数。

#### 响应示例：按年查询

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "year": 2026,
    "month": 0,
    "activeDayCount": 2,
    "days": [
      {
        "date": "2026-01-15",
        "count": 1
      },
      {
        "date": "2026-06-10",
        "count": 2
      }
    ],
    "summary": {
      "sessionCount": 3,
      "activeDayCount": 2,
      "totalMinutes": 360,
      "averageMinutes": 120,
      "averageRating": 4,
      "sessionCost": 240,
      "racketCost": 1599,
      "stringingCost": 80,
      "totalCost": 1919,
      "trainingCount": 1,
      "singlesCount": 0,
      "doublesCount": 1,
      "matchCount": 1
    },
    "charts": {
      "weeklySessions": [],
      "ratingTrend": [],
      "expenseBreakdown": [],
      "sessionTypeBreakdown": [
        {
          "key": "training",
          "label": "训练",
          "value": 1,
          "percent": 33.3
        },
        {
          "key": "singles",
          "label": "单打",
          "value": 0,
          "percent": 0
        },
        {
          "key": "doubles",
          "label": "双打",
          "value": 1,
          "percent": 33.3
        },
        {
          "key": "match",
          "label": "比赛",
          "value": 1,
          "percent": 33.3
        }
      ]
    }
  }
}
```

#### 日历统计字段说明

`summary` 字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| sessionCount | number | 查询范围内打球记录数 |
| activeDayCount | number | 查询范围内有打球记录的天数 |
| totalMinutes | number | 查询范围内累计打球分钟数 |
| averageMinutes | number | 查询范围内平均单次打球时长，单位分钟，保留 1 位小数 |
| averageRating | number | 查询范围内平均评分，保留 1 位小数 |
| sessionCost | number | 查询范围内打球消费合计 |
| racketCost | number | 查询范围内球拍购买费用合计 |
| stringingCost | number | 查询范围内穿线费用合计 |
| totalCost | number | `sessionCost + racketCost + stringingCost` |
| trainingCount | number | 查询范围内训练记录数 |
| singlesCount | number | 查询范围内单打记录数 |
| doublesCount | number | 查询范围内双打记录数 |
| matchCount | number | 查询范围内比赛记录数，包含单打比赛和双打比赛 |

`charts.sessionTypeBreakdown` 用于训练/单打/双打/比赛占比展示，固定返回 4 项；无记录时 `value` 和 `percent` 为 `0`。

---

## 16. 获取本月统计

### 16.1 GET /api/stats/month

获取当前登录用户当前自然月统计。

统计口径：

- 使用服务器 `Asia/Shanghai` 时区
- 按 `tennis_sessions.date` 统计
- 不按 `created_at` 统计
- 默认只统计未删除记录

#### 请求

```http
GET /api/stats/month
Authorization: Bearer <token>
```

#### 响应 data

```json
{
  "monthCount": 8,
  "monthMinutes": 960,
  "monthCost": 320,
  "totalCount": 42
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| monthCount | number | 本月打球记录数 |
| monthMinutes | number | 本月总时长，单位分钟 |
| monthCost | number | 本月总花费 |
| totalCount | number | 当前用户总记录数 |

#### 响应示例

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

#### curl 示例

```bash
curl http://localhost:8081/api/stats/month \
  -H 'Authorization: Bearer <token>'
```

### 16.2 GET /api/stats/charts

获取统计图表数据，支持按月或按年查询。

统计口径：

- 使用服务器 `Asia/Shanghai` 时区
- 按 `tennis_sessions.date` 统计
- 默认只统计未删除记录
- `charts.sessionCategoryCountBreakdown` 按一级类型统计打球次数占比，分母为 `summary.sessionCount`
- `charts.sessionSubCategoryCountBreakdown` 按一级类型 + 二级类型统计打球次数占比，分母为 `summary.sessionCount`
- `charts.sessionCategoryDurationBreakdown` 按一级类型统计打球时长占比，分母为 `summary.totalMinutes`
- `charts.sessionSubCategoryDurationBreakdown` 按一级类型 + 二级类型统计打球时长占比，分母为 `summary.totalMinutes`
- `charts.sessionCategoryCostBreakdown` 按一级类型统计打球费用占比，分母为 `summary.sessionCost`
- `charts.sessionSubCategoryCostBreakdown` 按一级类型 + 二级类型统计打球费用占比，分母为 `summary.sessionCost`

#### 请求

```http
GET /api/stats/charts?period=month&year=2026&month=6
Authorization: Bearer <token>
```

| 参数 | 类型 | 必填 | 示例 | 说明 |
|---|---|---|---|---|
| period | string | 是 | `month` | `month` 或 `year` |
| year | number | 是 | `2026` | 查询年份，范围 `2000-2100` |
| month | number | period=month 时是 | `6` | 查询月份，范围 `1-12` |

#### 响应 data 关键字段

| 字段 | 类型 | 说明 |
|---|---|---|
| period | string | 查询周期，`month` 或 `year` |
| year | number | 查询年份 |
| month | number | 查询月份；按年查询时为 `0` |
| rangeText | string | 查询范围展示文案 |
| summary.sessionCost | number | 查询范围内打球费用合计 |
| summary.totalMinutes | number | 查询范围内打球总分钟数 |
| summary.sessionCount | number | 查询范围内打球记录数 |
| charts.expenseBreakdown | array | 打球、球拍、穿线总消费占比 |
| charts.sessionTypeBreakdown | array | 兼容旧统计的训练/单打/双打/比赛次数占比 |
| charts.sessionCategoryCountBreakdown | array | 按一级类型统计打球次数占比，固定返回日常球局、训练、比赛 |
| charts.sessionSubCategoryCountBreakdown | array | 按二级类型统计打球次数占比，固定返回 6 个合法类型组合 |
| charts.sessionCategoryDurationBreakdown | array | 按一级类型统计打球时长占比，固定返回日常球局、训练、比赛 |
| charts.sessionSubCategoryDurationBreakdown | array | 按二级类型统计打球时长占比，固定返回 6 个合法类型组合 |
| charts.sessionCategoryCostBreakdown | array | 按一级类型统计打球费用占比，固定返回日常球局、训练、比赛 |
| charts.sessionSubCategoryCostBreakdown | array | 按二级类型统计打球费用占比，固定返回 6 个合法类型组合 |

#### 类型占比示例

```json
{
  "sessionCategoryCountBreakdown": [
    { "key": "daily", "label": "日常球局", "value": 6, "percent": 60 },
    { "key": "training", "label": "训练", "value": 2, "percent": 20 },
    { "key": "match", "label": "比赛", "value": 2, "percent": 20 }
  ],
  "sessionSubCategoryCountBreakdown": [
    { "key": "daily_singles", "label": "日常球局 · 打单", "value": 3, "percent": 30 },
    { "key": "daily_doubles", "label": "日常球局 · 双打", "value": 3, "percent": 30 },
    { "key": "training_serve", "label": "训练 · 发球", "value": 1, "percent": 10 },
    { "key": "training_other", "label": "训练 · 其他", "value": 1, "percent": 10 },
    { "key": "match_singles", "label": "比赛 · 单打", "value": 1, "percent": 10 },
    { "key": "match_doubles", "label": "比赛 · 双打", "value": 1, "percent": 10 }
  ],
  "sessionCategoryDurationBreakdown": [
    { "key": "daily", "label": "日常球局", "value": 720, "percent": 60 },
    { "key": "training", "label": "训练", "value": 240, "percent": 20 },
    { "key": "match", "label": "比赛", "value": 240, "percent": 20 }
  ],
  "sessionSubCategoryDurationBreakdown": [
    { "key": "daily_singles", "label": "日常球局 · 打单", "value": 360, "percent": 30 },
    { "key": "daily_doubles", "label": "日常球局 · 双打", "value": 360, "percent": 30 },
    { "key": "training_serve", "label": "训练 · 发球", "value": 120, "percent": 10 },
    { "key": "training_other", "label": "训练 · 其他", "value": 120, "percent": 10 },
    { "key": "match_singles", "label": "比赛 · 单打", "value": 120, "percent": 10 },
    { "key": "match_doubles", "label": "比赛 · 双打", "value": 120, "percent": 10 }
  ],
  "sessionCategoryCostBreakdown": [
    { "key": "daily", "label": "日常球局", "value": 180, "percent": 60 },
    { "key": "training", "label": "训练", "value": 60, "percent": 20 },
    { "key": "match", "label": "比赛", "value": 60, "percent": 20 }
  ],
  "sessionSubCategoryCostBreakdown": [
    { "key": "daily_singles", "label": "日常球局 · 打单", "value": 100, "percent": 33.3 },
    { "key": "daily_doubles", "label": "日常球局 · 双打", "value": 80, "percent": 26.7 },
    { "key": "training_serve", "label": "训练 · 发球", "value": 60, "percent": 20 },
    { "key": "training_other", "label": "训练 · 其他", "value": 0, "percent": 0 },
    { "key": "match_singles", "label": "比赛 · 单打", "value": 60, "percent": 20 },
    { "key": "match_doubles", "label": "比赛 · 双打", "value": 0, "percent": 0 }
  ]
}
```

#### curl 示例

```bash
curl 'http://localhost:8081/api/stats/charts?period=month&year=2026&month=6' \
  -H 'Authorization: Bearer <token>'
```

---

## 17. 前端联调建议

### 17.1 小程序开发者工具

本地 HTTP 调试需要在微信开发者工具中开启：

```text
详情 -> 本地设置 -> 不校验合法域名、web-view、TLS 版本以及 HTTPS 证书
```

### 17.2 真机调试

真机不能访问电脑的 `localhost`。

需要：

1. 手机和电脑在同一局域网
2. 后端监听 `0.0.0.0:PORT`
3. 小程序请求地址使用电脑局域网 IP

例如：

```text
http://192.168.1.8:8081
```

### 17.3 登录后保存 token

前端登录成功后需要保存：

```ts
wx.setStorageSync('token', loginResp.token)
```

后续请求携带：

```http
Authorization: Bearer <token>
```

---

## 18. 完整联调流程示例

### 18.1 获取 token

```bash
TOKEN=$(curl -s -X POST http://localhost:8081/api/auth/wechat-login \
  -H 'Content-Type: application/json' \
  -d '{"code":"test_code"}' | jq -r '.data.token')
```

### 18.2 新增记录

```bash
curl -X POST http://localhost:8081/api/sessions \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"date":"2026-01-15","durationMinutes":120,"rating":5,"category":3,"subCategory":2,"matchRank":1,"courtName":"奥森网球场","partner":"张三","cost":120,"racketName":"Wilson Blade","shoeName":"Asics Gel Resolution","note":"双打比赛冠军"}'
```

### 18.3 查看列表

```bash
curl http://localhost:8081/api/sessions \
  -H "Authorization: Bearer $TOKEN"
```

### 18.4 查看本月统计

```bash
curl http://localhost:8081/api/stats/month \
  -H "Authorization: Bearer $TOKEN"
```

---

## 19. 注意事项

1. `date` 字段必须使用 `YYYY-MM-DD`。
2. `createdAt` 和 `updatedAt` 为服务端时间。
3. 删除接口是软删除，普通列表和统计不会返回已删除数据。
4. 普通查询只能访问当前 token 对应用户的数据。
5. 当前列表接口支持分页和按日期筛选。
6. 新客户端优先使用 `category`、`subCategory`、`typeText`；`type` 和 `typeLabel` 为兼容旧客户端保留。

---

## 20. 球拍管理接口

### 20.1 球拍状态

| 值 | 含义 |
|---:|---|
| 1 | 主力拍 |
| 2 | 在用 |
| 3 | 已退役 |

### 20.2 RacketResponse

```json
{
  "id": 1,
  "name": "EZONE 主力拍",
  "brand": "Yonex",
  "model": "EZONE 100",
  "status": 1,
  "imageUrl": "",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80,
  "usageCount": 12,
  "usageMinutes": 1440,
  "usageHours": 24,
  "totalMinutes": 1440,
  "totalHours": 24,
  "createdAt": "2026-05-30T12:00:00+08:00",
  "updatedAt": "2026-05-30T12:00:00+08:00"
}
```

说明：

- `stringName`、`tension`、`lastStringDate`、`lastStringCost` 来自最近一条穿线记录。
- `usageCount`、`usageMinutes`、`usageHours` 通过打球记录中的 `racketId` 实时统计。
- `totalMinutes`、`totalHours` 为兼容旧前端保留，当前与 `usageMinutes`、`usageHours` 一致。
- 默认列表不返回已退役球拍。

### 20.3 获取球拍列表

```http
GET /api/rackets
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| includeRetired | boolean | 是否包含已退役球拍，默认 `false` |

响应 data：

```json
[
  {
    "id": 1,
    "name": "EZONE 主力拍",
    "brand": "Yonex",
    "model": "EZONE 100",
    "status": 1,
    "stringName": "Poly Tour Pro",
    "tension": 48,
    "lastStringDate": "2026-05-10",
    "lastStringCost": 80,
    "usageCount": 12,
    "usageMinutes": 1440,
    "usageHours": 24,
    "totalMinutes": 1440,
    "totalHours": 24
  }
]
```

统计口径：

- 按 `tennis_sessions.racket_id = racket.id` 关联统计。
- 只统计当前登录用户自己的未删除打球记录。
- 统计 `COUNT(tennis_sessions.id)` 得到 `usageCount`。
- 统计 `SUM(tennis_sessions.duration_minutes)` 得到 `usageMinutes`。
- `usageHours = usageMinutes / 60`，当前向下取整。
- 列表返回哪些球拍仍由 `includeRetired` 和球拍删除状态控制。

新增字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| usageCount | number | 该球拍关联打球记录次数 |
| usageMinutes | number | 该球拍累计使用分钟数 |
| usageHours | number | 该球拍累计使用小时数，当前向下取整 |
| totalMinutes | number | 兼容旧字段，等于 `usageMinutes` |
| totalHours | number | 兼容旧字段，等于 `usageHours` |

### 20.4 获取球拍统计

统计当前登录用户的球拍数量和球拍相关消费。

```http
GET /api/rackets/stats
Authorization: Bearer <token>
```

统计口径：

- `racketCount`：当前用户未删除球拍总数，包含主力拍、在用、已退役，不包含已逻辑删除球拍。
- `racketCost`：当前用户未删除球拍的购买费用合计，累加 `racket.purchase_price`，空值按 `0` 处理。
- `stringingCost`：当前用户未删除球拍的未删除穿线记录费用合计，累加 `racket_stringing_record.cost`。
- 已删除球拍不计入球拍总数、球拍花费和穿线费用。
- 已删除穿线记录不计入穿线费用。

响应 data：

```json
{
  "racketCount": 3,
  "racketCost": 4597,
  "stringingCost": 320,
  "totalCost": 4917,
  "racketCostText": "4597.00",
  "stringingCostText": "320.00",
  "totalCostText": "4917.00"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| racketCount | number | 球拍总数，排除已删除球拍 |
| racketCost | number | 球拍购买费用合计 |
| stringingCost | number | 球拍穿线费用合计 |
| totalCost | number | 球拍购买费用 + 穿线费用 |
| racketCostText | string | 格式化后的球拍购买费用，保留两位小数 |
| stringingCostText | string | 格式化后的穿线费用，保留两位小数 |
| totalCostText | string | 格式化后的总费用，保留两位小数 |

curl 示例：

```bash
curl http://localhost:8081/api/rackets/stats \
  -H 'Authorization: Bearer <token>'
```

### 20.5 获取可选球拍列表

新增打球记录时使用，只返回：

```text
status IN (1, 2)
```

接口：

```http
GET /api/rackets/selectable
Authorization: Bearer <token>
```

### 20.6 添加球拍

```http
POST /api/rackets
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "name": "EZONE 主力拍",
  "brand": "Yonex",
  "model": "EZONE 100",
  "imageUrl": "",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80
}
```

规则：

- `name` 必填。
- 新增球拍默认 `status = 2`，表示在用。
- 如果传入穿线相关字段，会自动创建一条穿线记录。

### 20.7 获取球拍详情

```http
GET /api/rackets/:id
Authorization: Bearer <token>
```

响应 data：

```json
{
  "racket": {
    "id": 1,
    "libraryId": 1,
    "name": "EZONE 主力拍",
    "brand": "Yonex",
    "model": "EZONE 100",
    "status": 1,
    "imageUrl": "",
    "purchaseDate": "2026-01-01",
    "purchasePrice": 1599,
    "stringName": "Poly Tour Pro",
    "tension": 48,
    "lastStringDate": "2026-05-10",
    "lastStringCost": 80,
    "usageCount": 12,
    "usageMinutes": 1440,
    "usageHours": 24,
    "afterStringingUsageCount": 3,
    "afterStringingUsageMinutes": 360,
    "afterStringingUsageHours": 6,
    "totalMinutes": 1440,
    "totalHours": 24,
    "createdAt": "2026-05-30T12:00:00+08:00",
    "updatedAt": "2026-05-30T12:00:00+08:00"
  },
  "stringingRecords": [
    {
      "id": 1,
      "racketId": 1,
      "stringName": "Poly Tour Pro",
      "tension": 48,
      "cost": 80,
      "stringDate": "2026-05-10",
      "createdAt": "2026-05-30T12:00:00+08:00",
      "updatedAt": "2026-05-30T12:00:00+08:00"
    }
  ]
}
```

详情页使用统计口径：

- 累计使用统计：当前用户未删除打球记录中，按 `tennis_sessions.racket_id = racket.id` 统计。
- `usageCount = COUNT(tennis_sessions.id)`。
- `usageMinutes = COALESCE(SUM(tennis_sessions.duration_minutes), 0)`。
- `usageHours = usageMinutes / 60`，当前向下取整。
- 最近一次穿线按 `string_date DESC, id DESC` 取第一条未删除穿线记录。
- 如果存在最近一次穿线记录，穿线后使用统计按 `tennis_sessions.date >= lastStringDate` 统计，包含穿线当天。
- 如果没有穿线记录，`afterStringingUsageCount`、`afterStringingUsageMinutes`、`afterStringingUsageHours` 均为 `0`。

详情新增字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| usageCount | number | 该球拍累计关联打球记录次数 |
| usageMinutes | number | 该球拍累计使用分钟数 |
| usageHours | number | 该球拍累计使用小时数，当前向下取整 |
| afterStringingUsageCount | number | 最近一次穿线后累计使用次数 |
| afterStringingUsageMinutes | number | 最近一次穿线后累计使用分钟数 |
| afterStringingUsageHours | number | 最近一次穿线后累计使用小时数，当前向下取整 |

### 20.8 编辑球拍

```http
PUT /api/rackets/:id
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "name": "EZONE 主力拍",
  "brand": "Yonex",
  "model": "EZONE 100",
  "status": 1,
  "imageUrl": "",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80
}
```

说明：

- `status = 1` 时会自动将其他主力拍改为在用。
- 传入穿线相关字段时，会新增一条穿线记录。

### 20.9 设置主力拍

```http
POST /api/rackets/:id/set-primary
Authorization: Bearer <token>
```

规则：

- 系统只允许当前用户存在一支主力拍。
- 设置成功后，原主力拍会自动变成在用。

### 20.10 退役球拍

```http
POST /api/rackets/:id/retire
Authorization: Bearer <token>
```

规则：

- 退役后 `status = 3`。
- 默认球拍列表不展示。
- `/api/rackets/selectable` 不返回退役球拍。
- 历史打球记录和统计不受影响。

### 20.11 新增穿线记录

```http
POST /api/rackets/:id/stringing-records
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "cost": 80,
  "stringDate": "2026-05-10"
}
```

### 20.12 打球记录关联球拍

新增/编辑打球记录支持传入：

```json
{
  "racketId": 1,
  "racketName": "EZONE 主力拍"
}
```

`racketId` 用于球拍累计使用时长统计，`racketName` 用于前端展示兼容。

---

## 21. 球拍库与添加球拍选择接口

### 21.1 获取球拍库，按品牌分类

添加球拍页面使用。返回系统球拍库中的球拍，并按品牌分组。

```http
GET /api/racket-library
Authorization: Bearer <token>
```

响应 data 示例：

```json
[
  {
    "brand": "Yonex",
    "items": [
      {
        "id": 1,
        "brand": "Yonex",
        "model": "EZONE 100",
        "releaseYear": 2025,
        "weight": 300,
        "headSize": 100,
        "imageUrl": ""
      }
    ]
  },
  {
    "brand": "Wilson",
    "items": [
      {
        "id": 7,
        "brand": "Wilson",
        "model": "Blade 98 16x19",
        "releaseYear": 2024,
        "weight": 305,
        "headSize": 98,
        "imageUrl": ""
      }
    ]
  }
]
```

### 21.2 从球拍库添加到我的球拍

用户选择球拍库中的球拍后，新增我的球拍时传 `libraryId`。

```http
POST /api/rackets
Authorization: Bearer <token>
Content-Type: application/json
```

请求体示例：

```json
{
  "libraryId": 1,
  "name": "EZONE 主力拍",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80
}
```

说明：

- `libraryId` 有值时，后端会从球拍库补全 `brand`、`model`、`imageUrl`。
- 如果请求体里也传了 `brand`、`model`、`imageUrl`，以前端传入值为准。
- `name` 仍可自定义，比如“EZONE 主力拍”。

### 21.3 添加库中没有的球拍

如果球拍库没有对应球拍，用户可以手动输入。

```http
POST /api/rackets
Authorization: Bearer <token>
Content-Type: application/json
```

请求体示例：

```json
{
  "name": "我的老款备用拍",
  "brand": "Volkl",
  "model": "V-Cell 10",
  "purchasePrice": 1200
}
```

说明：

- 手动输入时不传 `libraryId`，或传 `0`。
- `name` 必填。
- `brand`、`model` 可选。

### 21.4 获取我的球拍，用于新增打球记录时选择

新增打球记录页面使用。只返回当前用户未退役球拍：

```text
status IN (1, 2)
```

接口：

```http
GET /api/my-rackets
Authorization: Bearer <token>
```

响应 data 示例：

```json
[
  {
    "id": 1,
    "libraryId": 1,
    "name": "EZONE 主力拍",
    "brand": "Yonex",
    "model": "EZONE 100",
    "status": 1,
    "stringName": "Poly Tour Pro",
    "tension": 48,
    "lastStringDate": "2026-05-10",
    "totalHours": 86
  }
]
```

新增打球记录时，把选中的我的球拍写入：

```json
{
  "racketId": 1,
  "racketName": "EZONE 主力拍"
}
```

---

## 22. 球拍接口重构说明

### 22.1 添加球拍不再包含穿线信息

添加球拍只维护球拍本体信息，不创建穿线记录。

```http
POST /api/rackets
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "libraryId": 1,
  "name": "EZONE 主力拍",
  "brand": "Yonex",
  "model": "EZONE 100",
  "imageUrl": "",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599
}
```

说明：

- `name` 必填。
- `libraryId` 可选。
- 从球拍库选择时，后端可根据 `libraryId` 补全 `brand`、`model`、`imageUrl`。
- 库中没有的球拍，用户可不传 `libraryId`，直接手动输入 `name`、`brand`、`model`。
- 当前接口不再接收/处理 `stringName`、`tension`、`lastStringDate`、`lastStringCost`。

### 22.2 编辑球拍，支持主力/在用/退役

```http
PUT /api/rackets/:id
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "libraryId": 1,
  "name": "EZONE 主力拍",
  "brand": "Yonex",
  "model": "EZONE 100",
  "status": 1,
  "imageUrl": "",
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599
}
```

说明：

- `status = 1`：设为主力拍，后端会自动将当前用户其他主力拍改为在用。
- `status = 2`：在用。
- `status = 3`：退役。
- 编辑球拍不再接收/处理穿线信息。

### 22.3 删除球拍

```http
DELETE /api/rackets/:id
Authorization: Bearer <token>
```

响应：

```json
{
  "deleted": true
}
```

说明：

- 当前实现为逻辑删除，更新 `deleted_at`。
- 删除后默认列表、我的球拍列表不再返回。
- 历史打球记录中的 `racketId` 不会被清空。

### 22.4 新增穿线记录

穿线信息通过独立接口维护。

```http
POST /api/rackets/:id/stringing-records
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "stringName": "Poly Tour Pro",
  "tension": 48,
  "cost": 80,
  "stringDate": "2026-05-10"
}
```

说明：

- `stringName` 必填。
- `stringDate` 必填，格式 `YYYY-MM-DD`。
- 球拍列表和详情中的当前球线信息来自最近一条穿线记录。

### 22.5 保留我的球拍接口

新增打球记录选择球拍时继续使用：

```http
GET /api/my-rackets
Authorization: Bearer <token>
```

只返回未删除且未退役的球拍：

```text
status IN (1, 2)
```
