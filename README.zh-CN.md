# Yiqi Docs

[English](./README.md) | 简体中文

Yiqi Docs 是一个源码开放、强制非商用的多站点文档系统，适合承载产品说明、API 文档、教程、内部知识库等 Markdown 内容。项目包含 Go 后端、Vue 管理后台、公开文档阅读页、SQLite 存储，以及“前端嵌入 Go 服务端”的单二进制部署方案。

## 效果截图

公开文档首页：

![公开文档首页](./docs/screenshots/public-doc-overview.jpg)

API 详情页：

![API 详情页](./docs/screenshots/api-detail.jpg)

## 功能特性

- 多站点文档：通过路径区分多个文档站点，例如 `/api-docs`、`/user-guide`。
- 树形页面：支持多级页面结构，自动生成侧边栏导航和页面路径。
- Markdown 渲染：公开页支持 Markdown、代码高亮和页面目录提取。
- 管理后台：支持站点、页面、排序、用户、媒体上传管理。
- 权限控制：`super_admin` 可管理全部内容；`admin` 只管理自己拥有的站点。
- 默认 SQLite：不依赖外部数据库，适合轻量部署。
- 单二进制部署：构建 Vue 前端并嵌入 Go 服务端，一个可执行文件即可运行。
- 可选抓取脚本：可从公开或共享文档源初始化站点内容。

## 技术栈

后端：

- Go
- Gin
- GORM
- SQLite
- JWT Cookie 鉴权

前端：

- Vue 3
- Vite
- Pinia
- Vue Router
- Markdown-It
- Highlight.js
- md-editor-v3

## 项目结构

```text
.
├── backend/
│   ├── cmd/
│   │   ├── server/      # HTTP 服务入口
│   │   └── scrape/      # 可选初始化/抓取命令
│   ├── internal/
│   │   ├── api/         # HTTP handlers 与路由
│   │   ├── auth/        # 密码哈希与 JWT
│   │   ├── config/      # 环境变量配置
│   │   ├── model/       # GORM 模型
│   │   ├── scraper/     # 初始化与抓取逻辑
│   │   └── store/       # 数据访问层
│   └── web/             # 嵌入式前端资源
├── frontend/
│   ├── src/
│   │   ├── api/         # API 客户端与类型
│   │   ├── components/  # 文档布局组件
│   │   ├── stores/      # Pinia 状态
│   │   └── views/       # 公开页与后台页面
├── docs/screenshots/    # README 截图
└── build.sh             # 生产构建脚本
```

## 环境要求

- Go，版本以 `backend/go.mod` 声明为准
- Node.js 20+ 或 22+
- pnpm
- 通过 Go 依赖提供 SQLite 支持

如需安装 pnpm：

```bash
corepack enable
corepack prepare pnpm@latest --activate
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/mujiangliu/yiqi-docs.git
cd yiqi-docs
```

### 2. 启动后端

```bash
cd backend
JWT_SECRET=change-me \
SEED_ADMIN_USER=admin \
SEED_ADMIN_PASS=change-me \
go run ./cmd/server
```

后端默认监听 `http://localhost:8080`。

### 3. 启动前端开发服务

另开一个终端：

```bash
cd frontend
pnpm install
pnpm dev
```

前端默认监听 `http://localhost:5173`。开发环境中，Vite 会把 `/api` 代理到 `http://localhost:8080`。

### 4. 打开页面

- 公开文档：`http://localhost:5173/<site-path>/`
- 管理后台：`http://localhost:5173/admin`
- 如果按上面的命令创建了初始账号：
  - 用户名：`admin`
  - 密码：你设置的 `SEED_ADMIN_PASS`

## 运行配置

| 变量 | 默认值 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| `PORT` | `8080` | 否 | HTTP 服务端口。 |
| `DB_PATH` | `./data.db` | 否 | SQLite 数据库文件路径。 |
| `JWT_SECRET` | 空 | 是 | JWT Cookie 签名密钥。 |
| `SEED_ADMIN_USER` | `admin` | 否 | 初始超级管理员用户名。 |
| `SEED_ADMIN_PASS` | 空 | 首次运行需要 | 初始超级管理员密码；为空时跳过创建。 |
| `SEED_CONTENT_ADMIN_USER` | `content` | 仅抓取脚本需要 | 抓取脚本使用的内容管理员用户名。 |
| `SEED_CONTENT_ADMIN_PASS` | 空 | 仅抓取脚本需要 | 抓取脚本使用的内容管理员密码。 |
| `APIFOX_SHARED_DOC_TOKEN` | 空 | 否 | 可选 Apifox shared-doc token；为空时抓取脚本会尝试公开页面抓取。 |

不要把生产环境变量提交到仓库。

## 初始化与抓取

可选抓取命令可以创建初始用户、创建文档站点并写入内容。

```bash
cd backend
DB_PATH=../local.db \
JWT_SECRET=change-me \
SEED_ADMIN_USER=admin \
SEED_ADMIN_PASS=change-me \
SEED_CONTENT_ADMIN_USER=content \
SEED_CONTENT_ADMIN_PASS=change-me \
go run ./cmd/scrape
```

如果设置了 `APIFOX_SHARED_DOC_TOKEN`，抓取脚本会优先尝试 shared-doc JSON API；如果没有设置或请求失败，则在可行时回退到公开页面抓取。

## 生产构建

在项目根目录运行：

```bash
./build.sh
```

脚本会依次执行：

1. 使用 pnpm 安装前端依赖。
2. 构建 Vue 生产资源。
3. 将 `frontend/dist` 复制到 `backend/web/dist`。
4. 构建 Go 服务端二进制 `./jiaocheng-web`。
5. 构建可选抓取二进制 `./scrape`。

## 运行生产二进制

```bash
DB_PATH=./data.db \
JWT_SECRET=change-me \
PORT=8080 \
./jiaocheng-web
```

访问：

```text
http://localhost:8080/<site-path>/
```

由于前端已经嵌入 Go 服务端，同一个服务会同时处理 API 路由和 SPA 路由。

## 部署示例

### PM2

```bash
pm2 start ./jiaocheng-web \
  --name yiqi-docs \
  --interpreter none
```

生产环境建议使用 ecosystem 文件：

```js
module.exports = {
  apps: [
    {
      name: "yiqi-docs",
      script: "./jiaocheng-web",
      interpreter: "none",
      cwd: "/opt/yiqi-docs",
      env: {
        PORT: "8090",
        DB_PATH: "/opt/yiqi-docs/data/docs.db",
        JWT_SECRET: "replace-with-a-long-random-secret"
      }
    }
  ]
}
```

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name docs.example.com;

    location / {
        proxy_pass http://127.0.0.1:8090;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

HTTPS 可以使用 Certbot、acme.sh 或云平台证书服务配置。

## API 概览

公开接口：

- `GET /api/sites/:path`：获取已发布文档站点及全部页面。
- `GET /api/media/:hash`：获取上传的媒体资源。

鉴权接口：

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`

管理接口：

- `GET /api/admin/sites`
- `POST /api/admin/sites`
- `PUT /api/admin/sites/:id`
- `DELETE /api/admin/sites/:id`
- `GET /api/admin/sites/:id/pages`
- `POST /api/admin/sites/:id/pages`
- `PUT /api/admin/pages/:id`
- `DELETE /api/admin/pages/:id`
- `POST /api/admin/pages/reorder`
- `GET /api/admin/users`
- `POST /api/admin/users`
- `PUT /api/admin/users/:id`
- `POST /api/admin/users/:id/reset-password`
- `DELETE /api/admin/users/:id`

## 安全说明

- 不要提交运行时数据库、WAL/SHM 文件、日志、证书、私钥或 `.env` 文件。
- 生产环境必须设置强随机 `JWT_SECRET`。
- 生产环境建议把 Go 服务放在 HTTPS 反向代理后面。
- 如果是内部系统，建议在网络层或身份层限制后台访问。
- 部署或迁移前请备份 SQLite 数据库。

## 常用开发命令

```bash
# 前端类型检查
pnpm --dir frontend typecheck

# 后端测试
cd backend && go test ./...

# 生产构建
./build.sh
```

## License

本项目使用 [Yiqi Docs Non-Commercial License v1.0](./LICENSE)。

你可以将本项目用于个人学习、研究、测试、非商业内部实验、非商业修改和非商业分发。任何商业用途都必须提前联系版权持有人或仓库所有者购买商业授权。

商业用途包括但不限于：商业产品、商业 SaaS、客户项目、商业交付、转售、咨询交付、付费服务或其他营利活动。
