# v1 Foundation 本地发布记录

## 发布范围

v1 Foundation 交付模块化后台平台的基础能力：

- Cookie/HttpOnly Session 登录、登出、当前用户和 RBAC；
- Settings、Media、Audit 管理模块；
- 通用 Resource Engine，以及 users、roles、permissions 资源页面；
- OpenAPI 契约、生成的 TypeScript 客户端和 `OpenApiDataProvider`；
- 模块运行时、模块生成器、内置插件注册表和插件 SDK 契约；
- shadcn-vue UI 基线与本地仓库级验收门禁。

## 本地验收命令

在仓库根目录执行：

```bash
pnpm run release:check
```

该命令覆盖后端测试、Admin lint/typecheck/组件测试/构建、OpenAPI 生成差异、模块契约和 UI 兼容性检查。E2E 需要本地后端和 Admin 服务已启动时显式开启：

```bash
RELEASE_CHECK_E2E=1 pnpm run release:check
```

本次发布前已单独执行完整 Playwright E2E，结果为 5/5 通过，覆盖登录会话、RBAC 导航、通用资源 CRUD、筛选/排序/批量删除，以及用户/角色/权限资源路由。

## 发布边界

- GitHub Actions 保持 `workflow_dispatch` 手动触发，不启用 push 或 pull request 自动 CI。
- 本次 E2E 使用内存资源数据。用户、角色、权限资源路由仍使用内存仓储；仓库已包含 MySQL 示例资源仓储，但本次验收不代表 MySQL 持久化验证通过。
- 插件 SDK 当前定义注册、启用/禁用和生命周期契约，不执行不受信任的动态代码、包安装、签名验证或沙箱隔离。
- 不引入第二套 UI 框架，不改变既有 shadcn-vue 的颜色、圆角、阴影、间距或字体基线。

## 验收状态

本记录随 `chore(release): verify v1 foundation` 提交。后续模块必须继续通过 `scripts/release-check.sh`，并在模块完成后单独提交；远程推送不会触发 GitHub 自动 CI。

2026-09-10：`pnpm run release:check` 返回 `release-check: passed`，包括后端测试、前端 18 个测试文件 / 36 个测试、生产构建、OpenAPI 无差异、模块契约和 UI 边界检查。`bash scripts/release-check-test.sh` 与 `git diff --check` 通过。E2E 在同次修复中单独执行，5/5 通过，未在此次 release-check 中重复执行。

本次已补充并通过用户 CRUD、模块注册、插件启用/停用的专门浏览器验收；完整 Playwright 回归为 7/7 通过。
