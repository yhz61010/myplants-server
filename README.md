# MyPlants Server

## 项目概述
MyPlants Server 是一个轻量级后端服务，面向植物知识展示与栽培经验分享。服务采用 Gin 实现 RESTful API，使用 GORM + SQLite 作为数据存储，适合部署在 1GB 内存的 VPS 上。

核心功能：用户资料管理、植物知识库（上传/查询）、以及类似时间线的栽培分享（含分页）。

开发约定：所有接口使用 JSON/UTF-8，使用 JWT 做鉴权；图片采用外链（对象存储），数据库仅保存 URL；代码采用模块化（`cmd/`, `internal/`, `pkg/`）以便于维护与扩展。

## 项目结构
以下是针对“植物分享平台”项目的目录结构规划，基于 Go 最佳实践（如使用 `cmd/` 作为可执行文件入口，`internal/` 存放私有代码，`pkg/` 存放可重用代码）。结构设计简洁、模块化，便于维护和扩展。每个文件夹和文件都附带简述功能。

```
myplants-server/
├── cmd/                          # 可执行文件目录
│   └── myplants-server/          # 主应用程序
│       └── main.go               # 应用程序入口点，初始化服务器和路由
├── internal/                     # 私有代码目录（避免外部包导入）
│   ├── handlers/                 # HTTP 请求处理器
│   │   ├── diary.go              # 日记相关 API 处理器（如上传、查询日记）
│   │   └── user.go               # 用户相关 API 处理器（如登录、资料管理）
│   ├── models/                   # 数据模型定义
│   │   ├── diary.go              # 日记数据结构（GORM 模型）
│   │   └── user.go               # 用户数据结构（GORM 模型）
│   ├── database/                 # 数据库相关
│   │   └── db.go                 # 数据库连接、迁移和初始化（SQLite + GORM）
│   ├── middleware/               # 中间件组件
│   │   ├── auth.go               # 认证中间件（JWT 或会话）
│   │   └── cors.go               # CORS 处理中间件
│   └── routes/                   # 路由定义
│       └── routes.go             # API 路由注册和分组
├── pkg/                          # 可重用公共包
│   ├── config/                   # 配置管理
│   │   └── config.go             # 配置文件加载和环境变量处理
│   └── utils/                    # 工具函数
│       └── helpers.go            # 通用辅助函数（如时间格式化、错误处理）
├── tests/                        # 测试文件
│   ├── handlers_test.go          # 处理器单元测试
│   └── models_test.go            # 模型单元测试
├── docs/                         # 文档
│   └── README.md                 # 项目说明文档
├── scripts/                      # 脚本
│   └── build.sh                  # 构建脚本
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖校验
└── .gitignore                    # Git 忽略文件
```

### 说明：
- **cmd/**: 存放主程序入口，便于构建多个可执行文件（如未来添加 CLI 工具）。
- **internal/**: 核心业务逻辑，防止外部依赖。
- **pkg/**: 公共库，可被其他项目复用。
- **tests/**: 单元测试，确保代码质量。
- **docs/** 和 **scripts/**: 辅助文档和自动化脚本。
- 结构遵循 Go 社区约定（如 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)），适合小型到中型项目。如果项目复杂，可进一步拆分（如添加 `api/` 用于 OpenAPI 定义）。

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

运行（示例）：

```bash
export JWT_SECRET="your_super_secret"
export UPYUN_BUCKET="myplants"
export UPYUN_OPERATOR="op"
export UPYUN_PASSWORD="password"
GIN_MODE=release go run cmd/myplants-server/main.go
```

注意：务必将 `JWT_SECRET`、`UPYUN_PASSWORD` 等敏感信息妥善保管在环境变量或机密管理系统中，不要提交到版本控制。
