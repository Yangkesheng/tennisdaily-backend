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
| 3 | 季军 | 季军 |
| 4 | 四强 | 四强 |
| 5 | 八强 | 八强 |
| 6 | 16强 | 16强 |
| 7 | 小组赛 | 小组赛 |

说明：

- 当 `category` 为 `3` 时，`matchRank` 可以为 `1-7`
- 当 `category` 为 `1` 或 `2` 时，后端会强制将 `matchRank` 处理为 `0`
- 旧客户端仍可通过 `type = 4` 或 `type = 5` 提交比赛成绩

---

## 4. 数据模型

### 4.1 SessionResponse

打球记录响应结构：

```json
{
  "id": 1,
  "date": "2026-01-15 19:30",
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
  "shoeId": 1,
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
| date | string | 打球开始时间，格式 `YYYY-MM-DD HH:mm` |
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
| racketName | string | 球拍名称，按 racketId 关联我的球拍实时返回 |
| shoeId | number | 使用球鞋 ID，未选择时为 `0` |
| shoeName | string | 球鞋名称：`shoeId > 0` 时按鞋 ID 关联我的球鞋实时返回；无 `shoeId` 时返回创建时快照；球鞋已删除时为空字符串 |
| note | string | 备注 |
| createdAt | string | 创建时间，RFC3339 |
| updatedAt | string | 更新时间，RFC3339 |

### 4.2 CreateSessionRequest / UpdateSessionRequest

新增和更新记录请求结构一致：

```json
{
  "date": "2026-01-15 19:30",
  "durationMinutes": 120,
  "rating": 3,
  "category": 3,
  "subCategory": 2,
  "matchRank": 1,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 80,
  "shoeId": 1,
  "shoeName": "Asics Gel Resolution",
  "note": "今天状态不错"
}
```

字段规则：

| 字段 | 是否必填 | 规则 |
|---|---|---|
| date | 是 | 格式 `YYYY-MM-DD HH:mm`，旧格式 `YYYY-MM-DD` 仍兼容并按当天 `00:00` 处理 |
| durationMinutes | 否 | 默认 120；必须大于 0，最大 600 |
| rating | 否 | 默认 3；范围 1-5 |
| category | 新客户端必填 | 一级类型，允许 `1`、`2`、`3` |
| subCategory | 新客户端必填 | 二级类型，允许 `1`、`2`、`3`、`4`；必须属于当前 `category` |
| type | 旧客户端必填 | 旧类型枚举，允许 1-5；未传新字段时后端按旧类型映射 |
| matchRank | 否 | 允许 0-7；非比赛类型会被强制改为 0 |
| courtName | 否 | 字符串 |
| partner | 否 | 字符串，搭档名称 |
| cost | 否 | 数字 |
| shoeId | 否 | 数字，我的球鞋 ID；不传或传 `0` 表示不关联球鞋 |
| shoeName | 否 | 字符串，兼容旧客户端保留；新客户端传 `shoeId` 即可 |
| note | 否 | 字符串 |

注意：

- 新客户端应提交 `category` 和 `subCategory`，后端会同步生成兼容旧字段 `type`
- 旧客户端仍可只提交 `type`，后端会自动映射出 `category` 和 `subCategory`
- 新增/编辑不再接收 `racketName`，响应中的 `racketName` 由后端按 `racketId` 关联我的球拍实时返回
- 响应中的 `shoeName` 优先按 `shoeId` 关联我的球鞋实时返回；`shoeId` 为 `0` 时回退到请求里的 `shoeName` 快照；关联球鞋已删除时返回空字符串
- `date` 必须是 `YYYY-MM-DD HH:mm`，旧格式 `YYYY-MM-DD` 仍兼容并按当天 `00:00` 处理，不要传完整 ISO 时间
- `durationMinutes` 传 `0` 时后端会使用默认值 `120`
- `rating` 传 `0` 时后端会使用默认值 `3`

---

## 5. 接口列表总览

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| GET | `/health` | 否 | 健康检查 |
| POST | `/api/auth/wechat-login` | 否 | 微信登录 |
| POST | `/api/auth/phone-login` | 否 | 手机号登录 |
| GET | `/api/session-config` | 否 | 获取打球记录配置 |
| GET | `/api/sessions` | 是 | 获取打球记录列表 |
| POST | `/api/sessions` | 是 | 新增打球记录 |
| GET | `/api/sessions/latest` | 是 | 获取最近一次打球记录 |
| GET | `/api/sessions/calendar` | 是 | 获取日历记录标记，支持按月或按年 |
| GET | `/api/sessions/:id` | 是 | 获取单条打球记录 |
| PUT | `/api/sessions/:id` | 是 | 更新打球记录 |
| DELETE | `/api/sessions/:id` | 是 | 删除打球记录，软删除 |
| GET | `/api/stats/month` | 是 | 获取本月统计 |
| GET | `/api/rackets/stats` | 是 | 获取球拍统计 |
| GET | `/api/shoe-brands` | 是 | 获取球鞋品牌 |
| GET | `/api/shoe-series` | 是 | 获取球鞋系列，支持按品牌、性别过滤 |
| GET | `/api/shoe-library` | 是 | 获取球鞋库，按品牌分类 |
| GET | `/api/shoe-library/stats` | 是 | 获取球鞋库品牌/系列数量统计 |
| GET | `/api/shoes` | 是 | 获取我的球鞋列表 |
| POST | `/api/shoes` | 是 | 新增我的球鞋 |
| GET | `/api/shoes/:id` | 是 | 获取球鞋详情 |
| PUT | `/api/shoes/:id` | 是 | 更新球鞋 |
| DELETE | `/api/shoes/:id` | 是 | 删除球鞋，软删除 |
| GET | `/api/shoes/stats` | 是 | 获取球鞋统计 |
| GET | `/api/my-shoes` | 是 | 获取可选球鞋，新增打球记录选择用 |
| GET | `/api/admin/permissions` | 是 | 查询当前用户是否为管理员 |
| POST | `/api/admin/shoe-brands` | 是（管理员） | 新增球鞋品牌 |
| POST | `/api/admin/shoe-series` | 是（管理员） | 新增球鞋系列 |
| POST | `/api/admin/shoe-library` | 是（管理员） | 新增球鞋库鞋款，配色合并为单条 |

---

## 6. 打球记录配置

### 6.1 GET /api/session-config

获取打球记录表单和展示所需的分类配置、比赛成绩、默认值等。

#### 请求

```http
GET /api/session-config
```

#### 响应 data

```json
{
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
  "matchRanks": [
    { "value": 0, "label": "" },
    { "value": 1, "label": "冠军" },
    { "value": 2, "label": "亚军" },
    { "value": 3, "label": "季军" },
    { "value": 4, "label": "四强" },
    { "value": 5, "label": "八强" },
    { "value": 6, "label": "16强" },
    { "value": 7, "label": "小组赛" }
  ]
}
```

#### 字段说明

| 字段 | 说明 |
|---|---|
| `categories` | 打球记录两级分类，前端可直接用于一级/二级联动选择 |
| `subCategories[].legacyType` | 当前二级分类对应的旧 `type`，用于旧客户端兼容展示 |
| `matchRanks` | 比赛成绩枚举 |

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
| `date` | string | 否 | `2026-01-15` | 按自然日筛选，格式 `YYYY-MM-DD` |

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
    "date": "2026-01-15 19:30",
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
  "date": "2026-01-15 19:30",
  "durationMinutes": 120,
  "rating": 4,
  "category": 1,
  "subCategory": 2,
  "matchRank": 0,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 80,
  "shoeName": "Asics Gel Resolution",
  "note": "今天状态不错"
}
```

#### 请求体示例：双打比赛冠军

```json
{
  "date": "2026-01-15 19:30",
  "durationMinutes": 120,
  "rating": 5,
  "category": 3,
  "subCategory": 2,
  "matchRank": 1,
  "courtName": "奥森网球场",
  "partner": "张三",
  "cost": 120,
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
    "date": "2026-01-15 19:30",
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
  -d '{"date":"2026-01-15","durationMinutes":120,"rating":5,"category":3,"subCategory":2,"matchRank":1,"courtName":"奥森网球场","partner":"张三","cost":120,"shoeId":1,"shoeName":"Asics Gel Resolution","note":"双打比赛冠军"}'
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
    "date": "2026-01-15 19:30",
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
  "date": "2026-01-16 19:30",
  "durationMinutes": 90,
  "rating": 4,
  "category": 3,
  "subCategory": 1,
  "matchRank": 2,
  "courtName": "国家网球中心",
  "partner": "李四",
  "cost": 100,
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
    "date": "2026-01-16 19:30",
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
  -d '{"date":"2026-01-16","durationMinutes":90,"rating":4,"category":3,"subCategory":1,"matchRank":2,"courtName":"国家网球中心","partner":"李四","cost":100,"shoeId":2,"shoeName":"Nike Vapor","note":"更新后的记录"}'
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
    "date": "2026-01-16 19:30",
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
      "shoeCost": 890,
      "totalCost": 2809,
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
| shoeCost | number | 查询范围内球鞋购买费用合计 |
| totalCost | number | `sessionCost + racketCost + stringingCost + shoeCost` |
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
- `charts.sessionCategoryCountBreakdown` 按配置中的一级类型统计打球次数占比，分母为 `summary.sessionCount`
- `charts.sessionSubCategoryCountBreakdown` 按配置中的一级类型 + 二级类型统计打球次数占比，分母为 `summary.sessionCount`
- `charts.sessionCategoryDurationBreakdown` 按配置中的一级类型统计打球时长占比，分母为 `summary.totalMinutes`
- `charts.sessionSubCategoryDurationBreakdown` 按配置中的一级类型 + 二级类型统计打球时长占比，分母为 `summary.totalMinutes`
- `charts.sessionCategoryCostBreakdown` 按配置中的一级类型统计打球费用占比，分母为 `summary.sessionCost`
- `charts.sessionSubCategoryCostBreakdown` 按配置中的一级类型 + 二级类型统计打球费用占比，分母为 `summary.sessionCost`

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
| summary.racketCost | number | 查询范围内球拍购买费用合计 |
| summary.stringingCost | number | 查询范围内穿线费用合计 |
| summary.shoeCost | number | 查询范围内球鞋购买费用合计 |
| summary.totalCost | number | 打球 + 球拍 + 穿线 + 球鞋费用合计 |
| summary.totalMinutes | number | 查询范围内打球总分钟数 |
| summary.sessionCount | number | 查询范围内打球记录数 |
| charts.expenseBreakdown | array | 打球、球拍、穿线、球鞋总消费占比 |
| charts.sessionCategoryCountBreakdown | array | 按配置中的一级类型统计打球次数占比 |
| charts.sessionSubCategoryCountBreakdown | array | 按配置中的二级类型统计打球次数占比 |
| charts.sessionCategoryDurationBreakdown | array | 按配置中的一级类型统计打球时长占比 |
| charts.sessionSubCategoryDurationBreakdown | array | 按配置中的二级类型统计打球时长占比 |
| charts.sessionCategoryCostBreakdown | array | 按配置中的一级类型统计打球费用占比 |
| charts.sessionSubCategoryCostBreakdown | array | 按配置中的二级类型统计打球费用占比 |

#### 类型占比示例

```json
{
  "sessionCategoryCountBreakdown": [
    { "key": "category_1", "category": 1, "label": "日常球局", "value": 6, "percent": 60 },
    { "key": "category_2", "category": 2, "label": "训练", "value": 2, "percent": 20 },
    { "key": "category_3", "category": 3, "label": "比赛", "value": 2, "percent": 20 }
  ],
  "sessionSubCategoryCountBreakdown": [
    { "key": "category_1_sub_1", "category": 1, "subCategory": 1, "label": "日常球局 · 打单", "value": 3, "percent": 30 },
    { "key": "category_1_sub_2", "category": 1, "subCategory": 2, "label": "日常球局 · 双打", "value": 3, "percent": 30 },
    { "key": "category_2_sub_3", "category": 2, "subCategory": 3, "label": "训练 · 发球", "value": 1, "percent": 10 },
    { "key": "category_2_sub_4", "category": 2, "subCategory": 4, "label": "训练 · 其他", "value": 1, "percent": 10 },
    { "key": "category_3_sub_1", "category": 3, "subCategory": 1, "label": "比赛 · 单打", "value": 1, "percent": 10 },
    { "key": "category_3_sub_2", "category": 3, "subCategory": 2, "label": "比赛 · 双打", "value": 1, "percent": 10 }
  ],
  "sessionCategoryDurationBreakdown": [
    { "key": "category_1", "category": 1, "label": "日常球局", "value": 720, "percent": 60 },
    { "key": "category_2", "category": 2, "label": "训练", "value": 240, "percent": 20 },
    { "key": "category_3", "category": 3, "label": "比赛", "value": 240, "percent": 20 }
  ],
  "sessionSubCategoryDurationBreakdown": [
    { "key": "category_1_sub_1", "category": 1, "subCategory": 1, "label": "日常球局 · 打单", "value": 360, "percent": 30 },
    { "key": "category_1_sub_2", "category": 1, "subCategory": 2, "label": "日常球局 · 双打", "value": 360, "percent": 30 },
    { "key": "category_2_sub_3", "category": 2, "subCategory": 3, "label": "训练 · 发球", "value": 120, "percent": 10 },
    { "key": "category_2_sub_4", "category": 2, "subCategory": 4, "label": "训练 · 其他", "value": 120, "percent": 10 },
    { "key": "category_3_sub_1", "category": 3, "subCategory": 1, "label": "比赛 · 单打", "value": 120, "percent": 10 },
    { "key": "category_3_sub_2", "category": 3, "subCategory": 2, "label": "比赛 · 双打", "value": 120, "percent": 10 }
  ],
  "sessionCategoryCostBreakdown": [
    { "key": "category_1", "category": 1, "label": "日常球局", "value": 180, "percent": 60 },
    { "key": "category_2", "category": 2, "label": "训练", "value": 60, "percent": 20 },
    { "key": "category_3", "category": 3, "label": "比赛", "value": 60, "percent": 20 }
  ],
  "sessionSubCategoryCostBreakdown": [
    { "key": "category_1_sub_1", "category": 1, "subCategory": 1, "label": "日常球局 · 打单", "value": 100, "percent": 33.3 },
    { "key": "category_1_sub_2", "category": 1, "subCategory": 2, "label": "日常球局 · 双打", "value": 80, "percent": 26.7 },
    { "key": "category_2_sub_3", "category": 2, "subCategory": 3, "label": "训练 · 发球", "value": 60, "percent": 20 },
    { "key": "category_2_sub_4", "category": 2, "subCategory": 4, "label": "训练 · 其他", "value": 0, "percent": 0 },
    { "key": "category_3_sub_1", "category": 3, "subCategory": 1, "label": "比赛 · 单打", "value": 60, "percent": 20 },
    { "key": "category_3_sub_2", "category": 3, "subCategory": 2, "label": "比赛 · 双打", "value": 0, "percent": 0 }
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
  -d '{"date": "2026-01-15 19:30","durationMinutes":120,"rating":5,"category":3,"subCategory":2,"matchRank":1,"courtName":"奥森网球场","partner":"张三","cost":120,"shoeId":1,"shoeName":"Asics Gel Resolution","note":"双打比赛冠军"}'
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

1. 新增/编辑记录的 `date` 字段推荐使用 `YYYY-MM-DD HH:mm`；旧格式 `YYYY-MM-DD` 仍兼容并按当天 `00:00` 处理。
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
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80,
  "latestStringingRecord": {
    "id": 1,
    "racketId": 1,
    "stringName": "Poly Tour Pro",
    "storeName": "冠军穿线工作室",
    "verticalTension": 48,
    "horizontalTension": 46,
    "cost": 80,
    "stringDate": "2026-05-10 10:00",
    "createdAt": "2026-05-10T10:00:00+08:00",
    "updatedAt": "2026-05-10T10:00:00+08:00"
  },
  "afterStringingUsageCount": 5,
  "afterStringingUsageMinutes": 480,
  "afterStringingUsageHours": 8,
  "stringHealth": {
    "state": "good",
    "display": "状态良好 · 预计还可打 7h",
    "score": 46.625,
    "remainingHours": 7
  },
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

- `stringName`、`storeName`、`verticalTension`、`horizontalTension`、`lastStringDate`、`lastStringCost` 来自最近一条穿线记录。
- `latestStringingRecord` 为最近一次完整穿线记录；没有穿线记录时返回 `null`。
- `afterStringingUsageCount`、`afterStringingUsageMinutes`、`afterStringingUsageHours` 为最近一次穿线之后（打球记录 `date >= string_date`）的累计使用统计；没有穿线记录时为 `0`。
- `stringHealth` 为聚酯线衰减健康度，仅在存在穿线记录时返回对象，没有穿线记录时返回 `null`。
- `usageCount`、`usageMinutes`、`usageHours` 通过打球记录中的 `racketId` 实时统计。
- `totalMinutes`、`totalHours` 为兼容旧前端保留，当前与 `usageMinutes`、`usageHours` 一致。
- 默认列表不返回已退役球拍。

`stringHealth` 对象字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| state | string | 状态标识，取值为 `fresh` / `peak` / `good` / `decline` / `dead` / `expired` |
| display | string | 前端直接展示的状态文案 |
| score | number | 健康分，范围 `0-100`，可作为进度条百分比 |
| remainingHours | number | 预计剩余可打小时数，四舍五入为整数 |

健康度计算公式（当前版本统一按聚酯线，不区分线径、打法、气温、湿度等修正系数）：

```text
effective_wear = afterStringingUsageMinutes / 60 + daysSinceStringing * 0.18
score = max(0, 100 * (1 - effective_wear / 16))
remainingHours = max(0, 16 - effective_wear)
```

其中 `daysSinceStringing` 为最近一次穿线日期到当前日期的自然天数差。`0.18`（每日静置衰减系数）、`16`（标准可用寿命）以及下面状态分档的阈值和展示文案均可在 `config.yaml` 的 `polyesterStringHealth` 段配置调整，上式为默认值。

状态分档（对应 `config.yaml` 默认配置）：

| state | score 区间 | display |
|---|---|---|
| fresh | `>= 85` | `新上线 · 手感正脆` |
| peak | `70 - 84` | `巅峰期 · 预计还可打 {remainingHours}h` |
| good | `45 - 69` | `状态良好 · 预计还可打 {remainingHours}h` |
| decline | `25 - 44` | `开始衰减 · 建议近期重穿` |
| dead | `10 - 24` | `手感变死 · 建议重穿` |
| expired | `< 10` | `已超期 · 不建议比赛使用` |

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
    "libraryId": 1,
    "name": "EZONE 主力拍",
    "brand": "Yonex",
    "model": "EZONE 100",
    "status": 1,
    "releaseYear": 2025,
    "weight": 300,
    "headSize": 100,
    "stringPattern": "16x19",
    "fileId": "cloud://tennisdaily/rackets/ezone100.png",
    "stringName": "Poly Tour Pro",
    "storeName": "冠军穿线工作室",
    "verticalTension": 48,
    "horizontalTension": 46,
    "lastStringDate": "2026-05-10",
    "lastStringCost": 80,
    "latestStringingRecord": {
      "id": 1,
      "racketId": 1,
      "stringName": "Poly Tour Pro",
      "storeName": "冠军穿线工作室",
      "verticalTension": 48,
      "horizontalTension": 46,
      "cost": 80,
      "stringDate": "2026-05-10 10:00",
      "createdAt": "2026-05-10T10:00:00+08:00",
      "updatedAt": "2026-05-10T10:00:00+08:00"
    },
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
| releaseYear | number | 关联球拍库的发布年份，无 `libraryId` 时为 `0` |
| weight | number | 关联球拍库的重量参数，单位克，无 `libraryId` 时为 `0` |
| headSize | number | 关联球拍库的拍面大小数值，无 `libraryId` 时为 `0` |
| stringPattern | string | 关联球拍库的穿线模式，例如 `16x19`，无 `libraryId` 时为空字符串 |
| fileId | string | 关联球拍库的球拍图片文件 ID，无 `libraryId` 时为空字符串 |
| usageCount | number | 该球拍关联打球记录次数 |
| usageMinutes | number | 该球拍累计使用分钟数 |
| usageHours | number | 该球拍累计使用小时数，当前向下取整 |
| afterStringingUsageCount | number | 最近一次穿线之后的使用次数，无穿线记录时为 `0` |
| afterStringingUsageMinutes | number | 最近一次穿线之后的累计使用分钟数，无穿线记录时为 `0` |
| afterStringingUsageHours | number | 最近一次穿线之后的累计使用小时数，向下取整，无穿线记录时为 `0` |
| stringHealth | object/null | 聚酯线健康度，字段与计算公式见 20.2；无穿线记录时为 `null` |
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
  "status": 1,
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80
}
```

规则：

- `name` 必填。
- `status` 可选，支持 `1` 主力拍、`2` 在用；不传默认 `2`。
- `status = 1` 时会自动将当前用户其他主力拍改为在用。
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
    "purchaseDate": "2026-01-01",
    "purchasePrice": 1599,
    "releaseYear": 2025,
    "weight": 300,
    "headSize": 100,
    "stringPattern": "16x19",
    "fileId": "cloud://tennisdaily/rackets/ezone100.png",
    "stringName": "Poly Tour Pro",
    "storeName": "冠军穿线工作室",
    "verticalTension": 48,
    "horizontalTension": 46,
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
      "storeName": "冠军穿线工作室",
      "verticalTension": 48,
      "horizontalTension": 46,
      "cost": 80,
      "stringDate": "2026-05-10 19:30",
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
- 最近一次穿线按 `string_date DESC, id DESC` 取第一条未删除穿线记录，`string_date` 精确到分钟。
- 如果存在最近一次穿线记录，穿线后使用统计按 `tennis_sessions.date >= 最近一次穿线记录的 string_date` 统计，精确到分钟。

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
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599,
  "stringName": "Poly Tour Pro",
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
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
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "cost": 80,
  "stringDate": "2026-05-10 19:30"
}
```

### 20.12 编辑穿线记录

```http
PUT /api/rackets/:id/stringing-records/:recordId
Authorization: Bearer <token>
Content-Type: application/json
```

请求体同新增穿线记录，返回更新后的穿线记录。

### 20.13 删除穿线记录

```http
DELETE /api/rackets/:id/stringing-records/:recordId
Authorization: Bearer <token>
```

响应：

```json
{
  "deleted": true
}
```

### 20.14 打球记录关联球拍

新增/编辑打球记录支持传入：

```json
{
  "racketId": 1
}
```

`racketId` 用于球拍累计使用时长统计；球拍名称不再作为打球记录的存储字段，接口返回的 `racketName` 按 `racketId` 关联我的球拍当前名称。

---

## 21. 球拍库与添加球拍选择接口

### 21.1 获取球拍品牌

添加球拍页面使用。返回系统维护的球拍品牌列表。

```http
GET /api/racket-brands
Authorization: Bearer <token>
```

响应 data 示例：

```json
[
  {
    "id": 1,
    "name": "Yonex",
    "fileId": "cloud://tennisdaily/brands/yonex.png"
  },
  {
    "id": 2,
    "name": "Wilson",
    "fileId": "cloud://tennisdaily/brands/wilson.png"
  }
]
```

### 21.2 获取球拍系列

添加球拍页面使用。返回系统维护的球拍系列列表，支持按品牌过滤。

```http
GET /api/racket-series
GET /api/racket-series?brandId=1
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选；传入后只返回该品牌下的系列 |

响应 data 示例：

```json
[
  {
    "id": 101,
    "brandId": 1,
    "name": "EZONE"
  },
  {
    "id": 102,
    "brandId": 1,
    "name": "VCORE"
  }
]
```

### 21.3 获取球拍库，按品牌分类

添加球拍页面使用。返回系统球拍库中的球拍，并按品牌分组。

```http
GET /api/racket-library
GET /api/racket-library?brandId=1
GET /api/racket-library?brandId=1&seriesId=101
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选；传入后只返回该品牌球拍 |
| seriesId | number | 系列 ID，可选；传入后只返回该系列球拍，可与 `brandId` 组合使用 |

响应 data 示例：

```json
[
  {
    "brandId": 1,
    "brand": "Yonex",
    "items": [
      {
        "id": 1,
        "brandId": 1,
        "brand": "Yonex",
        "seriesId": 101,
        "series": "EZONE",
        "model": "EZONE 100",
        "releaseYear": 2025,
        "weight": 300,
        "headSize": 100,
        "stringPattern": "16x19",
        "fileId": "cloud://tennisdaily/rackets/ezone100.png",
        "imageUrl": ""
      }
    ]
  },
  {
    "brandId": 2,
    "brand": "Wilson",
    "items": [
      {
        "id": 7,
        "brandId": 2,
        "brand": "Wilson",
        "seriesId": 201,
        "series": "Blade",
        "model": "Blade 98 16x19",
        "releaseYear": 2024,
        "weight": 305,
        "headSize": 98,
        "stringPattern": "16x19",
        "fileId": "cloud://tennisdaily/rackets/blade98.png",
        "imageUrl": ""
      }
    ]
  }
]
```

### 21.4 获取球拍库品牌和系列数量统计

返回球拍库中每个品牌的球拍数量，以及每个品牌下每个系列的球拍数量。

```http
GET /api/racket-library/stats
Authorization: Bearer <token>
```

响应 data 示例：

```json
[
  {
    "brandId": 1,
    "brand": "Yonex",
    "count": 12,
    "series": [
      {
        "seriesId": 101,
        "series": "EZONE",
        "count": 5
      },
      {
        "seriesId": 102,
        "series": "VCORE",
        "count": 4
      }
    ]
  },
  {
    "brandId": 2,
    "brand": "Wilson",
    "count": 9,
    "series": [
      {
        "seriesId": 201,
        "series": "Blade",
        "count": 6
      }
    ]
  }
]
```

说明：

- 品牌按品牌球拍数量 `count` 降序返回。
- 每个品牌下的系列按系列球拍数量 `count` 降序返回。
- 统计来源为 `racket_library`。

### 21.5 从球拍库添加到我的球拍

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
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "lastStringDate": "2026-05-10",
  "lastStringCost": 80
}
```

说明：

- `headSize` 表示拍面大小数值，例如 `98`，展示单位可拼接为 `sq.in.`。
- `stringPattern` 表示穿线模式，例如 `16x19`、`18x20`、`16/19`。
- `fileId` 表示球拍图片文件 ID，可用于前端按文件 ID 加载图片资源。
- `brandId`、`seriesId`、`series` 已从 `racket_library` 返回，前端可直接用于品牌/系列展示或筛选。
- `GET /api/racket-library` 支持通过 query 参数 `brandId`、`seriesId` 过滤球拍库列表。
- `GET /api/racket-library/stats` 返回球拍库品牌和系列数量统计。
- `brandId`、`seriesId` 目前由球拍库自身字段承载，未引入独立品牌/系列表。
- `libraryId` 有值时，后端会从球拍库补全 `brand`、`model`。
- 如果请求体里也传了 `brand`、`model`，以前端传入值为准。
- `name` 仍可自定义，比如“EZONE 主力拍”。

### 21.6 添加库中没有的球拍

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
- 如果球拍库中没有对应系列，`series` 可不传。
- 当前接口不再接收/处理 `stringName`、`verticalTension`、`horizontalTension`、`lastStringDate`、`lastStringCost`。

### 21.6 获取我的球拍，用于新增打球记录时选择

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
    "storeName": "冠军穿线工作室",
    "verticalTension": 48,
    "horizontalTension": 46,
    "lastStringDate": "2026-05-10",
    "totalHours": 86
  }
]
```

新增打球记录时，把选中的我的球拍写入：

```json
{
  "racketId": 1
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
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599
}
```

说明：

- `name` 必填。
- `libraryId` 可选。
- 从球拍库选择时，后端可根据 `libraryId` 补全 `brand`、`model`。
- 库中没有的球拍，用户可不传 `libraryId`，直接手动输入 `name`、`brand`、`model`。
- 当前接口不再接收/处理 `stringName`、`verticalTension`、`horizontalTension`、`lastStringDate`、`lastStringCost`。

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
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1599
}
```

说明：

- `status = 1`：设为主力拍，后端会自动将当前用户其他主力拍改为在用。
- `status = 2`：在用。
- `status = 3`：退役。
- 编辑球拍不再接收/处理穿线信息。

### 22.3 获取主力球拍

```http
GET /api/my-rackets/primary
Authorization: Bearer <token>
```

响应 `data` 为主力球拍的 `RacketResponse`，无主力球拍时返回 `null`。返回字段包含球拍信息、`latestStringingRecord` 最新一条穿线信息、`usageCount` 累计使用次数、`usageHours` 累计使用小时数，以及 `afterStringingUsageCount` / `afterStringingUsageMinutes` / `afterStringingUsageHours` 穿线后使用统计和 `stringHealth` 聚酯线健康度。

说明：

- 只返回当前用户未删除且 `status = 1` 的主力球拍。
- 响应包含 `latestStringingRecord` 最近一次穿线记录。
- 响应包含 `usageCount`、`usageMinutes`、`usageHours`，表示该球拍在打球记录中的累计使用次数和累计时间。
- 响应包含 `afterStringingUsageCount`、`afterStringingUsageMinutes`、`afterStringingUsageHours`，表示最近一次穿线之后的累计使用统计。
- 响应包含 `stringHealth` 聚酯线健康度，字段结构和计算公式见 20.2；没有穿线记录时为 `null`。
- 如果历史数据异常存在多把主力拍，返回最近更新的一把。

### 22.4 删除球拍

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

### 22.5 新增穿线记录

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
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "cost": 80,
  "stringDate": "2026-05-10 19:30"
}
```

说明：

- `stringName` 必填。
- `storeName` 为穿线门店，可选，未填写时为空字符串。
- `verticalTension` 表示竖线磅数，`horizontalTension` 表示横线磅数。
- `stringDate` 必填，格式 `YYYY-MM-DD HH:mm`，旧格式 `YYYY-MM-DD` 仍兼容并按当天 `00:00` 处理。
- 球拍列表和详情中的当前球线信息来自最近一条穿线记录。

### 22.6 编辑穿线记录

```http
PUT /api/rackets/:id/stringing-records/:recordId
Authorization: Bearer <token>
Content-Type: application/json
```

请求体：

```json
{
  "stringName": "Poly Tour Pro",
  "storeName": "冠军穿线工作室",
  "verticalTension": 48,
  "horizontalTension": 46,
  "cost": 80,
  "stringDate": "2026-05-10 19:30"
}
```

说明：

- `recordId` 为穿线记录 ID，必须属于当前用户和路径中的球拍。
- `stringName` 和 `stringDate` 必填。
- `storeName` 为穿线门店，可选，未填写时为空字符串。
- `verticalTension` 表示竖线磅数，`horizontalTension` 表示横线磅数。
- `stringDate` 格式为 `YYYY-MM-DD HH:mm`，旧格式 `YYYY-MM-DD` 仍兼容并按当天 `00:00` 处理。
- 响应返回更新后的穿线记录。

### 22.7 删除穿线记录

```http
DELETE /api/rackets/:id/stringing-records/:recordId
Authorization: Bearer <token>
```

响应：

```json
{
  "deleted": true
}
```

说明：

- `recordId` 为穿线记录 ID，必须属于当前用户和路径中的球拍。
- 当前实现为逻辑删除，更新 `deleted_at`。
- 删除后球拍列表和详情中的最近穿线信息会自动排除该记录。

### 22.8 保留我的球拍接口

新增打球记录选择球拍时继续使用：

```http
GET /api/my-rackets
Authorization: Bearer <token>
```

只返回未删除且未退役的球拍：

```text
status IN (1, 2)
```

---

## 23. 球鞋管理接口

球鞋功能整体流程与球拍一致：球鞋库（品牌/系列/性别）→ 添加我的球鞋（库选择补全 + 手动输入）→ 我的球鞋管理（列表/详情/编辑/主力/退役/删除）。

打球记录已支持 `shoeId` 关联我的球鞋，响应 `shoeName` 按 `shoeId` 实时返回；未关联（`shoeId = 0`）的记录回退到快照 `shoeName`。

### 23.1 球鞋状态

| 值 | 含义 |
|---:|---|
| `1` | 主力鞋 |
| `2` | 在用 |
| `3` | 已退役 |

### 23.2 获取球鞋品牌

```http
GET /api/shoe-brands
Authorization: Bearer <token>
```

响应 data：

```json
[
  {
    "id": 1,
    "name": "Asics",
    "fileId": "cloud://tennisdaily/shoe-brands/asics.png"
  }
]
```

说明：`shoe_brands.slug` 为内部品牌标识，不返回给前端。

### 23.3 获取球鞋系列

```http
GET /api/shoe-series
GET /api/shoe-series?brandId=1
GET /api/shoe-series?brandId=1&gender=1
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选 |
| gender | number | 性别，可选；`0` 未知 / `1` 男 / `2` 女 / `3` 童，不传返回全部 |

响应 data：

```json
[
  {
    "id": 101,
    "brandId": 1,
    "gender": 1,
    "name": "Gel Resolution"
  }
]
```

### 23.4 获取球鞋库，按品牌分类

```http
GET /api/shoe-library
GET /api/shoe-library?brandId=1&seriesId=101&gender=1
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选 |
| seriesId | number | 系列 ID，可选 |
| gender | number | 性别，可选；`0` 未知 / `1` 男 / `2` 女 / `3` 童 |

规则：

- 同一型号不同配色各占一行，每条记录有独立 `id`。
- 每个品牌分组下按 `release_year DESC, id ASC` 排序。

响应 data：

```json
[
  {
    "brandId": 1,
    "brand": "Asics",
    "items": [
      {
        "id": 1,
        "brandId": 1,
        "brand": "Asics",
        "seriesId": 101,
        "series": "Gel Resolution",
        "model": "Gel Resolution 9",
        "gender": 1,
        "releaseYear": 2025,
        "colorway": "White/Orewood Brown",
        "weight": "14.6 ounces (size 10.5)",
        "width": "Snug Medium",
        "surface": "Hard (all court)",
        "price": 1090,
        "colorwayCount": 5,
        "fileId": "cloud://tennisdaily/shoes/gel9.png",
        "imageUrl": ""
      }
    ]
  }
]
```

说明：`product_code` 为内部字段，不返回给前端。

### 23.5 获取球鞋库品牌和系列数量统计

```http
GET /api/shoe-library/stats
Authorization: Bearer <token>
```

按品牌分组统计行数（含不同配色行，品牌总数不区分性别）；品牌下 `series` 为按性别分组的哈希，key 为性别字符串（`"1"` 男 / `"2"` 女 / `"0"` 未知 / `"3"` 童），值为该性别下的系列统计数组。

响应 data 示例：

```json
[
  {
    "brandId": 1,
    "brand": "Asics",
    "count": 13,
    "series": {
      "1": [ { "seriesId": 101, "series": "Gel Resolution", "count": 8 } ],
      "2": [ { "seriesId": 101, "series": "Gel Resolution", "count": 5 } ]
    }
  }
]
```

### 23.6 获取我的球鞋列表

```http
GET /api/shoes
GET /api/shoes?includeRetired=true
Authorization: Bearer <token>
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| includeRetired | boolean | 是否包含已退役球鞋，默认 `false` |

排序：`status ASC, created_at DESC`（主力鞋在前）。

响应 data：`ShoeResponse` 数组，字段见 23.7。

### 23.7 添加球鞋

```http
POST /api/shoes
Authorization: Bearer <token>
Content-Type: application/json
```

从球鞋库选择：

```json
{
  "libraryId": 1,
  "name": "Gel Resolution 主力鞋",
  "size": "42",
  "colorway": "White/Orewood Brown",
  "status": 1,
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1090
}
```

手动输入库中没有的球鞋：

```json
{
  "name": "我的旧款备用鞋",
  "brand": "Babolat",
  "model": "Jet Mach III",
  "size": "42.5",
  "purchasePrice": 899
}
```

请求体字段：

| 字段 | 类型 | 必填 | 说明 |
|---|---|:---:|---|
| libraryId | number | 否 | 球鞋库配色行 ID，不同配色各占一行；传 `0` 或不传表示手动输入 |
| name | string | 是 | 球鞋名称 |
| brand | string | 否 | 品牌；从库选择时后端按 `libraryId` 补全 |
| model | string | 否 | 型号；从库选择时后端按 `libraryId` 补全 |
| status | number | 否 | `1` 主力鞋 / `2` 在用，默认 `2`；创建时不接受 `3` |
| size | string | 否 | 尺码，如 `42`、`42.5` |
| colorway | string | 否 | 配色；从库选择时可由库补全 |
| purchaseDate | string | 否 | 购买日期，格式 `YYYY-MM-DD` |
| purchasePrice | number | 否 | 购买价格 |

规则：

- `libraryId` 对应记录必须存在，否则返回 `40001 invalid request`。
- 显式传入的 `brand`、`model` 优先于库补全；`name` 为空时用 `brand + " " + model` 补全。
- `gender`、`releaseYear`、`fileId` 通过 `libraryId` 关联球鞋库实时补全，不冗余到 `my_shoes` 表。
- `status = 1` 时，自动将当前用户其他主力鞋改为在用。

响应 data（`ShoeResponse`）：

```json
{
  "id": 1,
  "libraryId": 1,
  "name": "Gel Resolution 主力鞋",
  "brand": "Asics",
  "model": "Gel Resolution 9",
  "status": 1,
  "gender": 1,
  "size": "42",
  "colorway": "White/Orewood Brown",
  "purchaseDate": "2026-07-15",
  "purchasePrice": 1090,
  "releaseYear": 2025,
  "fileId": "cloud://tennisdaily/shoes/gel9.png",
  "usageCount": 4,
  "usageMinutes": 480,
  "usageHours": 8,
  "totalMinutes": 480,
  "totalHours": 8,
  "wear": {
    "state": "good",
    "display": "状态良好 · 预计还可打 49h",
    "score": 82.5,
    "remainingHours": 49
  },
  "createdAt": "2026-08-05T10:00:00+08:00",
  "updatedAt": "2026-08-05T10:00:00+08:00"
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 我的球鞋 ID |
| libraryId | number | 球鞋库配色行 ID，手动输入时为 `0` |
| name | string | 球鞋名称 |
| brand | string | 品牌 |
| model | string | 型号 |
| status | number | `1` 主力鞋 / `2` 在用 / `3` 已退役 |
| gender | number | 性别，关联球鞋库补全，手动输入时为 `0` |
| size | string | 尺码 |
| colorway | string | 配色 |
| purchaseDate | string | 购买日期，`YYYY-MM-DD`，可为空字符串 |
| purchasePrice | number | 购买价格，可为 `null` |
| releaseYear | number | 关联球鞋库的上市年份，无 `libraryId` 时为 `0` |
| fileId | string | 关联球鞋库的图片文件 ID，无 `libraryId` 时为空字符串 |
| usageCount | number | 关联打球记录（按 `shoeId`）的上场次数 |
| usageMinutes | number | 关联打球记录的上场总分钟数 |
| usageHours | number | 上场总分钟数换算的小时数（分钟数 `/60` 取整） |
| totalMinutes / totalHours | number | 累计使用时长，当前口径与 `usageMinutes` / `usageHours` 一致 |
| wear | object | 磨损度估算，无购买日期时按场上用时估算；`null` 时表示未启用磨损度配置 |
| createdAt / updatedAt | string | 创建/更新时间，RFC3339 格式 |

`wear` 字段说明（口径参考球线健康度估算）：

| 字段 | 类型 | 说明 |
|---|---|---|
| state | string | 状态 key：`fresh` 全新 / `good` 状态良好 / `worn` 开始磨损 / `tired` 磨损明显 / `dead` 缓震衰减 / `expired` 已超期 |
| display | string | 展示文案，`{remainingHours}` 已替换为剩余小时数 |
| score | number | 磨损度得分，`0-100`，越高越新 |
| remainingHours | number | 预计剩余可使用小时数 |

磨损度估算规则（`config.yaml` 的 `shoeWear` 段可配置）：

- `effectiveWear = 上场小时数 + 购买后天数 × restWearPerDay`
- `score = max(0, 100 × (1 − effectiveWear / standardLifeHours))`
- 默认 `standardLifeHours = 60`（约 45-60 小时的中底寿命经验值）、`restWearPerDay = 0.06`（约 1000 天自然老化归零）

### 23.8 获取球鞋详情

```http
GET /api/shoes/:id
Authorization: Bearer <token>
```

响应 data 为单个 `ShoeResponse`；不存在或不属于当前用户时返回 `40401 not found`。

### 23.9 编辑球鞋

```http
PUT /api/shoes/:id
Authorization: Bearer <token>
Content-Type: application/json
```

请求体同 23.7，`status` 支持 `1`/`2`/`3`。`status = 1` 时自动将其他主力鞋改为在用。返回更新后的 `ShoeResponse`。

### 23.10 设置主力鞋

```http
POST /api/shoes/:id/set-primary
Authorization: Bearer <token>
```

规则：

- 同一用户同一时间只允许一支主力鞋。
- 设置成功后，原主力鞋自动变为在用。

### 23.11 退役球鞋

```http
POST /api/shoes/:id/retire
Authorization: Bearer <token>
```

规则：

- 退役后 `status = 3`。
- 默认列表不展示，`/api/shoes/selectable`、`/api/my-shoes` 不返回退役球鞋。

### 23.12 删除球鞋

```http
DELETE /api/shoes/:id
Authorization: Bearer <token>
```

规则：

- 逻辑删除（更新 `deleted_at`），不做物理删除。
- 删除后从列表、详情、主力鞋查询中消失。

### 23.13 获取球鞋统计

```http
GET /api/shoes/stats
Authorization: Bearer <token>
```

统计口径：

- `shoeCount`：当前用户未删除球鞋总数（含已退役，不含逻辑删除）。
- `shoeCost`：当前用户未删除球鞋的 `purchase_price` 合计，空值按 `0` 处理。
- `totalCost`：当前等于 `shoeCost`。

响应 data：

```json
{
  "shoeCount": 3,
  "shoeCost": 2879,
  "totalCost": 2879,
  "shoeCostText": "2879.00",
  "totalCostText": "2879.00"
}
```

### 23.14 可选球鞋 / 我的球鞋

```http
GET /api/shoes/selectable
GET /api/my-shoes
GET /api/my-shoes/primary
Authorization: Bearer <token>
```

说明：

- `/api/shoes/selectable`、`/api/my-shoes` 只返回状态 `1`、`2` 的未删除球鞋，供打球记录页选择；打球记录新增/编辑传 `shoeId` 即可关联。
- `/api/my-shoes/primary` 返回当前主力鞋，没有时返回 `null`；返回结构与 `/api/shoes/:id` 一致，同样带 `usageCount`、`usageMinutes`、`usageHours`、`totalMinutes`、`totalHours` 和 `wear` 磨损度估算。

### 23.15 管理员维护球鞋库接口

管理员白名单在 `config.yaml` 的 `admin.userIds` 配置（用户 ID 来自 JWT）。非管理员访问返回 `40301 forbidden`。

#### 23.15.1 查询管理员权限

```http
GET /api/admin/permissions
Authorization: Bearer <token>
```

响应 data：

```json
{
  "isAdmin": true
}
```

前端用它控制管理员入口显隐。

#### 23.15.2 新增球鞋品牌

```http
POST /api/admin/shoe-brands
Authorization: Bearer <token>
```

请求体：

```json
{
  "name": "Asics",
  "slug": "",
  "fileId": ""
}
```

规则：

- `name` 必填且唯一，重复返回 `invalid request: 品牌已存在`。
- `slug` 可选，内部品牌标识；为空时由品牌名自动生成（转小写、保留字母数字与连字符，纯中文退化为 `brand-<hash>`）。
- `fileId` 可选，品牌图片云存储文件 ID。

响应 data：

```json
{
  "id": 5,
  "name": "Asics",
  "fileId": ""
}
```

#### 23.15.3 新增球鞋系列

```http
POST /api/admin/shoe-series
Authorization: Bearer <token>
```

请求体：

```json
{
  "brandId": 5,
  "gender": 1,
  "name": "Gel Resolution"
}
```

规则：

- `brandId`、`name` 必填；`gender` 取值 `0` 未知 / `1` 男 / `2` 女 / `3` 童，默认 `0`。
- 同一品牌、性别下系列名唯一，重复返回 `invalid request: 该品牌下同名系列已存在`。

响应 data：

```json
{
  "id": 202,
  "brandId": 5,
  "gender": 1,
  "name": "Gel Resolution"
}
```

#### 23.15.4 新增球鞋库鞋款（配色合并为单条）

```http
POST /api/admin/shoe-library
Authorization: Bearer <token>
```

请求体：

```json
{
  "brandId": 5,
  "seriesId": 202,
  "model": "Gel Resolution 9",
  "gender": 1,
  "colorway": "Red/Black",
  "releaseYear": 2025,
  "weight": "14.6 ounces (size 10.5)",
  "width": "Snug Medium",
  "surface": "Hard (all court)",
  "price": 1090,
  "colorwayCount": 5,
  "fileId": "",
  "imageUrl": ""
}
```

规则：

- `brandId`、`seriesId`、`model` 必填；`seriesId` 必须属于 `brandId`。
- `colorway` 为单条配色字符串，前端多选主流颜色后合并为 `Red/Black` 形式提交；每次请求只插入一行 `shoe_library` 记录。
- 按 品牌+系列+型号+性别+配色 查重，配色已存在则返回 `invalid request: 配色 <名称> 已存在`。
- `gender`、`releaseYear`、`weight`、`width`、`surface`、`price`、`colorwayCount`、`fileId`、`imageUrl` 均可选。

响应 data：

```json
{
  "created": 1,
  "items": [
    {
      "id": 301,
      "brandId": 5,
      "brand": "Asics",
      "seriesId": 202,
      "series": "Gel Resolution",
      "model": "Gel Resolution 9",
      "gender": 1,
      "releaseYear": 2025,
      "colorway": "Red/Black",
      "weight": "14.6 ounces (size 10.5)",
      "width": "Snug Medium",
      "surface": "Hard (all court)",
      "price": 1090,
      "colorwayCount": 5,
      "fileId": "",
      "imageUrl": ""
    }
  ]
}
```
