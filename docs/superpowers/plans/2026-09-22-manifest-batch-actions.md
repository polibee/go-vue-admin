# Manifest 驱动批量 Action 实施计划

> **执行方式：** 当前 checkout 内 Native 执行。每个任务先写失败测试，再实现、验证并提交；阶段完成后只保留本地提交，不推送 GitHub。

## 目标

将现有用户批量设置状态接口收敛到 Manifest 驱动的统一 Action 入口，并让 ResourceList 按 Manifest 和当前权限显示批量操作。第一阶段只实现受控的 `set-status` Handler，不开放任意字段批量更新、SQL 或脚本。

## 全局约束

- 不新增依赖，不执行数据库迁移；
- PostgreSQL/Redis 继续由 Laragon 管理，测试不得修改配置；
- 资源、Action、字段和权限均由服务端 Manifest 决定；
- Action Handler 按领域放置，不把业务逻辑堆入通用 Controller；
- 旧 `/api/v1/admin/users/status` 不保留为公开兼容入口；
- 所有阶段性验证通过后再形成里程碑提交，不推送 GitHub。

## 文件边界

后端：

- `backend/app/core/admin/actions/`：通用 Action Handler 接口、注册表、请求归一化和结果类型；
- `backend/app/modules/users/actions/`：用户 `set-status` Handler，复用现有 `UserService` 的状态校验和最后管理员保护；
- `backend/app/core/admin/controllers/resource_action_controller.go`：统一入口的认证、Manifest/Action 查找、公共参数校验和 Handler 调度；
- `backend/routes/web.go`：注册统一 Action 路由并删除旧用户状态公开路由；
- `backend/app/openapi/spec.go`、`spec_test.go`：统一 Action 的 OpenAPI 契约；
- 测试文件按上述领域目录放置。

前端：

- `admin/src/generated/api.ts`、`admin/scripts/generate-api-client.mjs`：统一 Action Client 和请求/结果类型；
- `admin/src/core/resource/pages/ResourceListPage.vue`：按 Manifest Action 生成选择列、批量工具栏和 `kind` 对应参数表单；去除 `resourceName === 'users'` 的批量分支；
- `admin/src/lib/resource-actions.ts` 或新的 `admin/src/core/resource/actions.ts`：Action 可执行判断和通用调用辅助；
- `admin/src/locales/*`：批量 Action、成功/失败汇总和错误文案。

文档：

- `docs/resource-engine.md`：更新统一 Action API 和移除旧用户批量接口的说明；
- `docs/roadmap.md`：记录批量 Action 阶段完成状态和限制；
- 生成器 README/契约文档：说明 Action 由 Resource Manifest 生成和审核，Handler 仍需后端显式实现。

## 任务 1：建立 Action 执行契约与 Handler 注册边界

### RED

先新增失败测试，覆盖空 ID、重复 ID、超过上限、未声明 Action、未知 Handler Kind、未知参数、逐 ID 结果和重复 Kind 注册。

运行：

```text
go test ./app/core/admin/actions ./app/modules/users/actions -count=1
```

预期先因类型、注册表和 Handler 不存在而失败。

### GREEN

实现 `ActionRequest`、`ActionResult`、`ActionFailure`、最大 ID 数量、ID 归一化、Handler 接口、按 `Kind` 查找的注册表以及 `set-status` 参数解析。审计元数据只记录 Action、数量和结果，不记录敏感请求内容。

提交：`feat: add manifest action execution contract`

## 任务 2：接入统一后端入口和权限/数据范围

### RED

新增 Controller/路由测试覆盖统一入口、Action 权限、`own` 数据范围、部分成功结果和旧 `/api/v1/admin/users/status` 路由不存在。

运行：

```text
go test ./app/core/admin/controllers ./app/routes ./app/services/rbac -count=1
```

### GREEN

实现 Controller：查找 Manifest 和 Action，校验权限、ID 数量和参数，逐 ID 应用 ResourceScopeService，调度 Handler，返回统一结果并写入管理审计。将用户批量状态逻辑迁移到 `modules/users/actions`，删除旧公开路由、旧专用 Client 方法和不再使用的 Controller 方法；保留现有 UserService 作为领域服务。

提交：`feat: execute manifest batch actions`

## 任务 3：同步 OpenAPI、Client 和文档

### RED/GREEN

先更新 OpenAPI 测试，要求出现统一 Action 路径、请求体和结果 Schema，同时旧 `/admin/users/status` 不再出现。然后同步 `spec.go`、生成 Client 脚本和输出、Resource Engine、Roadmap、生成器 README。

验证：

```text
go test ./app/openapi ./app/core/admin/controllers -count=1
npx vue-tsc -b
```

提交：`feat: publish manifest action contract`

## 任务 4：前端 ResourceList 通用批量体验

### RED

补充前端契约测试，验证非 users Resource 声明 Action 后出现选择列，没有 Action 或没有权限时不出现入口，`user-status` 显示状态参数表单，执行后显示汇总并刷新列表。

### GREEN

重构 `ResourceListPage.vue`：选择列由可执行批量 Action 决定；Action 菜单由 Manifest 过滤；参数表单由 `kind` 分派；统一调用 `resourceAction` Client；清理 `bulkSetUserStatus`、`resourceName === 'users'` 判断和旧专用状态调用；行级 `set-status` 也调用统一 Action API。

提交：`feat: render manifest batch actions`

## 任务 5：阶段验证与里程碑

运行：

```text
go test ./... -count=1
go build -o storage/codex-backend-actions-check.exe .
npx vue-tsc -b
npm run build
git diff --check
```

删除临时二进制，确认 PostgreSQL/Redis 配置未变更，确认迁移未执行。完成真实路由 Contract Smoke（不执行写入破坏性操作，`set-status` 使用明确测试账号和可恢复状态）后提交 `feat: complete manifest batch actions milestone`。

## 验收结果

- users 的 `set-status` 通过统一 Action API 工作；
- ResourceList 不再硬编码 users 批量入口；
- Manifest Action、权限、数据范围和 Handler 四层校验均有测试；
- 旧专用用户批量状态 URL、Client 方法和前端调用被移除；
- 全量后端测试、构建、前端类型检查和生产构建通过；
- 不新增依赖、不执行迁移、不推送 GitHub。
