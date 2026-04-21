# Coplit 生成服务端聊天全内容

`/init` 这是一个介绍和学习植物知识，以及分享自己的养植和心得的 App 的服务端。我想用 Gin 作为服务端框架进行开发。目前主要功能包括，用户信息展示和修改。植物知识上传和显示，类似朋友圈这样的功能，可以让用户分享自己的日志。由于系统访问量不大，DB 采用 SQLite (配合 GORM)即可。

## 项目目录结构规划
我要开发植物分享平台，请帮我规划项目目录结构（符合 Go 最佳实践）。只需要列出文件夹和文件名，并简述功能，不要写具体代码。

## 开始生成目录结构
目录结构没有问题。开始生成，并将上述内容更新到文档。

## 业务功能实现
###
请帮我设计这个植物分享平台的 GORM 模型（SQLite）。需要包含：
Plant: 植物名称、科学分类（按植物的科和属进行分类）、内容介绍（长文本）、上传者、上传时间、多张图片 URL。
并将 DB 相关文档，放到项目目录下 docs 目录的 db 目录上。

### 
再帮我实现用户的 GORM 模型。包括：
User: 用户名、密码、头像、个人简介。

###
我想将 #sym:Plant 和 #sym:Diary 合成一张表，便于统一管理。

###
#sym:Content 删除 Name 字段，统一使用 Title；删除 Description 字段，统一使用 Content。删除 Family 和 Genus 字段，统一使用 TagsStr 管理，通过 “family_” 和 “genus_” 前缀来区分。

###
请基于 Gin 框架实现用户注册和登录接口。
使用 bcrypt 对密码进行加密。
登录成功后返回 JWT Token。
编写一个 AuthMiddleware 中间件，用于验证受保护的路由。

###
将 HTTP 文档，更新到 docs 下的 http-api 目录下。

###
根据项目实际情况，添加 .gitignore 内容。

### 
将所有代码提交到我的 GitHub 上，创建 myplants-server 仓库。我的 GitHub 地址：https://github.com/yhz61010 

###
请实现以下核心接口的 CRUD：
植物日记（对应 Content 模型）: 支持按关键字模糊查询标题和标签。获取所有用户分享的日记，必须支持分页功能（Limit/Offset），防止一次读取过多数据导致 VPS 内存溢出。

###
图片处理（外链化）。千万不要将用户上传的图片二进制数据直接存进 SQLite 数据库。写一个接口，将图片上传到像“又拍云”云存储上，对应的存储名字叫 myplants，具体域名为 myplants.leovp.com。该域名已支付 HTTPS。数据库里只存一个 URL 字符串。

###
部署到 VPS 时，确保在代码中加入(释放模式（Release Mode）)：
gin.SetMode(gin.ReleaseMode)