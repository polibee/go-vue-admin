# 总体架构

## 定位

```text
backend/             Goravel v1.18 后端
admin/               Vue 3 + shadcn-vue 前端
Admin Core           通用后台平台能力
Business Modules     CMS、Trading、Payment 等业务能力
```

依赖方向只能是 `Business Module -> Admin Core -> Goravel`。Admin Core 不得依赖具体业务模块。

## 分层

Goravel Foundation 直接提供路由、中间件、HTTP、认证、授权、验证、ORM、迁移、工厂、事务、缓存、事件、队列、文件、邮件、日志、加密、哈希、Session、本地化和测试。

Admin Core 只包含模块注册、Resource Registry、导航注册、平台权限、审计记录、API 契约、后台设置和生成器接口。后端多语言直接使用 Goravel；前端使用 vue-i18n。插件系统属于未来的扩展边界，第一阶段只预留 Contract，不进入运行时。

每个业务模块拥有自己的模型、服务、路由、权限、资源和页面，不把业务规则放进 Core。

## 后端调用链

```text
Route -> Controller -> Request/Authorization -> Service -> Repository -> Goravel ORM -> Database
```

Controller 不直接查询数据库；Service 不生成 HTTP 响应；Repository 不承担权限判断。

## 目录基线

```text
backend/
├── app/
│   ├── core/
│   │   ├── contracts/
│   │   ├── module/
│   │   ├── resource/
│   │   ├── navigation/
│   │   ├── authorization/
│   │   ├── audit/
│   │   └── openapi/
│   ├── modules/
│   │   ├── users/
│   │   ├── roles/
│   │   ├── permissions/
│   │   ├── settings/
│   │   └── audit/
│   ├── providers/
│   ├── http/
│   └── models/
├── bootstrap/
├── database/
├── routes/
├── config/
├── lang/
├── storage/
├── artisan
├── go.mod
└── main.go

admin/
├── src/
├── public/
├── package.json
└── vite.config.ts
```

只有跨模块、稳定且没有明确业务归属的能力才能进入 `core`。
# Admin/App API 边界

当前仓库明确区分两个未来 API 面：后台管理使用 `/admin/...` 页面、`/api/v1/admin/...` API 和 `admin.<resource>.<action>` 权限；面向终端用户使用根路径页面、`/api/v1/app/...` API 和 `app.<resource>.<action>` 权限。

当前只实现 Admin 面。领域 Service 可以复用，但 Controller、路由中间件、OpenAPI 契约和授权策略按 API 面分开。`admin:make-resource` 目前拒绝生成 App Resource，避免生成 App 权限却注册到 Admin Registry 的半套实现。
