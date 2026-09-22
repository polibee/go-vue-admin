# Resource Engine

Resource Engine 是标准 CRUD 加速器，不是强制所有业务使用的低代码平台。

## 三种使用等级

```text
Simple Resource    标准列表、表单、详情和基础操作
Extended Resource  标准资源 + 自定义 Action、Slot 或详情区
Custom Page        普通 Vue 页面 + Goravel Controller/Service
```

## 当前实现

资源 Registry 由 `backend/app/modules/admin/registry` 统一维护，底层契约位于
`backend/app/core/resource`。注册表拒绝空标识和重复资源，并按资源名稳定排序。

- `GET /api/v1/admin/registry` 返回当前 Admin 资源清单；
- 清单接口要求当前用户至少拥有一个已注册资源的 view 权限，结果中的每个资源仍按自身权限过滤；
- 当前登记 `users`、`roles`、`permissions` 三个基础资源；
- 列表、详情、新增、编辑、删除和 CSV 导出统一使用 `/api/v1/admin/{resource}` 与
  `/api/v1/admin/{resource}/{id}`；
- 导出接口为 `/api/v1/admin/{resource}/export`，复用搜索、筛选和排序参数；
- 列表支持 `page`、`per_page`、`search`、`sort`、`dir` 参数，返回 `data` 与 `meta`；
- 排序字段按资源白名单限制，搜索使用参数绑定。
- Resource Manifest 可声明 `data_scope: all|own`；`own` 资源必须同时声明整数型
  `owner_field`，后端列表、详情、导出和写操作会统一应用数据范围约束。
- Resource Field 可声明 `visible/readable/writable/sensitive`；后端响应裁剪、查询、导出和写入保护统一使用该策略，前端隐藏仅用于改善体验。

资源页面必须复用 Registry 契约，不能在页面内重新定义资源元数据。用户和角色的
统一 CRUD 入口仍由各自 Domain Service 承担安全校验；关系操作和用户批量状态属于
明确的业务 Action，继续使用专用关系端点。

用户资源写操作已通过独立的 `admin.users.manage` 权限保护：

- `POST /api/v1/admin/users` 创建用户并使用 Goravel Hash 保存密码；
- `PUT /api/v1/admin/users/:id` 编辑用户，密码留空时保持不变；
- `DELETE /api/v1/admin/users/:id` 删除用户并清理角色关联；
- `PUT /api/v1/admin/users/status` 批量设置用户状态，复用状态校验和最后管理员保护；
- 最后一个具备管理权限的活动管理员不能被删除。

用户资源 Manifest 通过 `actions` 声明 `set-status` Custom Action；列表页根据清单渲染单行状态操作，并复用批量状态接口。Manifest 只描述动作名称、类型和权限，不包含可执行代码。

Admin 前端提供 `/users/new`、`/users/:id/edit`、用户详情删除确认和批量状态设置流程。系统概览作为独立 Custom Page，通过 Overview Controller/Service 提供统计卡片和资源快捷入口；roles、permissions 的通用写表单仍沿用 RBAC 专用页面。

## Resource Definition

```text
name
label
route
permissions
fields
columns
filters
actions
navigation
data_scope
owner_field
```

标准能力：分页、搜索、排序、筛选、CSV 导出、新增、编辑、查看、删除、行操作、批量操作、字段权限和审计钩子。

多步骤流程、审批和状态机、实时数据、图表分析、复杂联动表单和高度定制详情页直接使用 Custom Page。

扩展优先使用明确的 Go 接口和 Vue slot/component，而不是无限增加 JSON 字段。禁止把任意可执行代码序列化进 Resource Manifest。

```text
模块注册 -> Resource 注册 -> 权限/导航注册 -> API 暴露 -> Vue 页面消费
```
