# 认证设计

## 目标

项目是前后端分离架构，采用：

~~~text
短时效 JWT Access Token
+
HttpOnly Secure Refresh Token
+
Refresh Token Rotation
~~~

Goravel v1.18 同时支持 JWT 和 Session 驱动；本项目选择 JWT 作为 Admin API 的访问凭证，避免前端和后端必须共享页面 Session。

## Token 规则

~~~text
Access Token
- 短时效
- 前端仅保存在内存
- 请求使用 Authorization: Bearer

Refresh Token
- HttpOnly
- Secure
- SameSite 按部署拓扑配置
- 服务器端可撤销
- Redis 保存会话和轮换状态
~~~

禁止把长期 Access Token 或 Refresh Token 写入 localStorage。

## 流程

~~~text
POST /api/v1/auth/login
        ↓
验证账号和密码
        ↓
返回 Access Token
写入 Refresh Token Cookie
        ↓
访问受保护 API
        ↓
Access Token 过期
        ↓
POST /api/v1/auth/refresh
        ↓
验证并轮换 Refresh Token
        ↓
返回新的 Access Token
~~~

登出时撤销 Refresh Token，并清除 Cookie。密码使用 Goravel Hash 能力，不自行实现加密或哈希。

## 认证边界

- Goravel Auth 负责身份验证基础能力；
- Admin Core 负责 JWT Claims、Refresh Token 会话和撤销；
- RBAC 负责角色、权限和数据范围；
- Controller/Service 边界必须再次验证权限；
- 前端路由守卫只负责用户体验，不能替代后端授权。

## Redis

Redis 用于：

- Refresh Token 会话；
- Token 撤销；
- 轮换重放检测；
- 登录限流；
- 验证码和临时状态。

## 失败场景

必须处理：

- Access Token 过期；
- Refresh Token 过期；
- Refresh Token 重放；
- 用户被停用；
- 用户角色被撤销；
- Redis 不可用；
- 多设备登录；
- 主动退出所有设备。

## API

~~~text
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
POST /api/v1/auth/logout-all
GET  /api/v1/auth/me
~~~

## 当前实现状态

第一条可运行垂直链路已完成：

- `backend/app/models/user.go` 用户模型；
- `users` 表迁移；
- Goravel JWT 登录、当前用户、刷新和登出接口；
- Goravel Hash 密码校验；
- 开发管理员 Seeder：`admin@example.com / Admin123!`；
- 前端登录页、内存 Access Token 和路由守卫。
- `roles`、`permissions`、`role_user`、`permission_role` RBAC 表；
- 管理员用户、角色、权限列表接口：
  `GET /api/v1/admin/users`、`GET /api/v1/admin/roles`、`GET /api/v1/admin/permissions`；
- `AdminUser` Seeder 创建 `super-admin` 及基础管理权限；
- 角色创建、详情、编辑、删除和权限分配接口；
- 用户角色查询和绑定接口；
- 基于 `admin.users.view`、`admin.roles.manage`、`admin.permissions.manage` 的后端细粒度授权；
- API 错误统一返回稳定 `code`，前端按模块语言包渲染中文或英文。

当前实现已接入独立的随机 Refresh Token：登录写入 HttpOnly Cookie，Redis 只保存 Token 哈希和用户 ID，刷新时一次性消费并轮换，重放返回 401，登出会撤销并清除 Cookie。生产环境仍需按部署拓扑配置 HTTPS、SameSite 和跨来源 Cookie 策略；`logout-all`、审计和登录限流仍属于后续增强。系统角色 `super-admin` 受到保护，最后一个具备管理权限的管理员不能被移除。

## 测试

必须验证登录、刷新、轮换、撤销、过期、重放检测、停用用户、权限变更和跨来源 Cookie 配置。
