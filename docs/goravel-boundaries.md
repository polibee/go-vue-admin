# Goravel 能力边界

项目基于 Goravel v1.18.x。数据库驱动从 Goravel 框架中独立安装：默认使用 `goravel/postgres`，需要 MySQL 兼容时再安装 `goravel/mysql`。

## 直接复用

| 能力 | 实现方 |
|---|---|
| HTTP、路由和中间件 | Goravel |
| Controller、Request、Response | Goravel |
| ORM、关系、事务、软删除 | Goravel ORM |
| Migration、Seeder、Factory | Goravel |
| Session、认证、授权 | Goravel 及其官方能力 |
| 验证、Hash、Crypt | Goravel |
| Cache、Event、Queue、Scheduler | Goravel |
| 文件、邮件、日志 | Goravel |
| Artisan 和测试 | Goravel |
| PostgreSQL/MySQL 驱动 | 对应 Goravel Driver Package |

## 本项目新增

| 能力 | 实现方 |
|---|---|
| 用户、角色、权限后台模型 | Admin Core/Platform Modules |
| 菜单和模块注册 | Admin Core |
| Resource Manifest | Admin Core |
| 标准 CRUD 页面 | Vue Admin |
| OpenAPI 和 Client 生成 | 项目工具链 |
| 审计业务记录 | Admin Core |
| CRUD Generator | 项目 CLI |

## 禁止重复实现

禁止创建第二套 ORM、数据库连接池、路由系统、验证器、事件总线、缓存抽象、任务队列、密码哈希工具或日志框架。

如果 Goravel 当前能力不足，优先写小型适配器；只有稳定且跨多个模块的新增能力才进入 Admin Core。业务默认使用 Goravel ORM、Model、Query、Relationship、Scope 和 Transaction。
