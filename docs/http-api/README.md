# HTTP API 文档

本文档列出 MyPlants Server 的 HTTP 接口说明（中文）。所有受保护的接口均使用 JWT Bearer Token 进行认证。

## 概览
- 默认返回内容分页大小建议：10-20 条，避免一次性返回大量数据导致 VPS 内存压力。
- 所有请求与响应使用 UTF-8 编码，JSON 格式。
- 图片处理：接口仅保存图片 URL，建议将图片上传到对象存储（Cloudflare R2、阿里云 OSS 等），数据库仅保存外链 URL。

---

## 认证

### 注册用户
- URL：`POST /api/auth/register`
- 描述：创建新用户账号
- 请求体（JSON）：
  ```json
  {
    "username": "string (required)",
    "password": "string (required, min 6)",
    "avatar": "string (optional, URL)",
    "bio": "string (optional)"
  }
  ```
- 响应：
  - 201 Created：返回创建的用户资源（不包含密码）
  - 400 Bad Request：参数校验失败
  - 409 Conflict：用户名已存在

### 登录
- URL：`POST /api/auth/login`
- 描述：用户登录，返回 JWT Token
- 请求体（JSON）：
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- 响应：
  - 200 OK：
    ```json
    {
      "token": "<jwt-token>",
      "user": {"id": 1, "username": "...", "avatar": "...", "bio": "..."}
    }
    ```
  - 401 Unauthorized：凭据错误

---

## 内容管理（统一 `Content` 表：日记 / 植物）
> 说明：`Content` 使用 `Type` 字段区分 `"diary"` 或 `"plant"`。植物的科/属信息以标签形式存入 `tags`，使用 `family_` 和 `genus_` 前缀区分。

### 创建内容（发布日记或植物）
- URL：`POST /api/diaries`
- 描述：创建一条内容（日记或植物条目）
- 认证：需要 `Authorization: Bearer <token>`
- 请求体（JSON）：
  ```json
  {
    "type": "diary" | "plant",
    "userId": "string",
    "title": "string",
    "content": "string",
    "images": ["string"],
    "tags": ["string"],
    "createTime": "string (optional, RFC3339)"
  }
  ```
- 响应：
  - 201 Created：返回创建的内容对象
  - 400 Bad Request：参数错误或时间格式不正确
  - 401 Unauthorized：未认证

### 按关键字模糊查询植物列表（植物知识库）
- URL：`GET /api/plants?query=<q>&limit=<n>&offset=<m>`
- 描述：按 `title` 或 `content` 模糊匹配植物条目，支持分页
- 参数：
  - `query` (可选)：搜索关键字
  - `limit` (可选)：每页数量，建议默认 10，最大不超过 50
  - `offset` (可选)：偏移量
- 响应：
  - 200 OK：
    ```json
    {
      "items": [ /* content objects */ ],
      "total": 123
    }
    ```

### 时间线动态（分页）
- URL：`GET /api/timeline?limit=<n>&offset=<m>`
- 描述：获取所有用户的栽培分享（按创建时间倒序），必须支持分页
- 参数：同上（`limit`/`offset`），默认 `limit=10`
- 响应：同上

---

## 图片上传建议
- 不要将图片二进制存入 SQLite。请实现独立的图片上传接口：
  1. 接收文件（multipart/form-data）
  2. 将文件上传到对象存储（R2/OSS 等）
  3. 返回可访问的 URL，前端将该 URL 放入 `images` 字段提交给内容创建接口

示例上传接口：
- `POST /api/upload`（返回 `{ "url": "https://..." }`）

### 又拍云 (UpYun) 集成说明

- 我们在服务端集成了又拍云 Go SDK（`github.com/upyun/go-sdk/v3/upyun`）。
- 环境变量：
  - `UPYUN_BUCKET`：又拍云空间名（例如 `myplants`）
  - `UPYUN_OPERATOR`：操作员用户名
  - `UPYUN_PASSWORD`：操作员密码

- 行为：当上述三个环境变量都存在时，`POST /api/upload` 会使用 UpYun SDK 将上传的文件写入又拍云，返回的公有访问 URL 格式为 `https://myplants.leovp.com/<保存路径>`（域名需在又拍云控制台或 DNS 配置中绑定并启用 HTTPS）。
- 对于大文件（>10MB）建议使用本地临时文件并通过 SDK 的 `LocalPath` 参数上传，以启用断点续传与 MD5 校验；当前实现会优先尝试使用流式 `Reader`，如需改为临时文件上传我可以修改。

示例：

请求（multipart/form-data）

```
POST /api/upload
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary...

file=<binary file>
```

响应

```json
{ "url": "https://myplants.leovp.com/20260421/abcdef123456_image.jpg" }
```

注意：返回的 URL 假定你的自定义域 `myplants.leovp.com` 已正确指向又拍云或静态文件所在位置并配置了 HTTPS；如果只是用于本地测试，服务器会把文件保存到 `./uploads/`（同样返回一个 `https://myplants.leovp.com/...` 风格的 URL，但你需在部署时将域名指向实际存储或配置静态文件服务）。

---

## 错误码与响应规范
- 遵循 HTTP 标准状态码
- 错误响应统一格式：
  ```json
  { "error": "human-readable message" }
  ```

---

## 安全与部署建议
- 将 `JWT` 密钥和敏感配置放入环境变量，而不是硬编码
- 部署时设置 `GIN` 的 Release 模式：在 `main` 中调用 `gin.SetMode(gin.ReleaseMode)`
- 对列表接口设置合理的 `limit` 上限，防止 OOM

---

如需将更多接口（如用户资料更新、关注/点赞、评论等）加入文档，我可以继续补充。