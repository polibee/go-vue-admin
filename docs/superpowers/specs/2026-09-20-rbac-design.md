# RBAC 与错误码本地化设计

## 目标

在现有 Goravel v1.18 + Vue 3 + shadcn-vue 管理后台上，完成第一版可用 RBAC：

- 统一前后端错误码，避免中英文混用；
- 管理角色；
- 为角色分配权限；
- 为用户绑定角色；
- 在后端对管理接口执行细粒度权限校验；
- 前端继续使用官方 shadcn-vue 组件和按语言目录组织的语言包。

本阶段不引入第二套权限库，不实现菜单拖拽、数据范围、审计中心和插件权限隔离。

## 已有基础

当前数据库已经存在：

- `users`；
- `roles`；
- `permissions`；
- `role_user`；
- `permission_role`。

现有 `super-admin` 角色拥有：

- `admin.users.view`；
- `admin.roles.manage`；
- `admin.permissions.manage`。

现有认证使用 JWT Access Token，RBAC 接口必须先验证 JWT，再验证权限。

## 错误码本地化

后端响应保留稳定的 `code`，`message` 只用于日志或开发辅助，不作为前端最终展示文本。

标准错误码：

| code | HTTP | 语义 |
| --- | ---: | --- |
| `VALIDATION_ERROR` | 422 | 请求参数无效 |
| `AUTH_INVALID_CREDENTIALS` | 401 | 账号或密码错误 |
| `AUTH_UNAUTHORIZED` | 401 | 未认证或 Token 无效 |
| `RBAC_FORBIDDEN` | 403 | 没有权限 |
| `RBAC_ROLE_NOT_FOUND` | 404 | 角色不存在 |
| `RBAC_PERMISSION_NOT_FOUND` | 404 | 权限不存在 |
| `RBAC_USER_NOT_FOUND` | 404 | 用户不存在 |
| `RBAC_SYSTEM_ROLE` | 409 | 系统角色不可删除或修改关键字段 |
| `RBAC_LAST_ADMIN` | 409 | 不能移除最后一个有效管理员 |
| `INTERNAL_ERROR` | 500 | 未知服务端错误 |

前端 `apiFetch` 将错误转换成 `ApiError`，由页面根据 `error.code` 调用：

```text
t(`errors.${error.code}`)
```

没有映射的错误码使用 `errors.unknown`。后端原始 `message` 不直接渲染，避免语言混用和泄漏内部信息。

## 授权模型

采用简单、可预测的 RBAC：

```text
User -> Role -> Permission
```

- 用户可以拥有多个角色；
- 用户有效权限是所有角色权限的并集；
- 一个接口声明一个或多个所需权限时，默认使用“满足任一权限”；
- `super-admin` 作为系统角色，仅通过角色名称精确识别；
- 角色、权限查询和授权判断集中在 RBAC Service，不在 Controller 中散落 SQL；
- 前端路由和按钮隐藏只用于体验，不能替代后端检查。

权限中间件提供明确的入口：

```go
middleware.RequirePermission("admin.roles.manage")
```

权限不足统一返回 `403 + RBAC_FORBIDDEN`。

## API 设计

### 角色

```text
GET    /api/v1/admin/roles
POST   /api/v1/admin/roles
GET    /api/v1/admin/roles/:id
PUT    /api/v1/admin/roles/:id
DELETE /api/v1/admin/roles/:id
PUT    /api/v1/admin/roles/:id/permissions
```

角色字段：`id`、`name`、`display_name`、`permissions`、时间字段。

`name` 是稳定标识，不允许修改系统角色的 `name`。权限更新采用完整集合提交，后端在一个事务中同步 `permission_role`。

### 权限

```text
GET /api/v1/admin/permissions
```

第一阶段权限由系统 Seeder 管理，不开放任意用户创建权限，避免权限标识失控。

### 用户角色

```text
GET /api/v1/admin/users
GET /api/v1/admin/users/:id/roles
PUT /api/v1/admin/users/:id/roles
```

角色绑定采用完整集合提交，后端校验用户和角色均存在，再在事务中同步 `role_user`。

安全约束：

- 禁止删除或停用最后一个有效管理员；
- 禁止移除当前操作员自己的最后一个管理角色；
- 不允许通过普通角色管理接口修改系统核心角色的关键属性。

## 服务端分层

```text
Route
  -> JWT Authentication
  -> Permission Middleware
  -> RBAC Controller
  -> RBAC Service
  -> Goravel ORM
```

Controller 只负责输入输出和 HTTP 状态码；Service 负责关系同步、系统角色约束、最后管理员约束和事务；Goravel ORM 负责数据库访问。

## 前端设计

`RBACView.vue` 从只读卡片升级为统一管理页面：

- 角色列表和角色编辑 Dialog；
- 权限 Checkbox 列表；
- 用户角色分配 Dialog；
- 加载使用 `Skeleton`；
- 空状态使用官方 `Empty`；
- 操作反馈使用现有 toast/Alert 机制；
- 所有标签、校验、错误码位于 `locales/zh-CN/rbac.json`、`locales/en-US/rbac.json` 和对应错误语言包。

不修改 `admin/src/components/ui` 下的官方组件源码，只组合官方组件。

## 测试策略

后端：

- 错误码和 HTTP 状态码测试；
- 未认证返回 `AUTH_UNAUTHORIZED`；
- 无权限返回 `RBAC_FORBIDDEN`；
- 角色权限完整集合同步；
- 用户角色完整集合同步；
- 系统角色保护；
- 最后管理员保护。

前端：

- 错误码只显示当前语言文本；
- 角色和权限加载、空状态、错误状态；
- 中英文切换；
- 角色权限选择和用户角色选择提交。

验收标准：登录后可以创建角色、为角色分配权限、为用户绑定角色；没有对应权限的 JWT 调用写接口必须被后端拒绝；中文和英文页面不能出现混合提示。

## 分阶段实现

1. 统一错误码和前端错误映射；
2. RBAC Service 与权限中间件；
3. 角色 CRUD 和角色权限同步；
4. 用户角色绑定与安全约束；
5. 前端管理 Dialog 和多语言；
6. 自动化测试、迁移/Seeder 验证和浏览器联调。
