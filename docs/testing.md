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
- 手动清理通过 `POST /api/v1/admin/audit-logs/cleanup`，支持 `retention`、`selected`、`filtered`、`all` 四种范围；
- `selected` 需要传 `ids`，`filtered` 使用当前 `action`/`user_id` 筛选条件且不受分页限制，`all` 删除全部日志；三种破坏性范围都必须传入 `confirmation: "DELETE"`；
- `filtered` 至少需要一个有效筛选条件，空筛选不会退化成全量删除；
- 定期清理使用 `go run . artisan admin:prune-audit-logs --days=365`；应用内默认每天 02:30 注册保留期清理任务，也可由 Laragon 或操作系统任务计划调用该命令；
- 清理动作返回删除数量，并在交互式后台清理完成后保留一条摘要审计事件；定期命令只执行保留期清理，不执行全量删除；
- 清理失败不得通过 HTTP 请求影响其他业务数据。

## 站内通知

- 通知迁移已审核并执行；可用 `go run . artisan migrate:status` 确认 `20260922000005_create_notifications_table` 为 `Ran`；
- 通知 API 仅允许登录用户访问自己的通知，支持分页、未读数量、单条已读和全部已读；
- AdminShell 顶部通知入口只支持站内相对路径跳转，不渲染 HTML，也不实现邮件、短信或 WebSocket 通道；
- 已在加载最新后端进程的浏览器实例中验收通知入口空状态，以及列表和未读数量 API 均返回成功；若页面仍显示接口失败，应先停止占用 3000/3001 的旧项目进程再重启后端，不要修改 Laragon 服务配置。
- 当前已接入用户创建、状态变更、密码重置、角色分配、角色权限变更、资源批量 Action、审计清理摘要和“撤销其他会话”通知；发布入口统一由通知 Service 提供，Controller 不直接写通知表；
- 系统公告的指定用户/角色投递尚未接入，需先补充公告受众字段和授权校验，再作为独立阶段实现。
