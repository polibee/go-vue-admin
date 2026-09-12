# 数据库与生产配置

## 数据库选择

资源模块统一使用 GORM。数据库方言通过 DB_CONNECTION 选择：

支持分项环境变量，也支持完整的 `DB_DSN`。切换到 Supabase 时使用 PostgreSQL 的连接池连接串并启用 TLS。

    RESOURCE_PROVIDER=database
    DB_CONNECTION=mysql

或：

    RESOURCE_PROVIDER=database
    DB_CONNECTION=postgres

迁移和资源连接池应使用同一组连接环境变量。应用负责连接池生命周期，并在 shutdown 时关闭资源连接。

## 上线前配置

- 使用独立数据库和最小权限账号。
- 将 APP_KEY、数据库密码、支付密钥放入密钥管理系统或部署平台 Secret。
- 关闭开发 bootstrap 账号和本地默认密码。
- 明确迁移执行窗口和回滚策略。
- 先备份，再执行数据库迁移。
- 为 API、数据库和队列配置监控与告警。

## 数据一致性

生成的 Repository 不能混用 Goravel DB 和模块私有连接。事务边界应由 Service 层明确管理，跨资源写入必须有失败补偿或回滚策略。
