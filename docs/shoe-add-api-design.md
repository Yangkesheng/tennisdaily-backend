# 添加球鞋 API 设计文档（独立设计稿）

> 状态：设计稿 v1.1，已按本设计实现（2026-08-05），见 `docs/api.md` 23. 球鞋管理接口
> 日期：2026-08-05
> 关联需求：增加添加网球鞋/球鞋功能，整体流程与添加球拍保持一致

> 2026-08-05 变更：`shoe_library` 取消发布状态（`draft / published / archived`），
> 库表删除 `status` 字段，C 端默认读取全部库数据，不再有发布流程（见 `migrations/017_drop_shoe_library_status.sql`）。

> 2026-08-05 更新：接入确认后的球鞋库表结构（`shoe_brands` / `shoe_series` / `shoe_library`），完整 DDL 见 `migrations/015_create_shoe_tables.sql`；`shoe`（我的球鞋）表未包含在确认的 DDL 中，本节保留 v1 草案，最终需确认（见 12.6）。

---

## 1. 背景与目标

当前小程序打球记录只有 `shoeName` 自由文本字段，用户无法像球拍一样维护自己的球鞋、从球鞋库选择、设置主力鞋并统计购买花费。

本次设计目标：

- 提供与球拍一致的"球鞋库 → 添加我的球鞋"完整流程。
- 支持从系统球鞋库选择（按品牌/系列筛选），也支持手动输入库中没有的球鞋。
- 我的球鞋支持主力鞋/在用/退役状态管理与逻辑删除。
- 本次只输出 API 设计与数据库设计，不写代码。

本次不包含（如需实现需另行确认，见第 12 节）：

- 球鞋磨损/更换记录。

> 2026-08-05 已确认并实现：打球记录支持 `shoeId` 关联；球鞋购买费用计入首页/统计 expense（新增 `shoeCost`）。见 `docs/api.md` 与 `docs/change-log.md`。

## 2. 设计原则

- 与球拍流程一一对应，接口路径、请求/响应结构、状态枚举、软删除规则全部对齐球拍。
- 我的球鞋表名使用 `my_shoes`；品牌/系列/库表使用复数。
- 所有接口沿用统一响应格式 `{code, message, data}`。
- 所有用户数据归属只来自 JWT 中的当前用户 ID，不接收前端 `userId`。
- 所有查询/统计过滤 `deleted_at IS NULL`。
- 日期解析统一 `Asia/Shanghai`，`purchaseDate` 使用 `YYYY-MM-DD`。
- JSON 字段使用 camelCase，数据库字段使用 snake_case。

## 3. 与球拍流程的对应关系

| 球拍 | 球鞋 | 说明 |
|---|---|---|
| `racket_brands` | `shoe_brands` | 品牌表 |
| `racket_series` | `shoe_series` | 系列表 |
| `racket_library` | `shoe_library` | 球鞋库 |
| `racket` | `my_shoes` | 我的球鞋表 |
| `GET /api/racket-brands` | `GET /api/shoe-brands` | 品牌列表 |
| `GET /api/racket-series` | `GET /api/shoe-series` | 系列列表 |
| `GET /api/racket-library` | `GET /api/shoe-library` | 球鞋库（按品牌分组） |
| `GET /api/racket-library/stats` | `GET /api/shoe-library/stats` | 库品牌/系列统计 |
| `GET /api/rackets` | `GET /api/shoes` | 我的球鞋列表 |
| `POST /api/rackets` | `POST /api/shoes` | 添加球鞋（核心接口） |
| `GET /api/rackets/:id` | `GET /api/shoes/:id` | 球鞋详情 |
| `PUT /api/rackets/:id` | `PUT /api/shoes/:id` | 编辑球鞋 |
| `DELETE /api/rackets/:id` | `DELETE /api/shoes/:id` | 删除球鞋（逻辑删除） |
| `POST /api/rackets/:id/set-primary` | `POST /api/shoes/:id/set-primary` | 设置主力鞋 |
| `POST /api/rackets/:id/retire` | `POST /api/shoes/:id/retire` | 退役球鞋 |
| `GET /api/rackets/selectable` | `GET /api/shoes/selectable` | 可选球鞋（状态 1/2） |
| `GET /api/my-rackets` | `GET /api/my-shoes` | 打球记录页选择用 |
| `GET /api/my-rackets/primary` | `GET /api/my-shoes/primary` | 当前主力鞋 |
| `GET /api/rackets/stats` | `GET /api/shoes/stats` | 球鞋数量/花费统计 |

球拍添加流程中已取消穿线信息（见 `docs/api.md` 22.1），球鞋本身没有与穿线对应的子记录，因此添加球鞋流程更简单，只维护球鞋本体信息。

## 4. 名词说明

| 名词 | 含义 |
|---|---|
| 球鞋库（`shoe_library`） | 系统维护的球鞋型号库，按品牌/系列/型号组织 |
| 我的球鞋（`my_shoes`） | 当前用户添加到自己的球鞋 |
| 主力鞋 | 状态 `1`，同一用户同一时间最多一支 |
| 在用 | 状态 `2`，默认状态 |
| 已退役 | 状态 `3`，默认列表不展示 |

## 5. 新增球鞋页面流程（与添加球拍一致）

1. 进入"我的球鞋"页，点击添加球鞋。
2. 可选两条路径：
   - 从球鞋库选择：浏览/筛选品牌、系列、性别 → 选择型号/配色 → 自动带入名称、品牌、型号、配色、图片。
   - 手动输入：库中没有时直接填写名称、品牌、型号。
3. 补充填写尺码、配色、购买日期、购买价格，选择状态（主力鞋/在用）。
4. 提交 `POST /api/shoes`，创建成功返回完整球鞋信息，进入我的球鞋列表。
5. 后续可编辑、设置主力、退役、删除。

## 6. 数据模型设计

### 6.1 表结构（已确认，完整 DDL 见 `migrations/015_create_shoe_tables.sql`）

本轮确认的球鞋库表共三张：`shoe_brands`（品牌）、`shoe_series`（系列）、`shoe_library`（球鞋库，业务主表）。要点：

- 品牌表：`name`、`slug` 均唯一；`slug` 为内部品牌标识，不暴露给前端。
- 系列表：品牌下按 `gender` 区分系列，唯一键为 `(brand_id, gender, name)`。
- 球鞋库表：一个型号可有多行配色（`colorway` 不同配色各占一行），无唯一去重键；`product_code` 仅用于回源关联 `source_shoes`，不是去重键。
- 球鞋库 `status` 取值 `draft / published / archived`，所有 C 端查询只返回 `published`。
- `brand_id`、`series_id` 仅逻辑关联，不使用数据库外键（与球拍一致）。

```sql
-- 品牌表
CREATE TABLE IF NOT EXISTS `shoe_brands` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '品牌ID，主键',
  `name` VARCHAR(128) NOT NULL COMMENT '品牌名称，例如 Nike、adidas、Asics',
  `slug` VARCHAR(128) NOT NULL COMMENT '品牌标识，例如 nike、adidas、asics',
  `file_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '品牌图片文件ID',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shoe_brands_name` (`name`),
  UNIQUE KEY `uk_shoe_brands_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='球鞋品牌表';

-- 系列表
CREATE TABLE IF NOT EXISTS `shoe_series` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '系列ID，主键',
  `brand_id` BIGINT UNSIGNED NOT NULL COMMENT '品牌ID，仅逻辑关联 shoe_brands.id，不使用数据库外键',
  `gender` SMALLINT NOT NULL DEFAULT 0 COMMENT '性别：0未知 1男 2女 3童',
  `name` VARCHAR(128) NOT NULL COMMENT '系列名称，例如 Vapor、Court FF、Gel Resolution',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_shoe_series_brand_gender_name` (`brand_id`, `gender`, `name`),
  KEY `idx_shoe_series_brand_id` (`brand_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='球鞋系列表';

-- 球鞋库（业务主表）
CREATE TABLE IF NOT EXISTS `shoe_library` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '库表ID，主键',
  `brand_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '品牌ID，仅逻辑关联 shoe_brands.id',
  `brand` VARCHAR(50) NOT NULL COMMENT '品牌名称冗余字段，例如 Nike',
  `series_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '系列ID，仅逻辑关联 shoe_series.id',
  `series` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '系列名称冗余字段，例如 Vapor',
  `model` VARCHAR(100) NOT NULL COMMENT '型号，例如 Vapor 12',
  `gender` SMALLINT NOT NULL DEFAULT 0 COMMENT '性别：0未知 1男 2女 3童',
  `colorway` VARCHAR(128) NULL DEFAULT NULL COMMENT '配色，例如 White/Orewood Brown，不同配色各占一行',
  `product_code` VARCHAR(64) NULL DEFAULT NULL COMMENT '来源商品代码，仅作为回源关联 source_shoes 的列，不是去重键',
  `release_year` INT NOT NULL DEFAULT 0 COMMENT '上市年份',
  `weight` VARCHAR(128) NULL DEFAULT NULL COMMENT '单只重量，例如 14.6 ounces (size 10.5)',
  `width` VARCHAR(128) NULL DEFAULT NULL COMMENT '鞋宽（版型）：Snug Medium、Wide、Narrow',
  `surface` VARCHAR(255) NULL DEFAULT NULL COMMENT '适用场地：Hard (all court)、Clay 等',
  `price` DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '当前售价',
  `colorway_count` INT NOT NULL DEFAULT 0 COMMENT '该型号在售配色数量',
  `file_id` VARCHAR(255) NULL DEFAULT NULL COMMENT '首图云存储fileID',
  `image_url` VARCHAR(500) NULL DEFAULT NULL COMMENT '首图访问URL',
  `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '状态：draft、published、archived',
  `created_at` DATETIME(3) NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` DATETIME(3) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_shoe_library_brand_id` (`brand_id`),
  KEY `idx_shoe_library_series_id` (`series_id`),
  KEY `idx_shoe_library_brand` (`brand`),
  KEY `idx_shoe_library_model` (`model`),
  KEY `idx_shoe_library_gender` (`gender`),
  KEY `idx_shoe_library_colorway` (`colorway`),
  KEY `idx_shoe_library_product_code` (`product_code`),
  KEY `idx_shoe_library_surface` (`surface`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='球鞋库表（业务主表）';
```

#### 我的球鞋表 `my_shoes`（已确认）

已确认并追加到迁移脚本，实际 DDL 见 `migrations/015_create_shoe_tables.sql`：

```sql
CREATE TABLE IF NOT EXISTS my_shoes (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '球鞋ID',
  user_id BIGINT NOT NULL COMMENT '用户ID',
  library_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '球鞋库ID',

  name VARCHAR(100) NOT NULL COMMENT '球鞋名称',
  brand VARCHAR(50) DEFAULT NULL COMMENT '品牌',
  model VARCHAR(100) DEFAULT NULL COMMENT '型号',
  status TINYINT NOT NULL DEFAULT 2 COMMENT '状态:1主力鞋 2在用 3已退役',
  size VARCHAR(20) DEFAULT NULL COMMENT '尺码，例如 42、42.5',
  colorway VARCHAR(128) DEFAULT NULL COMMENT '配色，从库选择时默认带入，可修改',
  purchase_date DATE DEFAULT NULL COMMENT '购买日期',
  purchase_price DECIMAL(10,2) DEFAULT NULL COMMENT '购买价格',

  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间',

  CONSTRAINT my_shoes_status_check CHECK (status IN (1, 2, 3)),
  INDEX idx_my_shoes_user_deleted_status (user_id, deleted_at, status),
  INDEX idx_my_shoes_user_deleted_created (user_id, deleted_at, created_at),
  INDEX idx_my_shoes_user_library (user_id, library_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='我的球鞋表';
```

草案说明：

- `size` 是用户维度字段（每个人尺码不同），放在 `my_shoes` 表；球鞋库不保存尺码。
- `colorway` 同时在球鞋库和我的球鞋中出现：从库选择时默认带入，用户可修改。
- `release_year`、`file_id` 通过 `libraryId` 关联球鞋库实时补全，不冗余到 `my_shoes` 表（与球拍一致）。

### 6.2 状态枚举

| 值 | 含义 |
|---:|---|
| `1` | 主力鞋 |
| `2` | 在用 |
| `3` | 已退役 |

### 6.3 软删除

- `shoe.deleted_at` 逻辑删除。
- 列表、详情、统计、主力鞋查询全部过滤 `deleted_at IS NULL`。
- 删除后的历史数据不参与任何展示与统计（当前阶段球鞋尚未被打球记录引用，无历史影响）。

## 7. 接口设计

所有接口均在 `/api` 下，除登录/健康检查外均需要 JWT：

```http
Authorization: Bearer <token>
```

### 7.1 获取球鞋品牌

```http
GET /api/shoe-brands
```

响应 data：

```json
[
  {
    "id": 1,
    "name": "Asics",
    "fileId": "cloud://tennisdaily/shoe-brands/asics.png"
  },
  {
    "id": 2,
    "name": "Nike",
    "fileId": "cloud://tennisdaily/shoe-brands/nike.png"
  }
]
```

字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 品牌 ID |
| name | string | 品牌名称 |
| fileId | string | 品牌图片文件 ID，可为空字符串 |

说明：

- `shoe_brands.slug` 为内部品牌标识，不返回给前端。

### 7.2 获取球鞋系列

```http
GET /api/shoe-series
GET /api/shoe-series?brandId=1
GET /api/shoe-series?brandId=1&gender=1
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选；传入后只返回该品牌下的系列 |
| gender | number | 性别，可选；`0` 未知 / `1` 男 / `2` 女 / `3` 童，不传返回全部 |

响应 data：

```json
[
  {
    "id": 101,
    "brandId": 1,
    "gender": 1,
    "name": "Gel Resolution"
  },
  {
    "id": 102,
    "brandId": 1,
    "gender": 1,
    "name": "Court FF"
  }
]
```

字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 系列 ID |
| brandId | number | 品牌 ID |
| gender | number | 性别：`0` 未知 / `1` 男 / `2` 女 / `3` 童 |
| name | string | 系列名称 |

### 7.3 获取球鞋库，按品牌分类

```http
GET /api/shoe-library
GET /api/shoe-library?brandId=1
GET /api/shoe-library?brandId=1&seriesId=101
GET /api/shoe-library?brandId=1&gender=1
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| brandId | number | 品牌 ID，可选 |
| seriesId | number | 系列 ID，可选，可与 `brandId` 组合 |
| gender | number | 性别，可选；`0` 未知 / `1` 男 / `2` 女 / `3` 童 |

规则：

- 只返回 `status = 'published'` 的球鞋库数据（`draft`、`archived` 不返回）。
- `colorway` 不同配色各占一行，一个型号可能对应多条记录，每条记录有独立 `id`。

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

字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 球鞋库 ID |
| brandId | number | 品牌 ID |
| brand | string | 品牌 |
| seriesId | number | 系列 ID |
| series | string | 系列 |
| model | string | 型号 |
| gender | number | 性别：`0` 未知 / `1` 男 / `2` 女 / `3` 童 |
| releaseYear | number | 版本年份 |
| colorway | string | 配色 |
| weight | string | 单只重量（含尺码说明），可为空字符串 |
| width | string | 鞋宽（版型），可为空字符串 |
| surface | string | 适用场地，可为空字符串 |
| price | number | 当前售价 |
| colorwayCount | number | 该型号在售配色数量 |
| fileId | string | 球鞋图片文件 ID |
| imageUrl | string | 球鞋图片 URL，可为空字符串 |

说明：

- `product_code`、`status` 为内部字段，不返回给前端。
- 每个品牌分组下按 `release_year DESC, id ASC` 排序（与球拍库的型号排序一致）。

### 7.4 获取球鞋库品牌和系列数量统计

```http
GET /api/shoe-library/stats
```

返回结构与 `GET /api/racket-library/stats` 一致，统计来源为 `shoe_library`：

```json
[
  {
    "brandId": 1,
    "brand": "Asics",
    "count": 12,
    "series": [
      {
        "seriesId": 101,
        "series": "Gel Resolution",
        "count": 5
      }
    ]
  }
]
```

统计口径：

- 只统计 `status = 'published'` 的数据。
- 按 `brand` 分组统计行数（含不同配色行）；系列按 `series` 分组统计行数。

### 7.5 添加球鞋（核心接口）

```http
POST /api/shoes
Content-Type: application/json
```

#### 从球鞋库选择

请求体：

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

#### 手动输入库中没有的球鞋

请求体：

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
| libraryId | number | 否 | 球鞋库 ID（`shoe_library` 中 `published` 的配色行 ID，不同配色各占一行）；传 `0` 或不传表示手动输入 |
| name | string | 是 | 球鞋名称，可为"品牌 + 型号"自定义别名 |
| brand | string | 否 | 品牌；从库选择时后端按 `libraryId` 补全 |
| model | string | 否 | 型号；从库选择时后端按 `libraryId` 补全 |
| status | number | 否 | `1` 主力鞋 / `2` 在用，默认 `2`；创建时不接受 `3` |
| size | string | 否 | 尺码，如 `42`、`42.5` |
| colorway | string | 否 | 配色；从库选择时可由库补全 |
| purchaseDate | string | 否 | 购买日期，格式 `YYYY-MM-DD` |
| purchasePrice | number | 否 | 购买价格，建议非负 |

补全规则（与球拍一致）：

- `libraryId` 有效时，`brand`、`model` 为空则由球鞋库补全。
- `libraryId` 对应的记录必须是 `status = 'published'`，否则返回 `40001 invalid request`。
- 请求体显式传入的 `brand`、`model` 优先于库补全。
- `name` 为空时，用 `brand + " " + model` 补全。
- `colorway` 为空时，可由球鞋库补全（新增能力，球拍没有对应字段）。
- `gender`、`releaseYear`、`fileId` 通过 `libraryId` 关联球鞋库实时补全，不冗余到 `my_shoes` 表。
- 从库选择时，库中的 `price`（当前售价）可作为 `purchasePrice` 的前端默认值，由前端带入，后端不强制覆盖。
- `libraryId` 不存在时返回 `40001 invalid request`。

状态规则：

- `status = 1` 时，自动将当前用户其他主力鞋改为在用（与球拍一致，保证只有一支主力鞋）。

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
  "purchaseDate": "2026-01-01",
  "purchasePrice": 1090,
  "releaseYear": 2025,
  "fileId": "cloud://tennisdaily/shoes/gel9.png",
  "createdAt": "2026-08-05T10:00:00+08:00",
  "updatedAt": "2026-08-05T10:00:00+08:00"
}
```

`ShoeResponse` 字段说明：

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
| createdAt / updatedAt | string | 创建/更新时间，RFC3339 格式 |

### 7.6 获取我的球鞋列表

```http
GET /api/shoes
GET /api/shoes?includeRetired=true
```

查询参数：

| 参数 | 类型 | 说明 |
|---|---|---|
| includeRetired | boolean | 是否包含已退役球鞋，默认 `false` |

排序：`status ASC, created_at DESC`（与球拍一致，主力鞋在前）。

响应 data：`ShoeResponse` 数组，字段见 7.5。

### 7.7 获取球鞋详情

```http
GET /api/shoes/:id
```

响应 data 为单个 `ShoeResponse`；不存在或不属于当前用户时返回 `40401 not found`。

### 7.8 编辑球鞋

```http
PUT /api/shoes/:id
Content-Type: application/json
```

请求体同 7.5，扩展 `status` 支持 `3` 已退役（与球拍编辑一致）。

规则：

- 重新传 `libraryId` 时按库补全规则更新 `brand`、`model`。
- `status = 1` 时自动将其他主力鞋改为在用。
- 编辑后返回更新后的 `ShoeResponse`。

### 7.9 设置主力鞋

```http
POST /api/shoes/:id/set-primary
```

规则：

- 同一用户同一时间只允许一支主力鞋。
- 设置成功后，原主力鞋自动变为在用。

### 7.10 退役球鞋

```http
POST /api/shoes/:id/retire
```

规则：

- 退役后 `status = 3`。
- 默认球鞋列表不展示，`/api/shoes/selectable`、`/api/my-shoes` 不返回。

### 7.11 删除球鞋

```http
DELETE /api/shoes/:id
```

规则：

- 逻辑删除（写 `deleted_at`），不做物理删除。
- 删除后从列表、详情、主力鞋查询中消失。
- 若删除的是主力鞋，删除后不再存在主力鞋（与球拍一致，不自动指定新主力）。

### 7.12 可选球鞋 / 我的球鞋（为后续打球记录关联预留）

```http
GET /api/shoes/selectable
GET /api/my-shoes
GET /api/my-shoes/primary
```

说明：

- `/api/shoes/selectable`、`/api/my-shoes` 只返回状态 `1`、`2` 的未删除球鞋，用于打球记录页选择。
- `/api/my-shoes/primary` 返回当前主力鞋，没有时返回 `null`。
- 当前打球记录仍使用 `shoeName` 自由文本，这两个接口先随球鞋功能上线，供后续阶段接入。

### 7.13 获取球鞋统计

```http
GET /api/shoes/stats
```

设计（镜像 `GET /api/rackets/stats`）：

```json
{
  "shoeCount": 3,
  "shoeCost": 2879,
  "totalCost": 2879,
  "shoeCostText": "2879.00",
  "totalCostText": "2879.00"
}
```

统计口径：

- `shoeCount`：当前用户未删除球鞋总数（含已退役，不含逻辑删除）。
- `shoeCost`：当前用户未删除球鞋的 `purchase_price` 合计，空值按 `0` 处理。
- `totalCost`：当前阶段等于 `shoeCost`；若后续引入磨损/更换记录再累加。

> 注意：该接口只统计球鞋自身；是否将球鞋购买费用并入首页 `expense`（新增 `shoeCost` 或扩展 `racketCost` 语义）属于开放问题，见 12.2，需用户确认后再调整。

## 8. 校验规则汇总

- `name` 必填，去除首尾空格后不能为空。
- `status` 创建时只接受 `1`/`2`（默认 `2`），编辑时接受 `1`/`2`/`3`。
- `purchaseDate` 可选，格式 `YYYY-MM-DD`，解析失败返回 `40001`。
- `purchasePrice` 可选，建议校验非负。
- `size`、`colorway` 可选，做长度与内容安全校验（复用球拍的内容安全服务）。
- 所有操作只允许当前用户自己的球鞋，通过 JWT 用户 ID 校验归属。
- 软删除数据不参与任何查询。

## 9. 错误码

复用现有统一错误码，不新增：

| code | message | 场景 |
|---:|---|---|
| `40001` | invalid request | 参数缺失/格式错误、`libraryId` 不存在、创建传 `status=3` |
| `40101` | unauthorized | 未登录或 token 无效 |
| `40301` | forbidden | 操作其他用户的球鞋（正常情况下由归属校验拦截，不会出现） |
| `40401` | not found | 球鞋不存在或已删除 |
| `50001` | internal error | 服务内部错误 |

## 10. 路由注册清单

在 `cmd/api/main.go` 的 `authed` 分组下新增（与球拍顺序对齐）：

```go
authed.GET("/shoe-brands", shoeHandler.Brands)
authed.GET("/shoe-series", shoeHandler.Series)
authed.GET("/shoe-library", shoeHandler.Library)
authed.GET("/shoe-library/stats", shoeHandler.LibraryStats)
authed.GET("/my-shoes", shoeHandler.MyShoes)
authed.GET("/my-shoes/primary", shoeHandler.MyPrimaryShoe)
authed.GET("/shoes/stats", shoeHandler.Stats)
authed.GET("/shoes", shoeHandler.List)
authed.POST("/shoes", shoeHandler.Create)
authed.GET("/shoes/selectable", shoeHandler.Selectable)
authed.GET("/shoes/:id", shoeHandler.Detail)
authed.PUT("/shoes/:id", shoeHandler.Update)
authed.DELETE("/shoes/:id", shoeHandler.Delete)
authed.POST("/shoes/:id/set-primary", shoeHandler.SetPrimary)
authed.POST("/shoes/:id/retire", shoeHandler.Retire)
```

## 11. 实施顺序与检查项

建议按 AGENTS.md 第 19 节顺序实施：

1. `migrations/015_create_shoe_tables.sql`（已确认的 3 张库表 + `my_shoes` 我的球鞋表）。
2. `internal/model/shoe.go`：`Shoe`、`ShoeResponse`、`CreateShoeRequest`、`UpdateShoeRequest`、库相关 DTO。
3. `internal/repository/shoe_repository.go`：品牌/系列/库查询、CRUD、`SetPrimary`、`Retire`、`SoftDelete`、`Stats`。
4. `internal/service/shoe_service.go`：库补全、状态规则、内容安全、日期解析、响应组装。
5. `internal/handler/shoe_handler.go`：参数绑定、当前用户、统一响应、日志。
6. `cmd/api/main.go` 注册路由。
7. 更新 `docs/api.md`（新增球鞋章节）、`docs/change-log.md`。
8. 执行 `gofmt` 和 `go test ./...`。

## 12. 开放问题（需要用户确认）

### 12.1 打球记录是否改关联球鞋 ID

已实现（2026-08-05）：`tennis_sessions` 新增 `shoe_id`，`POST/PUT /api/sessions` 接收 `shoeId`；`SessionResponse.shoeName` 按 `shoeId` 关联我的球鞋实时返回，`shoeId = 0` 时回退快照 `shoe_name`。

说明：为保留旧客户端兼容，`shoe_name` 快照列保留（与球拍 `racket_name` 直接删除不同），`shoeId > 0` 时以实时关联为准。

### 12.2 球鞋费用是否计入首页/统计

已确认并实现（2026-08-05）：首页 `expense` 与统计接口新增 `shoeCost` 字段，`totalCost = sessionCost + racketCost + stringingCost + shoeCost`；`racketCost` 语义保持不变。

### 12.3 球鞋库数据来源

`shoe_library` 与球拍库一样是系统维护数据，且多了一个 `status`（`draft / published / archived`）发布状态。需确认：

- 初期库数据由迁移/脚本 seed 还是运营手工维护；
- 数据录入后的发布流程：谁负责把 `draft` 改为 `published`（建议由数据维护方确认后再发布，C 端只读 `published`）；
- 是否需要一个球鞋库管理接口（本次设计不提供，球拍库也没有）。

### 12.4 尺码体系

性别（`gender`：`0` 未知 / `1` 男 / `2` 女 / `3` 童）已按确认的库表落地；场地类型（`surface`）已作为库字段存在。

尺码体系仍需确认：建议统一使用 EU 码数字字符串（`42`、`42.5`）存入 `shoe.size`，不做 US/UK 换算；如需多体系，再单独设计。

### 12.5 是否需要磨损/更换记录

球鞋没有与穿线对应的记录类型。如需要"磨损程度/更换日期"等记录，请确认后另行设计，本次不包含。

### 12.6 我的球鞋表 `my_shoes` 的 DDL

已确认（2026-08-05）："我的球鞋"表名定为 `my_shoes`（不沿用单数 `shoe`），DDL 已追加到 `migrations/015_create_shoe_tables.sql`，并补充 `PRIMARY KEY (id)`。

- `shoe_library.price`（当前售价）是否作为添加时 `purchasePrice` 的后端兜底默认值（目前设计为前端默认带入，后端不强制）。

> 已实现（2026-08-05）：`my_shoes` 表已按确认结果追加到 `migrations/015_create_shoe_tables.sql` 并完成建表。
