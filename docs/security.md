# 安全规范

认证回答“用户是谁”，授权回答“用户能做什么”，两者不得合并成无法测试的巨型中间件。

本项目的前后端分离认证采用短时效 JWT Access Token 和 HttpOnly Refresh Token。详细流程见 [authentication.md](./authentication.md)。

后端必须验证身份认证、路由权限、Resource 权限、Action 权限和数据范围。隐藏菜单不等于实现权限。当前 RBAC 管理接口通过 Goravel 路由中间件和 Service 层双重校验，拒绝响应会中止后续路由链。

密码只能使用 Goravel 的成熟 Hash 能力保存，禁止明文和可逆加密。JWT、Refresh Token、数据库密码、Provider Key 和私钥不得提交 Git。敏感配置使用环境变量或 Secret Manager，仓库只保留 `.env.example`。

审计记录操作人、模块、资源、动作、资源 ID、结果、时间和请求关联信息。密码、Token、Secret、私钥等字段禁止进入日志和审计快照。

必须配置 CORS、Trusted Hosts、请求验证、速率限制扩展点、安全响应头和错误脱敏。文件上传、导入和批量操作必须经过权限、大小、类型、Schema 和 Service 校验，禁止直接写数据库。

RBAC 约束：系统角色 `super-admin` 不允许删除或修改权限；角色 ID、权限 ID 必须在 Service 层验证；最后一个具备管理权限的管理员不能解除自身或被其他操作移除。错误响应使用稳定错误码，避免把后端内部英文消息直接展示给用户。
