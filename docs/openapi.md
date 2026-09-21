# OpenAPI 与前后端契约

## 契约流程

```text
Goravel Controller/Request/Response
        ↓
OpenAPI Schema
        ↓
TypeScript API Client
        ↓
Vue Admin
```

OpenAPI 是前后端契约来源，前端不得另行维护重复的 TypeScript Model。

项目采用 OpenAPI 3.x，后端使用 `swaggest/openapi-go` 组织和生成 Schema，前端使用 OpenAPI TypeScript Generator 生成 Client。`goravel-crud` 自带的 Swagger 方案不作为项目契约。

## API 规范

统一使用 `/api/v1/`：

```text
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PATCH  /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

错误响应必须包含稳定机器码，例如 `VALIDATION_ERROR`、`AUTH_REQUIRED`、`PERMISSION_DENIED`、`RESOURCE_NOT_FOUND`、`RESOURCE_CONFLICT` 和 `INTERNAL_ERROR`。前端不得根据错误文案判断业务状态。

Scalar 只作为开发和受控环境的 API 浏览器，不是 API Contract Source。生产环境默认关闭或通过鉴权、网关和内网策略限制访问。

生成 Client 存放在 `admin/src/generated/`，标记为 `DO NOT EDIT`。后端相关 OpenAPI 适配代码位于 `backend/app/core/openapi/`。

当前已落地开发契约入口：`GET /api/openapi.json`。契约源位于 `backend/app/openapi/spec.go`，并由 `spec_test.go` 校验已实现的认证、资源和用户状态接口。该入口暂不代表生产公开文档；生产暴露前仍需增加鉴权或网关限制。

本地环境另提供 `GET /api/docs` Scalar 浏览页；该路由仅在 `APP_ENV` 非 `production` 时注册，并始终从本地 `/api/openapi.json` 读取契约。

前端已接入首版契约 Client：`admin/src/generated/api.ts`。该文件标记为生成产物，当前覆盖登录、当前用户、资源清单、资源列表和批量用户状态操作；后续生成器稳定后替换其生成实现。

Client 现在可通过 `pnpm generate:api` 从本地运行服务的 `/api/openapi.json` 重新生成；可用 `OPENAPI_BASE_URL` 指定契约服务地址。生成脚本为仓库内置无依赖脚本，不修改依赖锁文件。

本地服务启动后可运行 `powershell -File backend/scripts/contract-smoke.ps1`，执行只读契约烟测：OpenAPI、Scalar、登录、资源清单和 users 列表。脚本不会创建、删除或修改业务数据。
