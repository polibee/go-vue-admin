# Goravel 集成清单

## 基线

~~~text
Goravel v1.18.x
goravel/gin
goravel/postgres
goravel/redis
Vue 3
shadcn-vue
~~~

Goravel 项目在 backend/ 创建，前端项目在 admin/ 创建。安装器使用 latest 创建项目后，实际依赖版本必须固定在 backend/go.mod 和 backend/go.sum。

## 必须集成

### goravel/gin

作为 HTTP Driver。API、路由、中间件和 CORS 使用 Goravel 能力。

### goravel/postgres

作为默认数据库驱动。ORM 继续使用 Goravel ORM，业务代码不直接依赖 PostgreSQL 驱动 API。

### goravel/redis

用于 Cache、Queue、Refresh Token、限流和临时状态。

## 可选兼容

### goravel/mysql

当项目需要 MySQL 兼容时安装。数据库连接通过 DB_CONNECTION 选择：

~~~text
postgres
mysql
~~~

ORM 代码可以保持一致，但 Migration、JSON、索引、锁、排序和数据库函数仍需分别测试。不能因为使用 ORM 就宣称数据库完全无差异。

兼容矩阵至少包含：

~~~text
PostgreSQL Migration + Seed + CRUD
MySQL Migration + Seed + CRUD
~~~

第一阶段默认只以 PostgreSQL 作为发布数据库；MySQL 作为兼容性目标。

## 按需集成

### goravel/oss

只在 Media/File Module 中接入。Core 只依赖 Goravel Storage 接口，不直接依赖 OSS。第一阶段使用本地存储。

### goravel/openai

只在 AI Module 中接入，不进入 Admin Core。AI 不得成为框架运行依赖。

### goravel-workflow

只在 Workflow/Approval Module 中接入。它会发布流程表、配置和 Hook，属于业务域，不属于 Admin Foundation。

### goravel-socket

只在 Realtime/Notification Module 中接入。普通后台 CRUD 不需要 WebSocket。

## 不集成

### goravel-crud

不作为依赖。它自带 CRUD 路由、JWT Middleware、Swagger、面板和资源发布约定，与本项目的 Resource Engine、认证和 OpenAPI Contract 重复。

可以参考其 CRUD 生成思路，但不复制其运行时架构。

## 安装原则

不要把所有驱动都无条件安装。数据库驱动从 Goravel v1.16 起是独立包，只安装项目实际使用或明确兼容测试需要的驱动。

安装后必须：

- 注册对应 Service Provider；
- 配置 database.go；
- 配置环境变量；
- 执行 Migration；
- 执行 CRUD 和集成测试；
- 记录 go.mod/go.sum 版本。

## 目录

~~~text
backend/app/modules/
├── media/       可选 OSS
├── ai/          可选 OpenAI
├── workflow/    可选审批流
├── realtime/    可选 WebSocket
└── ...
~~~
