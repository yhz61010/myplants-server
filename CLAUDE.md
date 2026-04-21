# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 项目概述

MyPlants Server — 植物知识与栽培分享平台的轻量级 Go 后端。基于 Gin REST API + GORM/SQLite，面向 1核/1GB 内存 VPS 部署。

## 主要功能
- 用户资料管理（显示和修改）
- 植物知识上传和显示系统（包括科学分类、描述、图片等）
- 社交分享功能（类似动态/时间线），用于用户栽培日志和经验分享

## 开发约定
- 遵循标准的 Go 命名约定和项目结构
- 使用 Gin 进行 HTTP 路由和中间件
- 实现 RESTful API 端点
- 使用 JWT 做鉴权
- 确保正确的错误处理和日志记录

## 常用命令

```bash
# 构建
go build -o bin/myplants-server ./cmd/myplants-server
# 或: scripts/build.sh

# 开发运行
go run ./cmd/myplants-server/main.go

# 运行全部测试
go test ./tests/...

# 运行单个测试
go test ./tests/ -run TestRegisterAndLogin -v

# 整理依赖
go mod tidy
```

**运行时所需环境变量**：`JWT_SECRET`；图片上传还需要：`UPYUN_BUCKET`、`UPYUN_OPERATOR`、`UPYUN_PASSWORD`。未设置 `JWT_SECRET` 时，开发环境回退使用 `"dev-secret"`。

## 架构

- **internal/routes/routes.go** — 所有路由注册。三组路由：公开认证路由（`/api/auth`）、受保护 API（`/api`，使用 JWT 中间件）、管理员路由（`/api/admin`，使用管理员中间件）。同时在 `/admin` 下提供管理后台 HTML 页面。
- **internal/auth/auth.go** — JWT 签发与解析（HS256，`golang-jwt/jwt/v5`）。密钥在 init 时从 `JWT_SECRET` 环境变量加载。
- **internal/middleware/auth.go** — JWT 认证中间件，提取 Bearer token，将 `userId`/`username` 设置到 Gin 上下文。
- **internal/handlers/** — `user.go`（认证 + 用户 CRUD + 管理员用户管理）、`content.go`（日记/植物 CRUD，含分页）、`upload.go`（又拍云流式上传）。
- **internal/models/** — GORM 模型：`User` 和 `Content`（统一模型，通过 `Type` 字段区分日记/植物）。
- **internal/database/db.go** — SQLite 初始化及 AutoMigrate。连接池限制为 1 个连接（适配低内存 VPS）。全局 `database.DB` 供各处使用。
- **pkg/config/** — 配置辅助工具（预留）。
- **templates/** — 管理后台 HTML/静态资源，由 Gin 提供服务。
- **tests/** — 基于内存 SQLite（`setupInMemoryDB`）和 `httptest` 的 handler 级别测试。

## 关键约定

- 所有 API 响应为 JSON/UTF-8。图片存储为 URL（又拍云对象存储），数据库仅保存 URL 字符串。
- Content 模型是统一的：日记和植物共用同一张 `Content` 表，通过 `Type` 字段区分。日记专用端点是按类型过滤的别名路由。
- 数据库文件 `myplants.db` 在工作目录下创建——务必在项目根目录运行。
- Handler 通过全局 `database.DB`（或 `database.GetDB()`）访问数据库，不使用依赖注入。

## 核心技术栈约束

- **Web 框架**：Gin（生产环境必须使用 `gin.ReleaseMode`）
- **ORM**：GORM
- **数据库**：SQLite（单文件数据库，所有代码必须针对其优化）
- **部署环境**：1GB RAM / 1-Core CPU（严禁引入 Redis、Kafka 等额外中间件）
- **图片处理**：仅限通过又拍云 (UpYun) SDK 流式上传，数据库只存 URL

## 编程行为准则

### 先思后行
- **显式假设**：输出代码前，必须先列出假设
- **消除歧义**：需求不明确时，必须停止并提问，禁止凭空猜测
- **架构审查**：如果要求会导致内存占用过高或过于复杂，必须提出异议并给出极简替代方案

### 极简主义
- 严禁编写当前任务之外的"扩展性"代码
- 50 行能解决的，绝不写 200 行
- 逻辑重复不超过 3 次，不做函数抽取或接口抽象
- 优先使用流式处理，避免将大文件或大数据集一次性加载进内存

### 手术刀式修改
- 只修改与当前任务直接相关的行
- 严禁自动重构邻近代码、修改格式或更新注释
- 严格匹配现有项目的代码命名和结构风格

### 目标驱动与小步快跑
- 编码前按以下格式提供计划：
  1. [步骤 A] → verify: [如何验证]
  2. [步骤 B] → verify: [检查点]
- 逻辑较长时，自动拆分为多个步骤，确保每段代码输出都是完整的

## 内存控制专项（1GB RAM）

- 数据库连接池限制：`SetMaxOpenConns(1)`
- 所有 `List` 接口必须强制分页
- 严禁使用后台持续占用大量 CPU/内存的定时任务或缓存预热逻辑
