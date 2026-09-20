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
