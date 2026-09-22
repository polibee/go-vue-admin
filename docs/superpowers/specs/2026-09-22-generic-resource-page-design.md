# 通用 Resource 页面设计

## 目标

将普通资源的默认开发模式从“生成器生成每个资源的三套 Vue 页面”调整为“Resource Manifest 驱动通用页面”。新增普通资源后，后台自动获得菜单、路由、列表、搜索、筛选、表单、详情、权限、导出和批量操作；只有复杂业务才创建专用页面覆盖。

## 背景与问题

当前 admin:make-resource 会为资源生成列表页、表单页和详情页，并写入前端生成路由。这样可以快速得到可用页面，但普通资源之间产生重复页面，生成结果与通用 Resource Engine 的复用目标不一致，也使后续字段、批量操作和权限能力容易出现资源间分叉。

对照项目采用 ResourceDefinition -> Manifest -> GenericResourcePage 的方式。当前项目已经具备更丰富的 Manifest、字段权限、数据范围、软删除、筛选和批量选择能力，本次只调整默认页面策略，不降低这些已有能力。

## 设计原则

1. 普通资源默认使用通用页面。
2. Manifest 是菜单、路由、API、权限、表格、表单、详情和 Action 的唯一声明来源。
3. 复杂资源可以在自己的模块目录内覆盖列表、表单或详情页面。
4. 资源 URL 使用 /<resource>，不把 resources 暴露在 URL 中。
5. 迁移文件只生成、不自动执行。
6. 生成器不得覆盖人工创建的文件。
7. 后端权限和字段过滤始终在服务端执行，前端 Manifest 只用于展示和交互。

## 页面模式

Manifest 增加页面模式：

- generic：默认值，使用核心通用列表、表单和详情页面。
- custom：由业务模块提供一个或多个页面覆盖。

普通资源只生成：

~~~text
admin/src/modules/<resource>/resource.ts
~~~

复杂资源可以额外提供：

~~~text
admin/src/modules/<resource>/
├── pages/
│   ├── <Resource>ListPage.vue
│   ├── <Resource>FormPage.vue
│   └── <Resource>DetailPage.vue
├── components/
└── api.ts
~~~

覆盖必须显式声明，不能因为目录中存在文件就隐式改变页面模式。

## 前端边界

核心通用页面位于：

~~~text
admin/src/core/resource/pages/
├── ResourceListPage.vue
├── ResourceFormPage.vue
└── ResourceDetailPage.vue
~~~

通用页面通过 resourceName 加载 Manifest，不写死具体资源字段。组件职责保持清晰：

~~~text
admin/src/components/resource/
├── ResourceTable.vue
├── ResourceFilterBar.vue
├── ResourceForm.vue
├── ResourceDetail.vue
├── ResourceActionBar.vue
├── ResourceBulkActionBar.vue
└── ResourcePagination.vue
~~~

路由由 Registry/Manifest 统一生成：

~~~text
/<resource>
/<resource>/new
/<resource>/:id
/<resource>/:id/edit
~~~

路由必须继续经过权限守卫；不存在或无权访问的资源不能仅凭猜测 URL 打开。

## 后端边界

后端继续使用现有职责边界：

~~~text
backend/app/core/resource/
backend/app/modules/<resource>/resource/
backend/app/services/<domain>/
backend/app/core/admin/controllers/
~~~

通用资源请求通过 Manifest 查找字段、列、筛选、排序、关系和 Action，再交给资源服务或资源处理器执行。控制器不能通过不断增加资源名称分支实现业务差异；复杂业务逻辑必须放到对应模块。

资源权限、数据范围、字段权限、敏感字段脱敏和批量选择校验必须在后端执行。前端传来的字段列表、筛选条件和选择模式都必须经过 Manifest 和授权校验。

## 生成器行为

admin:make-resource <name> 作为普通资源的唯一公开入口：

1. 校验资源名称和字段定义。
2. 生成后端 Resource、服务边界、权限、菜单、测试和 README。
3. 生成迁移文件，但不执行迁移。
4. 生成前端 ResourceSpec，不生成专用页面。
5. 自动写入 Registry/路由发现所需的生成文件。
6. 运行冲突预检，禁止覆盖人工文件或产生半套写入。
7. 输出迁移审阅、执行和验证命令。

生成结果必须能自动出现在后台菜单、列表、表单、详情、权限、搜索、导出和批量 Action 中。

复杂页面不通过第二套生成命令解决。后续如需自定义页面，使用同一资源模块的显式覆盖声明和模块文件。

## 普通资源验收样例

使用 departments 作为无专用页面的标准验收资源：

1. 执行资源生成命令。
2. 审阅并手动执行迁移。
3. 启动 PostgreSQL、Redis、Go API 和 Go/Vue 管理面板。
4. 确认 departments 自动出现在正确菜单分组。
5. 验证列表、分页、搜索、筛选、排序、创建、编辑、详情和删除。
6. 验证导出、当前页选择、跨页选择、按筛选条件全选和批量 Action。
7. 验证资源权限、数据范围和字段权限。
8. 确认不存在 DepartmentsListPage.vue、DepartmentsFormPage.vue 和 DepartmentsDetailPage.vue。
9. 运行后端测试、前端类型检查、前端测试和生产构建。

## 复杂资源验收边界

选择订单或公告类复杂资源验证覆盖能力：

- 通用页面仍然是默认路径。
- 只有明确声明 custom 的资源使用专用页面。
- 自定义页面仍复用权限、数据范围、字段权限和审计边界。
- 自定义 Action 必须仍使用统一 Action 契约。

## 非目标

本阶段不做以下事情：

- 不新增第二个公开生成器入口。
- 不执行数据库迁移。
- 不重写已有业务服务。
- 不把所有现有资源一次性强制迁移为无专用页面。
- 不删除复杂资源的覆盖能力。
- 不把 PostgreSQL、Redis 替换为内存实现。
- 不改变请求/响应审计脱敏契约。

## 完成标准

- 普通资源生成后不再默认产生三套专用 Vue 页面。
- 通用页面可以仅凭 Manifest 完成资源 CRUD。
- 复杂资源可以显式覆盖页面而不污染核心目录。
- 资源菜单、路由、权限、搜索、导出和批量操作均由同一份 Manifest 驱动。
- departments 端到端验收通过。
- 全量后端测试、前端类型检查、前端测试和生产构建通过。
- 文档、生成器 Golden File 和示例输出同步更新。

## 最终决策：双模式 Resource Engine

本项目采用一个 Resource Engine、两种页面模式，不建立两套生成流程。

### 普通资源模式

普通资源默认使用：

~~~text
Resource Manifest
    -> Resource Registry
    -> Generic Resource List
    -> Generic Resource Form
    -> Generic Resource Detail
~~~

执行 admin:make-resource 后，普通资源只生成 ResourceSpec、后端资源契约、迁移、权限、菜单、测试和 README。前端路由根据 Manifest 的 resource name 进入核心通用页面，不生成资源专用的 ListPage、FormPage 或 DetailPage。

普通资源必须通过 Manifest 支持以下能力：

- 列表、分页、搜索、筛选、排序。
- 创建、编辑、详情和删除。
- 导出和统一批量 Action。
- 字段权限、数据范围和敏感字段处理。
- 资源菜单、权限守卫和错误状态。

### 复杂资源模式

复杂资源仍然使用同一份 Manifest 和同一套 API、权限、审计及 Action 契约，只允许替换页面表现层：

~~~text
Resource Manifest
    -> Resource Registry
    -> Custom Resource List/Form/Detail
~~~

复杂资源必须在 ResourceSpec 中显式声明 pageMode: custom，并在自身模块目录内提供覆盖页面。不能通过扫描文件是否存在来自动切换模式，也不能在通用页面中增加针对某个业务名称的分支。

### 页面解析优先级

页面解析顺序固定为：

1. 读取 Resource Manifest。
2. 未声明 pageMode 时使用 generic。
3. pageMode 为 generic 时强制使用核心通用页面。
4. pageMode 为 custom 时读取该资源模块声明的页面覆盖。
5. 覆盖页面缺失、导出错误或权限配置无效时，启动检查和测试必须失败，不能静默回退成半可用页面。

### 统一契约边界

无论使用 generic 还是 custom，以下能力都不能被页面覆盖绕过：

- 后端认证和资源权限。
- 数据范围和字段权限。
- Manifest 字段可读、可写、可搜索和可排序约束。
- 批量选择协议及 query selection 的服务端校验。
- 审计记录、敏感字段脱敏和错误响应契约。
- 统一 URL 规则和权限路由守卫。

### 选择标准

满足以下条件的资源使用 generic：

- 标准表格列表。
- 标准字段类型和筛选器。
- 标准 CRUD。
- 标准批量操作。
- 没有复杂状态流转或多步骤交互。

满足任意条件的资源才使用 custom：

- 多步骤或强业务约束表单。
- 订单、支付、审批、工作流等状态机。
- 复杂关系编辑或拖拽编排。
- 需要业务专属时间线、统计面板或实时交互。
- 通用表格无法表达的领域操作。

### 结果

最终开发者体验为：

~~~text
admin:make-resource departments
    -> 生成契约
    -> 自动进入菜单和路由
    -> 自动使用通用列表、表单、详情

orders 模块显式声明 pageMode: custom
    -> 仍复用同一份 Manifest、权限、API 和 Action
    -> 仅替换订单专用页面
~~~

因此，普通资源不会产生重复页面，复杂资源也不会被通用页面限制；两者共享同一个 Resource Engine 和安全边界。
