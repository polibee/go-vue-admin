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

项目采用 OpenAPI 3.x。当前契约由 `backend/app/openapi/spec.go` 明确维护，前端使用生成 Client。`goravel-crud` 自带的 Swagger 方案不作为项目契约。

## API 规范

认证接口统一使用 `/api/v1/auth/`，后台管理接口统一使用 `/api/v1/admin/`：

```text
GET    /api/v1/admin/users
POST   /api/v1/admin/users
GET    /api/v1/admin/users/{id}
PUT    /api/v1/admin/users/{id}
DELETE /api/v1/admin/users/{id}
```

错误响应必须包含稳定机器码，例如 `VALIDATION_ERROR`、`AUTH_REQUIRED`、`PERMISSION_DENIED`、`RESOURCE_NOT_FOUND`、`RESOURCE_CONFLICT` 和 `INTERNAL_ERROR`。前端不得根据错误文案判断业务状态。

Scalar 只作为开发和受控环境的 API 浏览器，不是 API Contract Source。生产环境默认关闭或通过鉴权、网关和内网策略限制访问。

生成 Client 存放在 `admin/src/generated/`，标记为 `DO NOT EDIT`。后端相关 OpenAPI 适配代码位于 `backend/app/core/openapi/`。

当前已落地的 Admin 契约规范入口：`GET /api/openapi/admin.json`。契约源位于 `backend/app/openapi/spec.go`，并由 `spec_test.go` 校验已实现的认证、资源和用户状态接口。该入口暂不代表生产公开文档；生产暴露前仍需增加鉴权或网关限制。

本地环境另提供 `GET /api/docs` Scalar 浏览页；该路由仅在 `APP_ENV` 非 `production` 时注册，并始终从本地 `/api/openapi/admin.json` 读取 Admin 契约。

前端已接入首版契约 Client：`admin/src/generated/api.ts`。该文件标记为生成产物，当前覆盖登录、当前用户、资源清单、资源列表、资源详情、资源新增/更新/删除和批量用户状态操作；通用资源的字段、选项、列和动作元数据也在 OpenAPI Schema 中声明。

Client 现在可通过 `pnpm generate:api` 从本地运行服务的 `/api/openapi/admin.json` 重新生成；可用 `OPENAPI_BASE_URL` 指定契约服务地址。生成脚本为仓库内置无依赖脚本，不修改依赖锁文件。

本地服务启动后可运行 `powershell -File backend/scripts/contract-smoke.ps1`，执行只读契约烟测：OpenAPI、Scalar、登录、资源清单和 users 列表。脚本不会创建、删除或修改业务数据。
# OpenAPI 契约边界

当前 `/api/openapi/admin.json` 是已实现的 Admin Core 契约，覆盖认证、Resource、RBAC、通知和审计等后台 API；它不是未来 C 端 API 的全站契约。当前不会注册 `/api/v1/app/*` 路由，也不会提前创建 C 端业务接口。

未来新增 App API 时，应单独发布 App 契约，并通过路由覆盖测试保证新增 Controller 不会遗漏文档。不要把面向用户的 API 直接追加到 Admin Resource 文档中。
