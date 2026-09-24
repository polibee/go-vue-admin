# Go Vue Admin 开发文档

本项目是基于 **Goravel v1.18 + Vue 3 + shadcn-vue** 的通用后台管理平台。项目只建设后台平台能力，不重复实现 Goravel 已经提供的 Web 框架能力。

## 阅读顺序

- [architecture.md](./architecture.md)：总体架构和依赖方向
- [goravel-boundaries.md](./goravel-boundaries.md)：Goravel 能力边界
- [module-system.md](./module-system.md)：模块组织和注册
- [resource-engine.md](./resource-engine.md)：标准 CRUD 资源
- [frontend.md](./frontend.md)：Vue Admin 和 shadcn-vue 规范
- [i18n.md](./i18n.md)：中文、英文和模块化语言包
- [openapi.md](./openapi.md)：前后端 API 契约
- [generator.md](./generator.md)：代码生成器
- [plugin-system.md](./plugin-system.md)：插件系统设计预留
- [security.md](./security.md)：认证、授权和敏感数据
- [authentication.md](./authentication.md)：JWT 和 Refresh Token 认证
- [integrations.md](./integrations.md)：Goravel 驱动和扩展包集成
- [testing.md](./testing.md)：测试和验收
- [roadmap.md](./roadmap.md)：分阶段路线图

## 核心原则

```text
Goravel 有的能力不重复实现
标准 CRUD 优先使用 Resource
复杂业务使用普通 Vue 页面和 Goravel Service
前端直接使用 shadcn-vue 官方组件
生成代码与人工代码分离
权限必须由后端最终裁决
```

当前实现已经覆盖：Goravel v1.18 基础工程、后台壳层、JWT 与可撤销 Refresh Token、RBAC、标准 Resource CRUD、OpenAPI、TypeScript Client、PostgreSQL、Redis、全局搜索、关系字段、字段和数据范围权限、批量 Action、通知以及请求/响应审计脱敏。

资源生成默认采用 `Resource Manifest + 通用列表/表单/详情页`；只有复杂业务才通过模块边界覆盖专用页面。插件系统、文件上传和导入暂不作为当前范围。生产使用前仍需完成真实环境配置、普通角色端到端授权验收、备份、监控、回滚和发布流程验证。
