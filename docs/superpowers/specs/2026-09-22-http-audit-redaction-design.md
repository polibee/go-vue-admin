# HTTP 请求响应审计与脱敏设计

## 目标

在现有 PostgreSQL `audit_logs.metadata` JSONB 字段中记录管理 API 的请求与响应摘要，统一执行递归脱敏和大小限制，避免密码、令牌、Cookie、授权头和常见个人敏感字段进入审计日志。

## 边界

- 只记录 `/api/v1/` HTTP API，不记录静态文件、Scalar 文档、OPTIONS 和审计日志查询本身。
- 不保存原始 Cookie、Authorization、Refresh Token 或完整二进制/CSV 响应。
- JSON 请求体、查询参数、路径参数和 JSON 响应记录为结构化数据；非 JSON 响应只记录 content type、状态和截断摘要。
- 单个请求或响应审计内容最多 32 KiB，超出时保留脱敏后的截断标记。
- 审计失败不得改变业务响应，只写服务端错误日志。
- 现有业务审计事件继续保留，且同样经过统一脱敏。

## 脱敏规则

大小写不敏感地脱敏以下键：`password`、`password_confirmation`、`token`、`access_token`、`refresh_token`、`authorization`、`cookie`、`set-cookie`、`secret`、`api_key`、`client_secret`。

递归处理 map、数组和嵌套对象；匹配键的值统一替换为 `[REDACTED]`。未知字段默认保留，避免破坏业务排查价值。请求/响应原文不进入日志。

## 验收标准

1. 登录、资源 CRUD、搜索、导出和关系 API 产生 HTTP 审计记录。
2. 审计记录包含 method、path、status、request、response 和 truncated 信息。
3. 密码、Bearer Token、Cookie 和嵌套敏感字段在数据库和管理端详情中均不可见。
4. 审计接口自身不会产生无限递归审计。
5. 审计写入失败不影响原请求状态码和响应体。
6. Go 单元测试、全量测试、前端构建和浏览器端到端验收通过。
