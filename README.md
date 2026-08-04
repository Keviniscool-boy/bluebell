# Bluebell

Bluebell 是一个前后端分离的社区论坛项目，后端使用 Go、Gin、MySQL 和 Redis，前端使用 Vue 2 与 Nginx。项目包含用户注册登录、JWT 鉴权、社区、帖子和投票功能，并提供 Swagger API 文档与 Docker Compose 一键启动环境。

## 功能

- 用户注册、登录和 JWT 鉴权
- 社区列表与详情
- 帖子发布、详情和分页排序
- Redis 驱动的帖子投票与热度排行
- Swagger API 文档
- MySQL、Redis、后端和前端容器化部署
- GitHub Actions 持续集成

## 技术栈

- 后端：Go 1.25、Gin、sqlx、go-redis、Viper、Zap
- 数据：MySQL 8.0、Redis 7
- 前端：Vue 2、Vue Router、Vuex、Axios、Nginx
- 工程：Docker Compose、Swagger、GitHub Actions

## 快速开始

### 环境要求

- Docker Desktop（包含 Docker Compose）
- Git

### 使用 Docker Compose

1. 创建本地环境变量文件：

   ```powershell
   Copy-Item .env.example .env
   ```

2. 修改 `.env`，为 MySQL 和 JWT 设置随机且安全的值。

3. 构建并启动全部服务：

   ```powershell
   docker compose up --build -d
   ```

4. 查看服务状态：

   ```powershell
   docker compose ps
   ```

启动后可访问：

- Web 应用：http://localhost:8080
- 后端 API：http://localhost:8888
- Swagger：http://localhost:8888/swagger/index.html

停止服务：

```powershell
docker compose down
```

MySQL 与 Redis 数据保存在 Docker 命名卷中。只有明确需要清空数据时才执行 `docker compose down -v`。

## 本地开发

后端需要先准备 MySQL、Redis，并将 `settings/config.example.yaml` 复制为 `settings/config.yaml` 后填写本地配置：

```powershell
go run .
```

前端开发服务器：

```powershell
Set-Location frontend
npm ci
npm run serve
```

开发服务器会将 `/api` 请求代理到 `http://127.0.0.1:8888`。

## 质量检查

```powershell
go test ./...
go vet ./...

Set-Location frontend
npm ci
npm run lint
npm run build
```

## 项目结构

```text
controller/   HTTP 参数处理与响应
dao/          MySQL 和 Redis 数据访问
logic/        业务逻辑
middlewares/  JWT 鉴权与限流
models/       数据模型与初始化 SQL
pkg/          JWT、Snowflake 等通用包
routes/       Gin 路由
settings/     配置加载与示例配置
frontend/     Vue 前端
docs/         Swagger 生成文件
```

## 参与贡献

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。建议从 Issue 建立上下文，在独立分支完成修改，并通过 Pull Request 合并到 `main`。

## 安全

不要提交 `.env`、`settings/config.yaml`、日志、数据库数据或真实凭据。发现安全问题时，请不要公开创建包含利用细节的 Issue，参见 [SECURITY.md](SECURITY.md)。

## 许可证

本项目基于 [Apache License 2.0](LICENSE) 开源。
