# Go Vue Admin

通用后台管理平台，基于 Goravel、Vue 3、TypeScript 和 shadcn-vue 构建。项目以 Resource Manifest 为统一契约，将后端 API、权限、菜单和 Vue 管理页面连接起来，用于快速构建可维护的业务后台。

> A modular admin platform built with Goravel, Vue 3, TypeScript, and shadcn-vue. Resource Manifests connect backend APIs, authorization, navigation, and Vue admin pages through one contract.

## 中文说明

### 当前能力

- JWT 登录、刷新、注销和 PostgreSQL 刷新令牌存储
- RBAC：用户、角色、权限、数据范围和字段权限
- Resource Registry 驱动的列表、表单、详情和删除流程
- 搜索、筛选、排序、分页、CSV 导出
- 当前页选择、跨页选择、按筛选条件全选
- 批量删除、批量更新和 Manifest 驱动的批量 Action
- `belongsTo` 关系选择、`hasMany` 详情摘要、表单分组和基础字段联动
- 全局搜索、菜单权限过滤和资源导航分组
- 审计日志、请求/响应脱敏、手动清理和定期清理命令入口
- `admin:make-resource` 资源生成器：生成资源骨架、权限、菜单、前端页面、迁移、测试和 README
- OpenAPI 文档与生成的 TypeScript API Client

### 资源生成

在 `backend/` 目录执行：

```powershell
go run . admin:make-resource announcements --fields="title:text:required,status:select:required:draft=Draft|published=Published"
```

生成器使用一份 ResourceSpec 生成后端资源、权限、菜单、列表/表单/详情页面、路由描述、测试、README 和迁移文件。迁移文件只生成不自动执行，执行前请人工审阅。生成资源会通过 discovery 自动进入后台面板，不需要手动维护重复菜单和路由。

### 本地开发

PostgreSQL 和 Redis 由 Laragon 管理，请先启动对应服务。

```powershell
# terminal 1
cd backend
go run .

# terminal 2
cd admin
pnpm install
pnpm dev
```

默认端口已明确隔离：

| 服务 | 地址 |
| --- | --- |
| Go API | `http://127.0.0.1:3000` |
| Go/Vue Admin | `http://127.0.0.1:5180` |
| Python Admin（其他项目） | `http://127.0.0.1:5173` |

Go/Vue 前端使用 `5180` 并启用 `strictPort`，避免与 Python 项目共用端口。后端 CORS 已允许该来源。

### 测试与构建

```powershell
cd backend
go test ./... -count=1

cd ..\admin
pnpm exec vue-tsc -b
pnpm run build
node --experimental-strip-types --test tests/*.test.ts
```

### 生产边界

当前平台核心能力已经完成并通过本地接口、浏览器、前端构建和后端测试验收，可以在此基础上开发业务资源。正式生产发布前，仍需针对实际业务完成迁移审阅、真实数据关系验收、部署配置、备份恢复、日志保留策略、负载测试和安全审计。

通用 CSV 导入、插件运行时、租户/部门扩展和复杂业务流程不属于当前通用 Resource 的默认能力；复杂业务应使用独立模块页面、Controller、Service 和 Repository。

## English

### Features

- JWT login, refresh, logout, and PostgreSQL-backed refresh-token storage
- RBAC with users, roles, permissions, data-scope rules, and field permissions
- Resource Registry driven list, form, detail, and delete flows
- Search, filters, sorting, pagination, and CSV export
- Current-page selection, cross-page selection, and select-all-by-filter
- Bulk delete, bulk update, and Manifest-driven bulk Actions
- `belongsTo` relation fields, `hasMany` detail summaries, form groups, and basic field dependencies
- Global search, permission-aware navigation, and grouped resource menus
- Audit logs with request/response redaction, manual cleanup, and scheduled cleanup entry points
- `admin:make-resource` for backend resources, permissions, menus, Vue pages, migrations, tests, and README files
- OpenAPI documentation and a generated TypeScript API client

### Resource generation

Run the generator from `backend/`:

```powershell
go run . admin:make-resource announcements --fields="title:text:required,status:select:required:draft=Draft|published=Published"
```

One ResourceSpec produces the backend resource, permissions, navigation, list/form/detail pages, route metadata, tests, README, and a migration file. Migrations are generated but never executed automatically and must be reviewed first. Generated resources enter the admin panel through discovery files without duplicate manual route and menu registration.

### Local development

Start PostgreSQL and Redis through Laragon first.

```powershell
cd backend
go run .

cd ..\admin
pnpm install
pnpm dev
```

The local ports are intentionally isolated:

- Go API: `http://127.0.0.1:3000`
- Go/Vue Admin: `http://127.0.0.1:5180`
- Python Admin (separate project): `http://127.0.0.1:5173`

### Status

The generic admin foundation is ready for business-module development. Production rollout still requires environment-specific migration review, deployment hardening, backup/recovery validation, load testing, and security review.

## Project layout

```text
backend/
  app/core/        platform contracts and shared admin infrastructure
  app/modules/     business modules and generated resources
  app/services/    domain services grouped by domain
  app/console/      console commands grouped by domain
admin/
  src/components/  shared UI components
  src/core/        frontend infrastructure and Resource Engine
  src/modules/     business pages and components
docs/              architecture, contracts, generator, testing, and roadmap
```

## License

See the repository license before distributing or deploying this project.
