# Resource Engine

Resource Engine 是标准 CRUD 加速器，不是强制所有业务使用的低代码平台。

## 两种页面模式

```text
Generic Resource   Manifest + 通用列表、表单、详情和基础操作
Custom Resource    Manifest + 显式声明的专用列表、表单或详情页面
```

## 当前实现

资源 Registry 由 `backend/app/modules/admin/registry` 统一维护，底层契约位于
`backend/app/core/resource`。注册表拒绝空标识和重复资源，并按资源名稳定排序。

- `GET /api/v1/admin/registry` 返回当前 Admin 资源清单；
- 清单接口要求当前用户至少拥有一个已注册资源的 view 权限，结果中的每个资源仍按自身权限过滤；
- 当前登记 `users`、`roles`、`permissions` 和 `departments` 资源；`departments` 是无专用 Vue 页面通用资源验收样例；
- 列表、详情、新增、编辑、删除和 CSV 导出统一使用 `/api/v1/admin/{resource}` 与
  `/api/v1/admin/{resource}/{id}`；
- 导出接口为 `/api/v1/admin/{resource}/export`，复用搜索、筛选和排序参数；
- 列表支持 `page`、`per_page`、`search`、`sort`、`dir` 参数，返回 `data` 与 `meta`；
- 排序字段按资源白名单限制，搜索使用参数绑定。
- Resource Manifest 可声明 `data_scope: all|own`；`own` 资源必须同时声明整数型
  `owner_field`，后端列表、详情、导出和写操作会统一应用数据范围约束。
- Resource Field 可声明 `visible/readable/writable/sensitive`；后端响应裁剪、查询、导出和写入保护统一使用该策略，前端隐藏仅用于改善体验。

资源页面必须复用 Registry 契约，不能在页面内重新定义资源元数据。普通资源默认使用
`admin/src/core/resource/pages/` 下的通用列表、表单和详情页面；只有 Manifest 显式声明
`pageMode: custom` 时才使用资源模块下的专用页面。用户和角色的
统一 CRUD 入口仍由各自 Domain Service 承担安全校验；关系操作继续使用专用关系端点，
批量操作统一通过 Manifest Action 执行。

用户资源写操作已通过独立的 `admin.users.manage` 权限保护：

- `POST /api/v1/admin/users` 创建用户并使用 Goravel Hash 保存密码；
- `PUT /api/v1/admin/users/:id` 编辑用户，密码留空时保持不变；
- `DELETE /api/v1/admin/users/:id` 删除用户并清理角色关联；
- `POST /api/v1/admin/users/actions/set-status` 批量设置用户状态，复用状态校验和最后管理员保护；
- 最后一个具备管理权限的活动管理员不能被删除。

资源 Manifest 通过 `actions` 声明可执行 Action；每个 Action 同时声明 `batch`，批量 Action 还必须声明服务端白名单中的 `kind` 与 `payload` 契约。列表页根据清单和当前权限渲染行级、批量操作。统一入口为 `POST /api/v1/admin/{resource}/actions/{action}`，请求使用 `ids` 与 `payload`，响应区分 `succeeded`、`failures` 和 `skips`。Manifest 只描述动作元数据，不包含可执行代码；未知 kind 不得被前端执行。

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

## Navigation contract

Registry 返回的 Resource Manifest 同时携带导航元数据：

- `navigation.group`：侧边栏分组，内置资源使用 `system`，生成资源默认使用 `business`；
- `navigation.order`：组内稳定排序；
- `navigation.hidden`：只隐藏菜单，不影响 API、路由或权限契约。

每个资源只有一个规范菜单和列表路由。`ResourceListPage` 是共享渲染组件，但只渲染路由对应的当前资源，不承担资源目录或跨资源切换职责。

标准能力：分页、搜索、排序、筛选、CSV 导出、新增、编辑、查看、删除、行操作、批量操作、字段权限和审计钩子。  

多步骤流程、审批和状态机、实时数据、图表分析、复杂联动表单和高度定制详情页才使用 Custom Resource；专用页面仍必须复用统一权限、数据范围、字段权限、Action 和审计契约。

扩展优先使用明确的 Go 接口和 Vue slot/component，而不是无限增加 JSON 字段。禁止把任意可执行代码序列化进 Resource Manifest。

```text
模块注册 -> Resource 注册 -> 权限/导航注册 -> API 暴露 -> Vue 页面消费
```
