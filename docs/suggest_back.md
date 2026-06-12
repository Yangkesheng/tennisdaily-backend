# 日历页后端接口设计

## 1. 背景与目标

当前小程序日历页面只展示当月日历，不再展示整年 12 个月。

现有前端实现仍可能通过：

```http
GET /api/sessions
```

获取全部打球记录，再由前端筛选当月日期。这种方式在记录增多后会造成不必要的数据传输和前端处理。

因此建议新增一个轻量级后端接口，专门为日历页提供某个月的打球日期数据。

目标：

- 按月查询当前用户的打球日期
- 返回当月有记录的日期列表
- 返回每天记录数量
- 返回当月有记录的天数
- 不返回完整打球记录详情
- 保持前端日历渲染逻辑简单

---

## 2. 接口定义

### Method

```http
GET
```

### Path

```http
/api/sessions/calendar
```

### 完整示例

```http
GET /api/sessions/calendar?year=2026&month=6
Authorization: Bearer <token>
```

---

## 3. 请求参数

### Query Parameters

| 参数 | 类型 | 必填 | 示例 | 说明 |
|---|---|---:|---|---|
| `year` | number | 是 | `2026` | 年份 |
| `month` | number | 是 | `6` | 月份，范围 `1-12` |

### 参数规则

- `year` 建议限制范围：`2000-2100`
- `month` 必须在 `1-12` 之间
- 日期统计按自然月计算
- 时区建议使用 `Asia/Shanghai`

---

## 4. 鉴权规则

该接口需要登录。

前端请求时必须携带：

```http
Authorization: Bearer <token>
```

后端从 JWT 中解析当前用户 ID，只返回当前用户自己的打球记录日历数据。

未登录或 token 无效时返回：

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

---

## 5. 响应结构

沿用项目统一响应格式：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### data 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `year` | number | 查询年份 |
| `month` | number | 查询月份 |
| `activeDayCount` | number | 当月有打球记录的天数 |
| `days` | array | 当月有记录的日期列表 |

### days[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| `date` | string | 日期，格式 `YYYY-MM-DD` |
| `count` | number | 当天打球记录数量 |

---

## 6. 成功响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "year": 2026,
    "month": 6,
    "activeDayCount": 4,
    "days": [
      {
        "date": "2026-06-02",
        "count": 1
      },
      {
        "date": "2026-06-08",
        "count": 2
      },
      {
        "date": "2026-06-15",
        "count": 1
      },
      {
        "date": "2026-06-29",
        "count": 1
      }
    ]
  }
}
```

---

## 7. 空数据响应示例

如果当月没有任何打球记录，返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "year": 2026,
    "month": 6,
    "activeDayCount": 0,
    "days": []
  }
}
```

---

## 8. 错误响应示例

### 8.1 参数缺失

请求：

```http
GET /api/sessions/calendar
```

响应：

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

### 8.2 月份非法

请求：

```http
GET /api/sessions/calendar?year=2026&month=13
```

响应：

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

### 8.3 未登录或 token 无效

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

### 8.4 服务端异常

```json
{
  "code": 500,
  "message": "internal error",
  "data": null
}
```

---

## 9. 前端 TypeScript 类型

建议放在：

```text
miniprogram/models/session.ts
```

```ts
export interface SessionCalendarDay {
  date: string
  count: number
}

export interface SessionCalendar {
  year: number
  month: number
  activeDayCount: number
  days: SessionCalendarDay[]
}
```

---

## 10. 前端 Service 建议

建议在：

```text
miniprogram/services/session-api-service.ts
```

新增：

```ts
export const getSessionCalendarFromApi = async (
  year: number,
  month: number,
): Promise<SessionCalendar> => {
  return request<SessionCalendar>({
    url: `/api/sessions/calendar?year=${year}&month=${month}`,
  })
}
```

页面使用方式：

```ts
const result = await getSessionCalendarFromApi(currentYear, currentMonth)

this.setData({
  currentYear: result.year,
  currentMonth: result.month,
  activeDayCount: result.activeDayCount,
  calendarMonths: [
    createCalendarMonth(
      result.year,
      result.month,
      result.days.map((day) => day.date),
    ),
  ],
})
```

---

## 11. 后端 Go DTO 设计

```go
type SessionCalendarQuery struct {
	Year  int `form:"year" binding:"required,min=2000,max=2100"`
	Month int `form:"month" binding:"required,min=1,max=12"`
}

type SessionCalendarDayResponse struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type SessionCalendarResponse struct {
	Year           int                          `json:"year"`
	Month          int                          `json:"month"`
	ActiveDayCount int                          `json:"activeDayCount"`
	Days           []SessionCalendarDayResponse `json:"days"`
}
```

---

## 12. 数据查询逻辑

### 数据表

```text
tennis_sessions
```

### 查询条件

```sql
user_id = 当前登录用户 ID
deleted_at IS NULL
date >= 当月第一天
date < 下月第一天
```

### SQL 示例

```sql
SELECT date, COUNT(*) AS count
FROM tennis_sessions
WHERE user_id = ?
  AND deleted_at IS NULL
  AND date >= ?
  AND date < ?
GROUP BY date
ORDER BY date ASC;
```

---

## 13. Repository 示例

```go
type SessionCalendarDay struct {
	Date  time.Time
	Count int64
}

func (r *SessionRepository) ListCalendarDays(
	ctx context.Context,
	userID uint,
	startDate time.Time,
	endDate time.Time,
) ([]SessionCalendarDay, error) {
	var rows []SessionCalendarDay

	err := r.db.WithContext(ctx).
		Table("tennis_sessions").
		Select("date, COUNT(*) AS count").
		Where(
			"user_id = ? AND deleted_at IS NULL AND date >= ? AND date < ?",
			userID,
			startDate,
			endDate,
		).
		Group("date").
		Order("date ASC").
		Scan(&rows).Error

	return rows, err
}
```

---

## 14. Service 示例

```go
func (s *SessionService) GetCalendar(
	ctx context.Context,
	userID uint,
	year int,
	month int,
) (*SessionCalendarResponse, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 1, 0)

	rows, err := s.repo.ListCalendarDays(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	days := make([]SessionCalendarDayResponse, 0, len(rows))
	for _, row := range rows {
		days = append(days, SessionCalendarDayResponse{
			Date:  row.Date.In(loc).Format("2006-01-02"),
			Count: row.Count,
		})
	}

	return &SessionCalendarResponse{
		Year:           year,
		Month:          month,
		ActiveDayCount: len(days),
		Days:           days,
	}, nil
}
```

---

## 15. Handler 示例

```go
func (h *SessionHandler) GetCalendar(c *gin.Context) {
	var query SessionCalendarQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, response.CodeInvalidRequest, "invalid request")
		return
	}

	userID := middleware.GetUserID(c)

	result, err := h.sessionService.GetCalendar(
		c.Request.Context(),
		userID,
		query.Year,
		query.Month,
	)
	if err != nil {
		response.Fail(c, response.CodeInternalError, "internal error")
		return
	}

	response.Success(c, result)
}
```

---

## 16. 路由注册

需要注意：

```http
GET /api/sessions/calendar
```

必须注册在：

```http
GET /api/sessions/:id
```

之前。

否则 `calendar` 可能会被 Gin 识别为 `:id`。

推荐顺序：

```go
sessions.GET("/calendar", sessionHandler.GetCalendar)
sessions.GET("", sessionHandler.List)
sessions.GET("/latest", sessionHandler.Latest)
sessions.GET("/:id", sessionHandler.GetByID)
```

或者：

```go
api.GET("/sessions/calendar", authMiddleware, sessionHandler.GetCalendar)
api.GET("/sessions/:id", authMiddleware, sessionHandler.GetByID)
```

---

## 17. 时区规则

前端日期格式统一为：

```text
YYYY-MM-DD
```

后端统计自然月时建议统一使用：

```text
Asia/Shanghai
```

月份区间：

```go
loc, _ := time.LoadLocation("Asia/Shanghai")

startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
endDate := startDate.AddDate(0, 1, 0)
```

查询范围采用左闭右开：

```text
[startDate, endDate)
```

也就是：

```sql
date >= startDate
AND date < endDate
```

---

## 18. 为什么不返回整月所有日期

不建议后端返回 1-31 号全部日期。

原因：

- 前端更适合生成日历格子
- 前端需要补齐月初空白格
- 前端需要标记今天
- 前端负责日历布局
- 后端只需要返回有记录的日期即可

后端返回：

```text
哪些日期有记录
每个日期几条记录
```

即可满足日历页需求。

---

## 19. 最终推荐协议

### Request

```http
GET /api/sessions/calendar?year=2026&month=6
Authorization: Bearer <token>
```

### Response

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "year": 2026,
    "month": 6,
    "activeDayCount": 4,
    "days": [
      {
        "date": "2026-06-02",
        "count": 1
      },
      {
        "date": "2026-06-08",
        "count": 2
      }
    ]
  }
}
```
