---
name: Cloud dev environment
description: Claude Code runs on cloud server; user develops locally and pushes to GitHub; some local files like cmd/myplants-server/main.go may not be in the repo
type: project
---

当前 Claude Code 运行在云端服务器上，代码来自 GitHub 仓库。用户在本地开发，部分文件（如 `cmd/myplants-server/main.go`）可能存在于本地但未提交到仓库。

**Why:** 用户说明了自己是本地开发，Claude Code 是云端代码，两者看到的文件可能不同。
**How to apply:** 如果发现关键文件在仓库中缺失，先确认是否为未提交的本地文件，而非不存在。不要轻易标注文件为"待创建"。
