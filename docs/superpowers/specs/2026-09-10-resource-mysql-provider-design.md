# Resource MySQL Provider Design

## Goal

将示例资源从当前的前端内存数据源扩展为可由后端选择的 MySQL 8 持久化 Provider，同时保留内存 Provider，确保开发、测试和数据库不可用时仍有明确的可运行模式。

## Scope

- 新增 `demo_resources` 数据库迁移。
- 在现有 `DemoResourceRepository` 契约下实现 MySQL Repository 和后端内存 Repository。
- 新增受认证保护的资源 CRUD API、分页、搜索、筛选、排序和批量删除。
- 通过 `RESOURCE_PROVIDER=memory|mysql` 选择后端 Provider；未配置时使用 `memory`，不静默切换。
- 前端保留现有 `MemoryResourceDataProvider`，并在配置选择 HTTP 模式时使用已经存在的 `HttpResourceDataProvider`。
- 统一使用现有响应封装、查询解析、RBAC 权限和 `ResourceDataProvider` 接口。

## Non-goals

- 本阶段不重构通用 Resource Engine，不引入第二套 API 类型系统。
- 本阶段不修改 RBAC 存储；当前 RBAC 继续使用内存仓储。
- 本阶段不自动修改 MySQL 用户、密码或数据库权限。
- 本阶段不在数据库不可用时偷偷回退到内存；Provider 选择错误或 MySQL 连接失败必须返回明确错误。

## Architecture

```text
HTTP Controller
      |
DemoResourceService  -- validation, normalization, domain errors
      |
DemoResourceRepository
      |-----------------------------|
MemoryDemoResourceRepository   MySQLDemoResourceRepository
                                      |
                                  Goravel DB facade
```

控制器只负责认证、请求边界、查询参数和错误映射；服务层不依赖 HTTP 或 SQL；Repository 负责存储细节。MySQL Repository 只使用白名单排序字段和显式字段映射，避免把用户输入拼入 SQL 片段。

## Runtime selection

后端读取 `RESOURCE_PROVIDER`：

| 值 | 行为 |
| --- | --- |
| `memory` 或空值 | 使用带有初始演示记录的内存 Repository |
| `mysql` | 使用 `demo_resources` MySQL Repository |
| 其他值 | 启动资源模块时返回配置错误，不自动猜测 |

前端通过 `VITE_RESOURCE_PROVIDER=memory|http` 选择数据源，默认 `memory`。设置为 `http` 时调用后端资源 API；HTTP Provider 和内存 Provider 共享同一个 `ResourceDataProvider` 接口。

## HTTP contract

- `GET /api/resources/demo`: 返回 `{ data, meta.pagination }`。
- `GET /api/resources/demo/:id`: 返回单条记录。
- `POST /api/resources/demo`: 接收 `id`, `name`, `status`, `owner`。
- `PUT /api/resources/demo/:id`: 接收至少一个可更新字段 `name`, `status`, `owner`。
- `DELETE /api/resources/demo/:id`: 删除单条记录。
- `POST /api/resources/demo/bulk-delete`: 接收 `{ ids: string[] }`。

所有端点要求当前会话拥有 `dashboard.view`；未认证返回 `401`，无权限返回 `403`，无效输入返回 `400`，记录不存在返回 `404`，存储错误返回 `503`。查询支持现有 `page`、`per_page`、`search`、`filter[...]`、`sort` 和 `sort_dir` 协议，`per_page` 最大值继续由共享解析器限制为 100。

## Persistence

迁移创建 `demo_resources`：字符串主键 `id`、`name`、`status`、`owner`、`created_at`、`updated_at`。`status` 由服务层限制为 `draft` 或 `active`，数据库使用可扩展字符串列而非 MySQL ENUM，避免后续状态扩展必须重建表结构。`id`、`status`、`owner` 建立查询索引，`name` 不作为唯一字段。

## Testing and verification

- 服务层单元测试覆盖输入规范化、状态校验、空更新、批量 ID 去重和 Repository 错误传播。
- MySQL Repository 使用可替换的 DB Query 边界测试字段映射、排序白名单和分页参数。
- HTTP 契约测试覆盖认证、响应封装和错误状态码。
- 前端 Provider 测试覆盖查询序列化、CRUD URL 和分页字段映射。
- 提交前运行前端全量测试、类型检查、生产构建，以及后端 `go test ./...` 和 `go vet ./...`。
- 如果当前 MySQL 凭据无法连接，只报告运行环境阻塞；不修改用户数据库权限。
