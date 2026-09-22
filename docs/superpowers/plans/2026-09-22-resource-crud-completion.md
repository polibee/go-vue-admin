# Resource CRUD 完整能力实施计划

## 执行约束

- 当前工作区直接开发，不推送 GitHub。
- 每个阶段先写失败测试，再实现，再运行对应测试。
- 不执行真实生产数据清理；审计清理验收使用隔离测试数据或现有测试数据库。
- 保持 Service 在 `backend/app/services/<domain>/`，Resource/Controller 在各自边界，前端共享组件在 `admin/src/components/`，资源运行时在 `admin/src/core/resource/`。

## 阶段 1：审计清理闭环

### 文件

- `backend/app/core/admin/controllers/audit_controller.go`
- `backend/app/services/audit/cleanup.go`
- `backend/app/services/audit/*_test.go`
- `admin/src/modules/audit/pages/AuditLogPage.vue`
- `admin/src/generated/api.ts`

### 工作

1. 为默认值、边界值、无权限和数据库删除结果增加测试。
2. 修复前后端错误码和响应解包不一致问题。
3. 增加清理按钮可见性、确认、结果提示和失败提示的前端测试。
4. 使用隔离旧日志记录验证接口删除数量、保留新记录并产生 `audit.cleanup`。

## 阶段 2：选择模型

### 文件

- `backend/app/core/admin/actions/selection.go`
- `backend/app/core/admin/actions/selection_test.go`
- `backend/app/core/admin/controllers/resource_action_controller.go`
- `backend/app/core/admin/controllers/resource_list_query.go`
- `admin/src/core/resource/lib/selection.ts`
- `admin/src/core/resource/pages/ResourceListPage.vue`

### 工作

定义 `Selection{Mode, IDs, Query, ExcludeIDs}`。`ids` 只执行明确 ID；`query` 重新执行筛选、排序无关的资源查询并应用数据范围。前端显示当前页选择和筛选结果全选状态，查询或翻页时不静默丢失模式。

## 阶段 3：批量操作

### 文件

- `backend/app/core/admin/actions/registry.go`
- `backend/app/core/admin/controllers/resource_action_controller.go`
- `backend/app/core/admin/controllers/resource_crud_controller.go`
- `backend/app/modules/admin/registry/registry.go`
- `admin/src/core/resource/pages/ResourceListPage.vue`
- `admin/src/generated/api.ts`

### 工作

增加标准 `bulk_delete`、`bulk_update` 和 Manifest 自定义 Action 处理器。删除和更新统一通过选择模型执行，返回逐项结果；前端提供通用确认和 Action 参数表单，移除用户状态专用分支。

## 阶段 4：Manifest 筛选器

### 文件

- `backend/app/core/resource/manifest.go`
- `backend/app/core/admin/controllers/resource_list_controller.go`
- `backend/app/core/admin/controllers/resource_list_capabilities.go`
- `admin/src/generated/api.ts`
- `admin/src/core/resource/pages/ResourceListPage.vue`

### 工作

增加 `filters` 声明及统一解析器，支持 select、multi-select、boolean、text、date range、relation option。前端根据声明渲染筛选控件，后端只接受声明过的字段并执行类型校验。

## 阶段 5：软删除、恢复和视图

### 文件

- `backend/app/core/resource/manifest.go`
- `backend/app/core/admin/controllers/resource_crud_controller.go`
- `backend/app/core/admin/controllers/resource_list_controller.go`
- `admin/src/core/resource/pages/ResourceListPage.vue`
- `admin/src/core/resource/lib/views.ts`
- `admin/src/stores/resource-views.ts`

### 工作

先支持 Manifest 的软删除声明、回收站过滤、恢复和永久删除权限；再以本地持久化为第一版实现列、排序、筛选视图保存，后续再接服务端用户视图。

## 验证

- `go test ./... -count=1`
- `npm run build`
- `node --experimental-strip-types --test tests/*.test.ts`
- `git diff --check`
- 前端启动后人工验收：审计清理、当前页选择、筛选全选、批量操作、回收站和视图保存。

## 里程碑

每个阶段形成一个本地提交；五阶段全部通过后再决定是否推送 GitHub。
