# 生产级检查清单

## 当前判断

项目已经达到可持续开发和内部验收的 Foundation 阶段，具备认证、RBAC、通用 CRUD、GORM 多数据库选择、OpenAPI 生成、插件生命周期和本地浏览器验收基础。

它还不能在没有业务侧加固的情况下直接宣称所有生产场景就绪。生产上线前必须完成下面的环境和业务检查。

## 已具备

- HttpOnly Cookie Session 和 CSRF 保护边界。
- 后端授权优先，前端菜单隐藏只作为体验优化。
- GORM Repository 和 MySQL/PostgreSQL 配置选择。
- OpenAPI、Schema 和生成客户端的漂移检查。
- 资源 CRUD、模块/插件注册和本地 E2E 验收。
- GitHub Actions 手动触发，不会因每次推送自动消耗 CI。

## 上线前必须补齐

- 正式身份源、强密码策略、MFA 或企业 SSO。
- 密钥托管、轮换和审计，尤其是支付和证书私钥。
- HTTPS、反向代理、可信 Origin 和 Cookie 安全属性。
- 数据库备份、恢复演练、迁移回滚和高可用方案。
- API 限流、请求体限制、上传扫描和日志脱敏。
- 结构化日志、指标、错误追踪和健康检查告警。
- 业务级授权测试、越权测试和依赖漏洞扫描。
- 生产构建、版本发布、回滚和数据库变更审批。

## 本地发布检查

    go -C backend test ./...
    go -C backend vet ./...
    pnpm run lint
    pnpm run typecheck
    pnpm run test
    pnpm run build
    pnpm run openapi:generate
    git diff --exit-code -- contracts/openapi contracts/schemas admin/src/generated

文档站只记录检查项，不替代真实的安全评审、灾备演练和业务验收。
