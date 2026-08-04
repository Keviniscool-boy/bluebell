# Bluebell

<div align="center">

**前后端分离的社区论坛项目** —— Go + Gin 后端、Vue 2 前端，JWT 鉴权、Redis 投票排行、Swagger 文档、Docker Compose 一键部署。

[![license](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8.svg)](go.mod)
[![Vue Version](https://img.shields.io/badge/Vue-2.6-4FC08D.svg)](frontend/package.json)

</div>

## 项目简介

Bluebell 是一个前后端分离的社区论坛项目，最初作为 Go Web 开发学习项目逐步演进而来。后端基于 Gin 提供 RESTful API，使用 JWT 实现无状态鉴权；MySQL 负责用户、社区与帖子等关系数据的持久化，Redis 承担投票计数、帖子热度与最新排序；前端使用 Vue 2 全家桶，通过 Nginx 反向代理与后端衔接。项目支持 Docker Compose 一键拉起全部服务，内置 Swagger 接口文档、pprof 性能分析与 GitHub Actions 持续集成。

## 功能特性

### 用户

- 注册 / 登录，密码 bcrypt 加密存储
- JWT 无状态鉴权，登录后可访问受保护接口

### 社区

- 社区列表与详情查询

### 帖子

- 发布帖子、查看帖子详情
- 分页列表，支持按发布时间 / 热度排序

### 投票与排行

- 基于 Redis 的帖子投票：赞同（+1）/ 反对（-1）/ 取消（0）
- 发帖一周内可投票，过期自动锁定
- 热度算法：分数 = 发帖时间 + 票数 × 432（约 200 票可抵消一天的帖子年龄）

### 工程化

- Swagger 接口文档：`/swagger/index.html`
- Docker Compose 一键部署（MySQL + Redis + 后端 + 前端）
- 优雅关机、配置文件热更新、接口限流、pprof 性能分析

## 技术栈

| 层     | 技术                                        |
| ------ | ------------------------------------------- |
| 后端   | Go 1.25 · Gin · sqlx · go-redis/v9 · Viper · Zap |
| 数据   | MySQL 8.0 · Redis 7                         |
| 鉴权/ID | JWT · Snowflake 分布式 ID                   |
| 前端   | Vue 2 · Vue Router · Vuex · Axios · Nginx   |
| 工程   | Docker Compose · Swagger · GitHub Actions   |

## 快速开始

### 方式一：Docker Compose（推荐，想快速体验用这个）

前置要求：Docker Desktop（含 Docker Compose）、Git。

1. 创建本地环境变量文件：

   ```powershell
   Copy-Item .env.example .env
   ```

2. 编辑 `.env`，填入 3 个随机值（这些是密码，**不要**提交到仓库）：

   ```text
   MYSQL_ROOT_PASSWORD=一段随机字符串
   MYSQL_PASSWORD=另一段随机字符串
   JWT_SECRET=至少 32 位的随机字符串
   ```

   可以用在线密码生成器，或在 PowerShell 里生成：

   ```powershell
   -join ((48..57)+(65..90)+(97..122) | Get-Random -Count 32 | ForEach-Object { [char]$_ })
   ```

3. 构建并启动全部服务（MySQL 首次启动会自动执行 `models/create_table.sql` 建表并写入示例社区）：

   ```powershell
   docker compose up --build -d
   ```

4. 查看服务状态：

   ```powershell
   docker compose ps
   ```

启动后访问：

| 服务          | 地址                                    |
| ------------- | --------------------------------------- |
| Web 应用      | http://localhost:8080                   |
| 后端 API      | http://localhost:8888                   |
| Swagger 文档  | http://localhost:8888/swagger/index.html |

停止服务：`docker compose down`（仅在确定要清空数据时使用 `docker compose down -v`）。

### 方式二：本地开发

需要：Go 1.25+、MySQL 8、Redis 7、Node.js。

1. **准备数据库**
   - 启动本地 MySQL，创建一个名为 `bluebell` 的数据库
   - 导入表结构与示例社区数据：

     ```powershell
     mysql -uroot -p bluebell -e "source models/create_table.sql"
     ```

   - 启动本地 Redis（默认端口 6379 即可）

2. **填写配置**
   - 复制 `settings/config.example.yaml` 为 `settings/config.yaml`
   - 填入你的 MySQL 连接信息，以及一个至少 32 位的 `jwt.secret`

3. **启动后端**（默认监听 8081 端口）：

   ```powershell
   go run .
   ```

4. **启动前端**（开发服务器会把 `/api/v1` 请求代理到本地后端）：

   ```powershell
   Set-Location frontend
   npm ci
   npm run serve
   ```

   > 💡 如果前端页面请求报网络错误，先检查 `frontend/vue.config.js` 里代理的 `target` 端口是否与后端实际端口一致。

## API 接口一览

> 全部接口前缀为 `/api/v1`，🔒 表示需要请求头携带 `Authorization: Bearer <token>`。

| 方法 | 路径                                    | 说明                                                         | 鉴权 |
| ---- | --------------------------------------- | ------------------------------------------------------------ | ---- |
| POST | `/signup`                               | 用户注册                                                     |      |
| POST | `/login`                                | 用户登录                                                     |      |
| GET  | `/ping`                                 | 连通性测试                                                   | 🔒   |
| GET  | `/community`                            | 社区列表                                                     |      |
| GET  | `/community/:id`                        | 社区详情                                                     |      |
| POST | `/post`                                 | 发布帖子                                                     | 🔒   |
| GET  | `/post/:id`                             | 帖子详情                                                     |      |
| GET  | `/posts/?page=&size=`                   | 帖子列表（分页）                                             |      |
| GET  | `/posts2/?order=time\|score&page=&size=` | 帖子列表（按时间 / 热度排序）                                |      |
| POST | `/vote`                                 | 帖子投票（direction：1 赞同 / -1 反对 / 0 取消）             | 🔒   |

完整的请求 / 响应示例见 Swagger 文档。

## 项目结构

```text
bluebell/
├── controller/   # HTTP 参数解析与统一响应
├── dao/          # 数据访问层（mysql / redis）
├── logic/        # 业务逻辑层
├── middlewares/  # JWT 鉴权、接口限流
├── models/       # 数据模型与建表 SQL
├── pkg/          # 通用包（JWT、Snowflake）
├── routes/       # Gin 路由注册
├── settings/     # 配置加载与示例配置
├── logger/       # Zap 日志封装
├── docs/         # Swagger 自动生成文件
├── frontend/     # Vue 2 前端
├── Dockerfile    # 后端镜像
├── docker-compose.yml
└── main.go       # 程序入口
```

## 质量检查

```powershell
go test ./...
go vet ./...

Set-Location frontend
npm ci
npm run lint
npm run build
```

## 参与贡献

从 Issue 建立上下文 → 在独立分支完成修改 → 通过 Pull Request 合并到 `main`。

## 许可证

本项目基于 [Apache License 2.0](LICENSE) 开源。
