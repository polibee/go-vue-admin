# Roadmap

## Phase 0：架构冻结

使用 Goravel Installer 创建 `backend/`，固定 Goravel v1.18.x、目录和模块注册，固定 shadcn-vue 官方组件约束，确定 Goravel Localization、vue-i18n 和中英文 Key 完整性检查。

## Phase 1：Admin Shell

已完成：Vue Admin 启动、shadcn-vue 官方组件全集、登录页、认证路由守卫、Sidebar、Header、Breadcrumb、错误页、加载页、Empty 状态和中英文切换。

## Phase 2：认证与 RBAC

基础版本已完成：用户表、JWT 登录、当前用户、刷新、登出、logout-all、Goravel Hash、中英文认证语言包、角色/权限表、管理员 RBAC 列表接口、角色 CRUD、权限分配、用户角色绑定、细粒度后端权限校验、初始 super-admin 数据和用户状态管理（active/disabled/locked）已完成；独立 Refresh Token Cookie、PostgreSQL 权威会话、Redis 镜像、自动故障切换、一次性轮换和重放检测已完成；认证审计日志、服务端筛选分页和前端详情查看已完成。菜单和更细粒度 Action 权限仍属于后续增强。

## Phase 3：Resource Engine

Resource Registry、字段、列、标准列表查询（分页、搜索、排序、用户状态筛选）、用户批量设置状态、声明式用户 Custom Action 和系统概览 Custom Page 已完成；users 资源已补齐创建、编辑、详情和删除闭环，metadata-driven 用户/角色表单已接入。

## Phase 4：OpenAPI 与 SDK

已完成第一步：入库 OpenAPI 3.0.3 契约源、开发访问入口和本地 Scalar 浏览页，覆盖当前认证、资源和用户状态接口；前端类型化 Client 已接入登录、当前用户、资源清单、资源列表和批量状态操作，并可由仓库内置无依赖脚本重新生成；契约单元测试与只读 Contract Smoke 已完成。

## Phase 5：Generator

Module、Resource、CRUD、权限、菜单和测试骨架生成，以及 Golden File Test。

## Phase 6：可选扩展

设置中心、文件管理、导入导出、Dashboard、全局搜索、插件 SDK 和 AI 辅助生成。

## Phase 7：插件系统

- Plugin Manifest；
- 插件发现和依赖解析；
- 插件状态机；
- Plugin Manager；
- 插件启用、停用和卸载 API；
- 插件 Capability 和 Locale Contract；
- 隔离进程运行时。

## 第一版明确不做

动态加载 Go 二进制插件、全能低代码 Schema、AI 运行时依赖、第二套 ORM/事件/缓存/任务系统，以及为每个页面设计独立视觉体系。

第一阶段只写插件设计和接口预留，不开发插件安装、执行、启停、卸载或插件市场。
