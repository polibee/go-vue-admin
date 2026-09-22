# Manifest Action 完整契约补齐实施计划

## 目标与约束

在现有批量 Action 基础上补齐 `batch`、`payload`、skipped 和生成器契约。保持当前目录边界、无新依赖、不执行迁移，阶段完成后只保留本地里程碑提交。

## 任务 1：Core Action 契约

### RED

新增测试覆盖：

- Manifest Action 的 `batch` 和 `payload` 元数据；
- 只有 `batch=true` 的 Action 可以进入批量候选；
- 缺失或未知 payload 契约被拒绝；
- Action 结果可以区分 failed 和 skipped；
- 旧 `params` 字段不被接受。

### GREEN

修改：

- `backend/app/core/resource/registry.go`：扩展 `resource.Action`；
- `backend/app/core/admin/actions/`：扩展请求、结果和 Handler 注册校验；
- `backend/app/modules/users/actions/`：将 `set-status` 定义为 `batch=true`、`payload=user-status`。

提交：`feat: complete manifest action contract`

## 任务 2：后端执行边界

### RED

新增 Controller 测试覆盖：

- `batch=false` Action 无法通过批量入口；
- `payload` 契约缺失或不匹配返回稳定校验错误；
- 范围外/无权限 ID 进入 skipped；
- 成功、失败、skipped 数量和明细准确；
- `set-status` 使用 `payload`，不接受 `params`。

### GREEN

修改 `resource_action_controller.go`：

1. 校验 Manifest Action 声明和 `Batch`；
2. 校验 `Payload` 是否是已注册契约；
3. 逐 ID 应用权限和数据范围；
4. 将越权 ID 放入 skipped，不交给 Handler；
5. 统一返回 succeeded/failed/skipped；
6. 保持领域事务由 Handler/Service 控制，不在 Manifest 中伪造跨资源事务。

提交：`feat: enforce complete batch action contract`

## 任务 3：Generator、OpenAPI 与 Client

### RED/GREEN

修改：

- `backend/app/generator/spec.go`：支持 Action 定义及 `batch/kind/payload`；
- `backend/app/generator/render.go`、Golden File：生成稳定 Action 元数据；
- `backend/app/generator` README 模板：输出 Handler 注册清单和 payload 契约说明；
- `backend/app/openapi/spec.go`、测试：同步 Action Schema、请求和响应；
- `admin/scripts/generate-api-client.mjs`、`admin/src/generated/api.ts`：同步类型化 Client；
- `announcements` 示例资源：验证生成结果。

默认标准 CRUD Action 使用 `batch=false`；生成器不生成任意业务 Handler。

提交：`feat: generate complete manifest action metadata`

## 任务 4：前端与文档

### RED/GREEN

修改：

- `admin/src/lib/resource-actions.ts`：批量候选必须满足 `batch=true`；
- `admin/src/core/resource/pages/ResourceListPage.vue`：使用 `payload`，展示 succeeded/failed/skipped 汇总；
- 未知 Kind 或 Payload 不显示可执行表单；
- 更新 `docs/resource-engine.md`、`docs/generator.md`、`docs/roadmap.md`。

提交：`feat: consume complete manifest action contract`

## 任务 5：验证与里程碑

```text
go test ./... -count=1
go build -o storage/codex-backend-action-contract-check.exe .
npx vue-tsc -b
npm run build
git diff --check
```

删除临时二进制，确认迁移未执行、PostgreSQL/Redis 配置未改变，完成 Generator Golden 和示例资源 Contract Smoke 后提交：

`feat: complete manifest action contract milestone`

不推送 GitHub。
