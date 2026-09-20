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
