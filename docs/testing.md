# 测试与验收

## 测试分层

```text
Unit
Integration
HTTP/API
Frontend Component
E2E
Contract
Golden File
```

## 第一阶段优先覆盖

- Goravel Service 业务规则；
- Repository 查询边界；
- 认证和 RBAC；
- Resource 列表、过滤、排序和分页；
- OpenAPI Client 生成；
- Generator 快照；
- Migration 从空库执行；
- 前端权限和错误状态。

## API 验收

```text
启动应用 -> 注册路由 -> 生成/检查 OpenAPI -> 生成 TypeScript Client -> 前端类型检查 -> HTTP Contract Test
```

## 完成标准

一个标准 Resource 必须具备后端权限、Goravel ORM 查询、请求验证、API 路由、OpenAPI 描述、generated Client、Vue 列表和表单、空/错/加载状态以及单元或集成测试。

## 当前阶段 E2E 验收

- `go run . artisan migrate:status`：确认权限范围和字段权限迁移已执行；
- `backend/scripts/contract-smoke.ps1`：验证 OpenAPI、登录、Registry、列表、生成资源和全局搜索；
- 浏览器验收：登录、动态菜单、资源列表、创建表单、编辑表单、详情、访问控制和审计日志；
- HTTP 验收：资源详情、权限过滤、搜索、CSV 导出，以及关系不存在/无权限边界；
- 真实关系成功路径需要 Registry 中存在带 `Relations` 的业务资源 fixture，当前内置资源暂未提供该 fixture；
- 请求/响应审计脱敏不作为本阶段 Resource E2E 的阻塞项，单独维护其测试和验收边界。

## 审计日志生命周期

- 后台审计页可由具备 `admin.settings.manage` 权限的管理员手动清理；
- 手动清理通过 `POST /api/v1/admin/audit-logs/cleanup`，请求体使用 `retention_days`，范围为 1 到 3650 天；
- 定期清理使用 `go run . artisan admin:prune-audit-logs --days=365`，由 Laragon 或操作系统任务计划按需调度；
- 清理动作只删除早于保留窗口的记录，并返回删除数量；清理本身保留一条摘要审计事件；
- 清理失败不得通过 HTTP 请求影响其他业务数据。
