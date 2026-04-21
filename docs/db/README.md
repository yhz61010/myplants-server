# 数据库设计文档

## 概述
本文档描述了 MyPlants Server 应用程序的数据库模式，使用 SQLite 和 GORM ORM。

## 模型

### Content
表示统一的內容条目（日记或植物）。

**字段：**
- `ID` (uint): 主键
- `Type` (string): 内容类型（"diary" 或 "plant"）
- `UserID` (string): 用户标识符
- `Title` (string): 统一标题/名称
- `Content` (string): 统一内容/描述（文本）
- `ImagesStr` (string): 图片 URL 的 JSON 字符串
- `TagsStr` (string): 标签的 JSON 字符串（包含 family_（科）和 genus_（属）前缀）
- `Images` ([]string): 图片的计算字段（不存储）
- `Tags` ([]string): 标签的计算字段（不存储）
- `CreatedAt` (time.Time): 创建时间戳
- `UpdatedAt` (time.Time): 更新时间戳
- `DeletedAt` (gorm.DeletedAt): 软删除时间戳

**注意：**
- 通过 Type 字段区分内容类型。
- 图片和标签以 JSON 字符串形式存储，以兼容 SQLite。
- 计算字段在数据库查询后填充。
- 植物的科和属信息通过 Tags 中的 "family_" 和 "genus_" 前缀区分。

### User
表示用户账户。

**字段：**
- `ID` (uint): 主键
- `Username` (string): 用户名（唯一）
- `Password` (string): 密码哈希（不暴露在 JSON 中）
- `Avatar` (string): 头像 URL
- `Bio` (string): 个人简介（长文本）
- `CreatedAt` (time.Time): 创建时间戳
- `UpdatedAt` (time.Time): 更新时间戳
- `DeletedAt` (gorm.DeletedAt): 软删除时间戳

**注意：**
- 密码存储为哈希值，不在 API 响应中返回。
- 用户名是唯一的。
- **引擎**: SQLite
- **ORM**: GORM v1.31.1
- **Driver**: gorm.io/driver/sqlite v1.6.0
- **File**: `myplants.db`

## 迁移
数据库模式在应用程序启动时使用 GORM 的 AutoMigrate 功能自动迁移。

## 关系
目前，模型之间没有定义明确的关联。用户信息通过字符串 ID 引用。

## 未来考虑
- 添加用户模型以进行适当的用户管理
- 实现日记和植物模型之间的关联
- 添加索引以优化性能
- 考虑添加验证约束