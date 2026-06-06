# TennisDaily 登录接口文档

> 本文档单独描述小程序登录、当前用户、退出登录、用户资料相关接口。业务接口仍统一使用 `Authorization: Bearer <token>` 鉴权。

---

## 1. 基础信息

### 1.1 Base URL

本地默认：

```text
http://localhost:8081
```

真机调试请使用电脑局域网 IP，例如：

```text
http://192.168.x.x:8081
```

### 1.2 Content-Type

有请求体的接口均使用：

```http
Content-Type: application/json
```

### 1.3 统一响应格式

成功：

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

失败：

```json
{
  "code": 40001,
  "message": "invalid request",
  "data": null
}
```

### 1.4 鉴权方式

登录成功后，后端返回 JWT token。后续需要登录态的接口请求头携带：

```http
Authorization: Bearer <token>
```

---

## 2. 登录流程

### 2.1 首次手机号登录

```text
用户进入登录页
-> 勾选协议
-> 点击手机号快捷登录
-> 小程序 getPhoneNumber 获取 phoneCode
-> 小程序 wx.login 获取 loginCode
-> POST /api/auth/phone-login
-> 后端获取 openid 和手机号
-> 后端查找或创建用户
-> 后端绑定微信身份和手机号
-> 后端签发 JWT
-> 前端保存 token
-> 进入首页
```

### 2.2 后续打开小程序

```text
小程序读取本地 token
-> GET /api/auth/me
-> token 有效：进入首页
-> token 无效：清除 token，回登录页
```

### 2.3 业务接口鉴权

所有业务接口继续使用：

```http
Authorization: Bearer <token>
```

---

## 3. 接口清单

| 方法 | 路径 | 鉴权 | 说明 |
|---|---|---|---|
| POST | `/api/auth/phone-login` | 否 | 手机号快捷登录 |
| GET | `/api/auth/me` | 是 | 获取当前用户、校验 token |
| POST | `/api/auth/logout` | 是 | 退出登录，MVP 阶段服务端直接返回成功 |
| PUT | `/api/auth/profile` | 是 | 更新用户昵称、头像 |
| POST | `/api/auth/wechat-login` | 否 | 旧版微信静默登录，兼容保留 |

---

## 4. 手机号快捷登录

### 4.1 POST /api/auth/phone-login

用于小程序手机号快捷登录。

后端完成：

- 校验 `loginCode`
- 校验 `phoneCode`
- 获取微信 `openid`
- 获取手机号
- 查找或创建用户
- 绑定微信身份
- 签发 JWT

#### 请求

```http
POST /api/auth/phone-login
Content-Type: application/json
```

#### 请求体

```json
{
  "loginCode": "wx.login 返回的 code",
  "phoneCode": "getPhoneNumber 返回的 code"
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| loginCode | string | 是 | `wx.login()` 返回的临时登录凭证 |
| phoneCode | string | 是 | `button open-type="getPhoneNumber"` 返回的手机号授权 code |

#### 响应 data

```json
{
  "token": "jwt-token",
  "user": {
    "id": 1,
    "phone": "13800138000",
    "maskedPhone": "138****8000",
    "nickname": "",
    "avatarUrl": "",
    "createdAt": "2026-05-30T12:00:00+08:00",
    "updatedAt": "2026-05-30T12:00:00+08:00"
  },
  "isNewUser": true
}
```

字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| token | string | 后端签发的 JWT |
| user | object | 当前登录用户 |
| isNewUser | boolean | 是否为新创建用户 |

user 字段说明：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | number | 用户 ID |
| phone | string | 手机号 |
| maskedPhone | string | 脱敏手机号，例如 `138****8000` |
| nickname | string | 用户昵称，当前可为空 |
| avatarUrl | string | 用户头像 URL，当前可为空 |
| createdAt | string | 创建时间，RFC3339 |
| updatedAt | string | 更新时间，RFC3339 |

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "phone": "13800138000",
      "maskedPhone": "138****8000",
      "nickname": "",
      "avatarUrl": "",
      "createdAt": "2026-05-30T12:00:00+08:00",
      "updatedAt": "2026-05-30T12:00:00+08:00"
    },
    "isNewUser": true
  }
}
```

#### curl 示例

```bash
curl http://localhost:8081/api/auth/phone-login \
  -X POST \
  -H 'Content-Type: application/json' \
  -d '{"loginCode":"wx-login-code","phoneCode":"phone-code"}'
```

---

## 5. 获取当前用户

### 5.1 GET /api/auth/me

用于：

- 校验 token 是否有效
- 小程序启动时判断是否需要进入登录页
- 我的页面展示账号状态
- 后续个人资料页面使用

#### 请求

```http
GET /api/auth/me
Authorization: Bearer <token>
```

#### 响应 data

```json
{
  "id": 1,
  "phone": "13800138000",
  "maskedPhone": "138****8000",
  "nickname": "",
  "avatarUrl": "",
  "createdAt": "2026-05-30T12:00:00+08:00",
  "updatedAt": "2026-05-30T12:00:00+08:00"
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "phone": "13800138000",
    "maskedPhone": "138****8000",
    "nickname": "",
    "avatarUrl": "",
    "createdAt": "2026-05-30T12:00:00+08:00",
    "updatedAt": "2026-05-30T12:00:00+08:00"
  }
}
```

#### token 无效响应

HTTP 状态码：

```text
401 Unauthorized
```

响应体：

```json
{
  "code": 40101,
  "message": "unauthorized",
  "data": null
}
```

#### curl 示例

```bash
curl http://localhost:8081/api/auth/me \
  -H 'Authorization: Bearer <token>'
```

---

## 6. 退出登录

### 6.1 POST /api/auth/logout

MVP 阶段服务端不维护 token 黑名单。前端退出登录时应清除本地 token，并跳转登录页。

#### 请求

```http
POST /api/auth/logout
Authorization: Bearer <token>
```

#### 响应 data

```json
{
  "success": true
}
```

#### 响应示例

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "success": true
  }
}
```

#### curl 示例

```bash
curl http://localhost:8081/api/auth/logout \
  -X POST \
  -H 'Authorization: Bearer <token>'
```

---

## 7. 更新用户资料

### 7.1 PUT /api/auth/profile

用于更新用户昵称、头像。

#### 请求

```http
PUT /api/auth/profile
Authorization: Bearer <token>
Content-Type: application/json
```

#### 请求体

```json
{
  "nickname": "阿卡门徒",
  "avatarUrl": "https://example.com/avatar.png"
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| nickname | string | 否 | 用户昵称 |
| avatarUrl | string | 否 | 用户头像 URL |

#### 响应 data

```json
{
  "id": 1,
  "phone": "13800138000",
  "maskedPhone": "138****8000",
  "nickname": "阿卡门徒",
  "avatarUrl": "https://example.com/avatar.png",
  "createdAt": "2026-05-30T12:00:00+08:00",
  "updatedAt": "2026-05-30T12:00:00+08:00"
}
```

#### curl 示例

```bash
curl http://localhost:8081/api/auth/profile \
  -X PUT \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"nickname":"阿卡门徒","avatarUrl":"https://example.com/avatar.png"}'
```

---

## 8. 旧版微信静默登录

### 8.1 POST /api/auth/wechat-login

该接口为旧版兼容接口。新登录页建议使用 `/api/auth/phone-login`。

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

---

## 9. 微信服务调用说明

### 9.1 使用 loginCode 获取 openid

后端调用微信接口：

```http
GET https://api.weixin.qq.com/sns/jscode2session
```

参数：

```text
appid=<小程序 appid>
secret=<小程序 secret>
js_code=<loginCode>
grant_type=authorization_code
```

### 9.2 使用 phoneCode 获取手机号

后端调用微信接口：

```http
POST https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=ACCESS_TOKEN
Content-Type: application/json
```

请求体：

```json
{
  "code": "phoneCode"
}
```

---

## 10. 用户查找与绑定规则

后端采用以下逻辑：

1. 根据 `appid + openid` 查找微信身份绑定。
   - 找到：获取 user，更新/绑定手机号，返回 token。
2. 根据手机号查找 user。
   - 找到：绑定当前 `appid + openid` 到该 user，返回 token。
3. 都没找到：
   - 创建新 user。
   - 保存手机号。
   - 创建微信身份绑定。
   - 返回 token。

---

## 11. 错误响应

| 场景 | HTTP | code | message |
|---|---:|---:|---|
| 参数缺失 | 400 | 40001 | 登录参数不完整 |
| loginCode 无效 | 400 | 40001 | 微信登录凭证无效，请重试 |
| phoneCode 无效 | 400 | 40001 | 手机号授权已失效，请重试 |
| 微信服务异常 | 502 | 50001 | 微信服务暂时不可用，请稍后重试 |
| token 无效 | 401 | 40101 | unauthorized |

---

## 12. 前端接入建议

### 12.1 登录页调用顺序

```text
1. 用户勾选协议
2. 点击手机号快捷登录按钮
3. 从 event.detail.code 获取 phoneCode
4. 调 wx.login 获取 loginCode
5. POST /api/auth/phone-login
6. 保存 token
7. 进入首页
```

### 12.2 保存 token

```ts
wx.setStorageSync('token', result.token)
```

### 12.3 后续请求携带 token

```http
Authorization: Bearer <token>
```

### 12.4 小程序启动校验登录态

```text
读取本地 token
-> 有 token：GET /api/auth/me
-> 成功：进入首页
-> 401：清除 token，进入登录页
```
