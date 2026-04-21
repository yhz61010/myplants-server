# MyPlants Server

## 项目概述
MyPlants Server 是一个轻量级后端服务，面向植物知识展示与栽培经验分享。服务采用 Gin 实现 RESTful API，使用 GORM + SQLite 作为数据存储，适合部署在 1GB 内存的 VPS 上。

核心功能：用户资料管理、植物知识库（上传/查询）、以及类似时间线的栽培分享（含分页）。

开发约定：所有接口使用 JSON/UTF-8，使用 JWT 做鉴权；图片采用外链（对象存储），数据库仅保存 URL；代码采用模块化（`cmd/`, `internal/`, `pkg/`）以便于维护与扩展。

## 项目结构

```
myplants-server/
├── cmd/                          # 可执行文件目录
│   └── myplants-server/
│       └── main.go               # 应用程序入口点，初始化服务器和路由
├── internal/                     # 私有代码目录
│   ├── auth/
│   │   └── auth.go               # JWT 签发与解析（HS256）
│   ├── database/
│   │   └── db.go                 # 数据库连接、迁移和初始化（SQLite + GORM）
│   ├── handlers/                 # HTTP 请求处理器
│   │   ├── content.go            # 内容相关 API 处理器（日记/植物的 CRUD，含分页）
│   │   ├── upload.go             # 图片上传处理器（又拍云/本地）
│   │   └── user.go               # 用户认证、CRUD 及管理员操作
│   ├── models/                   # 数据模型定义
│   │   ├── content.go            # 统一内容数据结构（GORM 模型，支持 diary/plant 类型）
│   │   └── user.go               # 用户数据结构（GORM 模型，含 IsAdmin 字段）
│   ├── middleware/               # 中间件组件
│   │   ├── auth.go               # JWT 认证中间件
│   │   └── cors.go               # CORS 处理中间件
│   └── routes/                   # 路由定义
│       └── routes.go             # API 路由注册和分组
├── pkg/                          # 可重用公共包
│   ├── config/
│   │   └── config.go             # 配置辅助（预留）
│   └── utils/
│       └── helpers.go            # 工具函数（预留）
├── templates/                    # 管理后台 HTML 页面
│   ├── index.html                # 管理面板首页
│   ├── login.html                # 管理员登录页
│   ├── users.html                # 用户管理列表
│   ├── user_detail.html          # 用户详情/编辑
│   ├── diaries.html              # 日记管理列表
│   └── diary_detail.html         # 日记详情/编辑
├── tests/                        # 测试文件
│   ├── diary_handlers_test.go    # 日记处理器测试
│   ├── user_handlers_test.go     # 用户处理器测试
│   ├── handlers_test.go          # 处理器测试（预留）
│   └── models_test.go            # 模型测试（预留）
├── docs/                         # 文档
│   ├── http-api/README.md        # HTTP API 文档
│   └── db/README.md              # 数据库设计文档
├── scripts/
│   └── build.sh                  # 构建脚本
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖校验
└── .gitignore                    # Git 忽略文件
```

### 说明：
- **cmd/**: 存放主程序入口。
- **internal/**: 核心业务逻辑，防止外部依赖。
- **pkg/**: 公共库，可被其他项目复用。
- **templates/**: 管理后台前端页面，由 Gin 静态文件服务提供。
- **tests/**: 基于内存 SQLite 和 httptest 的 handler 级别测试。
- **docs/**: HTTP API 文档和数据库设计文档。
- 结构遵循 Go 社区约定（如 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)）。

## API 文档
详细的 HTTP API 文档请参考： [docs/http-api/README.md](docs/http-api/README.md)

## 部署与环境变量
在生产部署到 VPS 时，建议按以下要求配置环境变量并以 Release 模式运行：

- `GIN_MODE=release`（`main.go` 已调用 `gin.SetMode(gin.ReleaseMode)`）
- `JWT_SECRET`：用于签发 JWT 的 HMAC 密钥（不要硬编码到代码中）
- 可选又拍云配置（用于图片上传）：
	- `UPYUN_BUCKET`（例如 `myplants`）
	- `UPYUN_OPERATOR`
	- `UPYUN_PASSWORD`

运行说明（示例）：

```bash
export JWT_SECRET="your_super_secret"
export UPYUN_BUCKET="myplants"
export UPYUN_OPERATOR="op"
export UPYUN_PASSWORD="password"
GIN_MODE=release go run cmd/myplants-server/main.go
```

备注：服务启动时会在项目根目录创建或使用 `myplants.db` 文件；如果你在不同工作目录启动服务，可能会看到不同的 SQLite 文件，因此建议在项目根目录执行运行命令以避免混淆。

注意：务必将 `JWT_SECRET`、`UPYUN_PASSWORD` 等敏感信息妥善保管在环境变量或机密管理系统中，不要提交到版本控制。
