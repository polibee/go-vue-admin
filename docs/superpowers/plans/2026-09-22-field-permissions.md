# Resource 字段权限实施计划

> **执行约束：** 按任务逐项执行；每项先写失败测试，再实现，再运行该项验证。迁移只生成、审阅、注册，不能由代码或生成器自动执行。阶段提交保持本地，完成独立验收后再考虑推送。

## 目标

在现有 Resource、RBAC 和数据范围权限基础上，增加字段级 `visible/readable/writable/sensitive` 策略，并让后端响应、查询、导出、创建、更新、RBAC 配置、前端页面和生成器共享同一份字段权限契约。

## 技术边界

- 不新增依赖，不改变 PostgreSQL/Redis 配置；
- 后端字段裁剪和写入保护是唯一可信边界；
- 默认字段策略保持现有资源行为：visible/readable/writable 为 true，sensitive 为 false；
- 角色字段覆盖只能收紧 Manifest 策略，不能重新开放 Manifest 禁止的字段；
- 多角色采用任一角色明确允许即可，但仍受 Manifest 上限约束；
- 迁移生成后人工审阅、注册和执行；本计划代码不自动执行迁移；
- `admin:make-resource` 仍是唯一公共资源生成入口；
- Service 放在 `backend/app/services/rbac/`，Resource 与 Controller 保持现有边界。

## 验收重点

- 不可读字段不能通过列表、详情、搜索、排序、导出或直接 API 泄露；
- 不可写字段在创建和更新时返回稳定的 `422 FIELD_PERMISSION_DENIED`，不能通过额外 JSON 属性、空值或专用资源路径绕过；
- 敏感字段默认不参与搜索和 CSV 导出；
- 数据范围过滤先应用，再进行字段裁剪，字段权限不能扩大 `all/own` 范围；
- 多角色字段权限合并正确；
- 生成器和 OpenAPI 产物确定性一致；
- `go test ./... -count=1`、`npx vue-tsc -b` 和可用环境下的生产构建通过。

## Task 1：Resource 字段策略契约与解析服务

### 文件

- 修改：`backend/app/core/resource/registry.go`
- 修改：`backend/app/core/resource/registry_test.go`
- 新增：`backend/app/services/rbac/field_permission_service.go`
- 新增：`backend/app/services/rbac/field_permission_service_test.go`

### 接口

为 `resource.Field` 增加 `Visible`、`Readable`、`Writable`、`Sensitive`。提供统一策略类型和默认值归一化，不让调用方直接解释零值。

服务提供：

```text
Resolve(userID, permission, manifest) -> effective field policies
ReadableFields(manifest, policy, export) -> fields
ValidateQueryField(field, policy, operation) -> error
ValidateWritablePayload(payload, manifest, policy) -> filtered payload/error
ProjectRecord(record, fields) -> response record
```

角色覆盖在本任务先使用内存/空覆盖接口，数据库持久化放到 Task 3；服务必须能够区分 Manifest 默认策略和覆盖后的有效策略。

### TDD 步骤

1. 先写测试：零值默认策略、显式策略保留、不可读字段裁剪、不可写字段拒绝、敏感字段导出/搜索限制。
2. 运行 `go test ./app/core/resource ./app/services/rbac -run 'Test.*Field|Test.*Policy' -count=1`，确认在实现前失败。
3. 实现字段策略类型、默认值和纯函数服务；字段名必须来自 Manifest，不接受客户端字段名作为 SQL 标识符。
4. 重跑聚焦测试，确认通过。
5. 提交：`feat: add resource field permission contract`。

## Task 2：统一应用到资源读取、查询和写入

### 文件

- 修改：`backend/app/core/admin/controllers/resource_list_controller.go`
- 修改：`backend/app/core/admin/controllers/resource_detail_controller.go`
- 修改：`backend/app/core/admin/controllers/resource_export_controller.go`
- 修改：`backend/app/core/admin/controllers/resource_crud_controller.go`
- 修改：`backend/app/core/admin/controllers/resource_scope.go`
- 测试：现有 Resource 列表、详情、导出、CRUD 测试；必要时新增字段权限边界测试。

### 行为

- 数据范围过滤先应用，再按有效字段策略裁剪响应；
- 列表/详情只返回 `visible && readable` 字段；
- 导出排除 `sensitive` 和不可读字段；
- 搜索和排序仅接受 `visible && readable && !sensitive` 字段；非法字段使用安全默认或稳定校验错误；
- 创建/更新拒绝不可写字段，返回 `FIELD_PERMISSION_DENIED`；
- users、roles 等专用资源也经过统一投影和写入校验，再进入原有 Service 业务校验；
- 不改变资源、权限和数据范围错误的既有语义，避免泄露字段存在性。

### TDD 步骤

1. 先写失败测试：列表/详情/导出裁剪、敏感搜索排序阻断、创建/更新拒绝不可写字段、own 范围与字段裁剪组合。
2. 运行 `go test ./app/core/admin/controllers -run 'Test.*Field|Test.*Resource|Test.*Export' -count=1`，确认失败。
3. 接入 FieldPermissionService，统一处理 map、模型 Public 输出和导出行；禁止在各 Controller 复制字段策略逻辑。
4. 重跑聚焦测试和现有 Controller 测试。
5. 提交：`feat: enforce resource field permissions`。

## Task 3：角色字段覆盖持久化与 RBAC API

### 文件

- 新增：`backend/database/migrations/20260922000002_create_permission_role_fields_table.go`
- 修改：`backend/bootstrap/migrations.go`
- 修改：`backend/app/services/rbac/rbac_role_service.go`
- 修改：`backend/app/core/admin/controllers/rbac_controller.go`
- 修改：`backend/app/openapi/spec.go`
- 修改：`backend/app/openapi/spec_test.go`
- 修改/新增：现有 RBAC Controller 测试文件

### 数据模型

`permission_role_field` 使用 `role_id`、`permission_id`、`field_name` 唯一约束，保存 `readable` 与 `writable` 覆盖。删除角色权限时同步删除覆盖；删除字段覆盖不影响 Manifest。

### API

扩展现有角色权限替换接口，保留 `permission_ids` 和 `scopes`，新增：

```json
{
  "fields": {
    "10": {
      "email": {"readable": true, "writable": false}
    }
  }
}
```

响应返回已保存的字段覆盖。未知 permission、未知 field、Manifest 禁止重新开放的策略和非法布尔值均返回 `422 VALIDATION_ERROR`。系统角色保护规则保持不变。

### TDD 步骤

1. 先写失败测试：保存/读取覆盖、未知字段拒绝、不能突破 Manifest 上限、多角色有效策略合并、删除权限同步清理。
2. 运行 `go test ./app/services/rbac ./app/core/admin/controllers -run 'Test.*Field|Test.*Permission' -count=1`，确认失败。
3. 实现事务内覆盖替换和字段策略解析；缺少覆盖时继承 Manifest 默认值。
4. 更新 OpenAPI 和测试，运行 `go test ./app/openapi ./app/services/rbac ./app/core/admin/controllers -count=1`。
5. 人工审阅迁移并注册；只有 PostgreSQL 可用且用户确认本地服务已启动时才执行迁移验证。
6. 提交：`feat: manage resource field permissions in rbac`。

## Task 4：RBAC 页面、生成器和前端资源页面

### 文件

- 修改：`admin/src/generated/api.ts`
- 修改：`admin/src/modules/rbac/pages/RBACPage.vue`
- 修改：`admin/src/locales/zh-CN/rbac.json`
- 修改：`admin/src/locales/en-US/rbac.json`
- 修改：`backend/app/console/resource_generator_command.go`
- 修改：`backend/app/generator/spec.go`
- 修改：`backend/app/generator/render.go`
- 修改：`backend/app/generator/frontend_render.go`
- 修改：`backend/app/generator/*_test.go`
- 修改：`backend/app/openapi/spec.go`
- 修改：`docs/generator.md`、`docs/resource-engine.md`、`docs/roadmap.md`

### 行为

- RBAC 页面在选中权限后显示紧凑的字段 readable/writable 配置；未选择权限时不允许编辑；
- 资源列表、表单、详情使用后端返回的字段策略，不再各自推断；
- 生成器支持字段策略选项，至少覆盖 `sensitive` 和 `readonly`，默认输出确定且兼容；
- Golden File 覆盖默认资源和受限字段资源；
- README 说明字段权限、迁移审阅和验证命令；
- OpenAPI 与生成客户端同步，不手工维护冲突类型。

### TDD/验证步骤

1. 先写失败 Golden/前端测试，再实现解析和渲染。
2. 运行 `go test ./app/generator -count=1`，确认失败后再修复。
3. 运行 `npx vue-tsc -b`，修复类型问题。
4. 运行 RBAC、OpenAPI、Resource 前端测试。
5. 使用临时资源进行真实 Artisan smoke，确认受限字段在后端、前端、README 和迁移产物一致；不执行生成迁移，并清理临时产物。
6. 提交：`feat: complete resource field permission milestone`。

## 最终验证与审查

执行：

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.tmp-gocache-field'); go test ./... -count=1; Remove-Item -LiteralPath '.tmp-gocache-field' -Recurse -Force -ErrorAction SilentlyContinue
npx vue-tsc -b
npm run build
```

若 Vite 仍被共享 Windows 原生模块或服务环境阻断，必须记录实际错误，不得将类型检查结果表述为生产构建通过。

最终进行一次独立提交范围自审，检查：

- 字段权限是否在所有读取、查询、导出和写入路径执行；
- 专用 users/roles 路径是否绕过统一策略；
- 多角色合并是否意外扩大 Manifest 权限；
- 未执行迁移、Redis/PostgreSQL 状态和 GitHub 推送状态是否如实记录。
