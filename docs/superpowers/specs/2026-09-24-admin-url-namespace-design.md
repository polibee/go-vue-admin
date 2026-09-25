# Admin URL Namespace Design

## 目标

明确 Go/Vue Admin 的后台页面、C 端页面和后端 API 的 URL 边界，避免前后端使用同名根路径，保证资源生成器、菜单、通知、OpenAPI 和测试使用同一套路由契约。

本设计只定义 URL 约定和迁移边界，不在本阶段直接修改运行时代码。

## 路由空间

| 类型 | 前缀 | 示例 | 所属实现 |
| --- | --- | --- | --- |
| 后台管理页面 | `/admin/` | `/admin/users` | Vue Router |
| C 端页面 | `/` | `/`, `/orders` | Vue Router 或业务前端 |
| 后台管理 API | `/api/v1/admin/` | `/api/v1/admin/users` | Go/Goravel |
| 认证 API | `/api/v1/auth/` | `/api/v1/auth/login` | Go/Goravel |
| OpenAPI 契约 | `/api/` | `/api/openapi.json` | Go/Goravel |

### 后台页面

所有通用后台页面必须位于 `/admin/` 下：

```text
/admin
/admin/login
/admin/users
/admin/users/new
/admin/users/:id
/admin/users/:id/edit
/admin/roles
/admin/permissions
/admin/audit-logs
/admin/settings
```

`/admin` 是后台首页，`/admin/login` 是后台登录页。后台页面不得继续占用 `/users`、`/roles`、`/permissions` 等根路径。

### C 端页面

根目录保留给面向最终用户的页面：

```text
/
/orders
/products
/account
```

C 端页面不得依赖 `/admin` 前缀，也不得复用后台页面的路由名称。未来 C 端与后台可以使用同名业务概念，但必须通过 URL 命名空间和 Router name 区分。

### 后端 API

后台 API 使用 `/api/v1/admin/`，不能使用与前端页面相同的根路径：

```text
GET    /api/v1/admin/users
GET    /api/v1/admin/users/{id}
POST   /api/v1/admin/users
PUT    /api/v1/admin/users/{id}
DELETE /api/v1/admin/users/{id}
GET    /api/v1/admin/registry
GET    /api/v1/admin/search
```

认证 API 保持 `/api/v1/auth/`；Admin OpenAPI 规范入口为 `/api/openapi/admin.json`，Scalar 为 `/api/docs`，不把 API 文档伪装成后台页面。

## Resource 契约

Resource Manifest 必须区分页面地址和 API 地址。页面地址属于前端 URL，API 地址属于后端 HTTP URL，二者不能使用一个含义模糊的 `route` 字段互相替代。

推荐契约：

```ts
type ResourceManifest = {
  name: string
  adminRoute: `/admin/${string}`
  apiBase: `/api/v1/admin/${string}`
}
```

资源生成器的默认结果：

```text
name:       departments
adminRoute: /admin/departments
apiBase:    /api/v1/admin/departments
```

资源页面内部跳转只能使用 `adminRoute` 或后台路由构造器；API Client 只能使用 `apiBase` 或 API Client 构造器。禁止通过字符串拼接把 `/admin/...` 当成 API 地址，也禁止把 `/api/...` 放进 Vue Router。

## 前端实现约束

1. Vue Router 使用 `/admin` 作为后台路由父级。
2. 后台所有登录、菜单、面包屑、通知跳转、列表、表单、详情和错误页都属于 `/admin` 子树。
3. 登录未认证时，后台页面统一重定向到 `/admin/login?redirect=...`。
4. 通知内部链接必须经过后台 URL 校验，只允许 `/admin/` 路径。
5. Vite 只代理 `/api` 到 Go 后端；`/admin` 和 C 端根路径由前端开发服务器处理。
6. Router name 使用 `admin-` 或资源域前缀，避免与 C 端 Router name 重复。

## 后端实现约束

1. 管理接口统一挂载在 `/api/v1/admin` 路由组下。
2. 后端不得注册 `/users`、`/roles`、`/permissions` 等页面同名接口。
3. Resource 动态路由和显式资源路由必须使用同一 API 前缀，避免重复注册和路径漂移。
4. OpenAPI 文档、生成 API Client 和契约测试必须反映 `/api/v1/admin/...`。
5. 后端返回的通知、菜单和资源元数据不得生成根路径后台链接。

## 迁移策略

本项目采用直接统一迁移，不保留旧根路径后台别名：

1. 先更新路由契约和 Resource Manifest 类型。
2. 将登录、首页、RBAC、资源、审计、设置和错误页面迁移到 `/admin` 子树。
3. 更新菜单、通知、面包屑、表单返回和详情跳转。
4. 更新生成器默认路由和生成文件。
5. 删除后端遗留的根路径页面接口。
6. 更新 OpenAPI、API Client、浏览器测试和文档。
7. 只验证新路径，不为旧后台根路径增加兼容跳转。

迁移过程中如果发现真实 C 端已经占用某个根路径，应优先调整后台路径，而不是让后台继续复用该根路径。

## 验收标准

### 后台页面

- `/admin/login` 可以登录并跳转到 `/admin`。
- `/admin/users`、`/admin/roles`、`/admin/permissions`、`/admin/audit-logs` 可正常加载。
- 资源生成后页面默认出现在 `/admin/<resource>`。
- 刷新后台深层链接不会回退到 C 端页面。
- 未登录访问后台页面只能跳转到 `/admin/login`。

### API

- 后台列表、详情、新增、更新、删除和批量 Action 请求全部使用 `/api/v1/admin/...`。
- 根路径 `/users`、`/roles`、`/permissions` 不再被后端管理接口占用。
- OpenAPI 路径、生成 API Client 和实际请求路径一致。
- `/admin/...` 页面请求不会被 Vite 代理到后端 API。

### 边界

- 根路径 C 端页面可以独立注册，不与后台 Router name 冲突。
- 通知链接只允许跳转后台 `/admin/...` 或经过明确允许的 C 端地址。
- Resource Manifest 同时提供 `adminRoute` 和 `apiBase`，不存在含义不明确的单一路由字段。

## 非目标

- 本阶段不实现 C 端业务页面。
- 本阶段不改变 API 版本号。
- 本阶段不新增依赖。
- 本阶段不保留旧后台根路径兼容别名。
- 本阶段不重构与 URL 命名空间无关的业务逻辑。
