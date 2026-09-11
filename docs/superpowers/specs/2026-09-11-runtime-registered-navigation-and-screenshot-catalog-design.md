# Runtime-Registered Navigation and Screenshot Catalog Design

## Goal

让模块和插件的业务页面能够按照声明自动注册到后台菜单与路由，同时保留独立的模块管理页和插件管理页。管理入口负责平台控制，业务入口负责实际功能，避免所有页面继续堆在“模块与插件”页面中。

## Decisions

### 1. 管理入口与业务入口分离

保留两个平台管理入口：

- `/admin/modules`：模块状态、版本、资源、依赖和配置。
- `/admin/plugins`：插件启停、版本、权限、依赖和配置。

模块或插件声明的业务菜单显示为独立侧边栏项目，不嵌套在管理页中。旧 `/admin/extensions` 继续作为兼容入口，但不再出现在主导航。

业务菜单采用自动分组策略：单个 owner 只有一个业务入口时直接提升为一级菜单；同一 owner 有多个入口时自动创建可折叠分组；Manifest 可以通过 `group`、`groupLabel` 和 `order` 覆盖默认分组。模块/插件管理菜单仍保持独立，配置页面不计入业务菜单数量。

### 2. Manifest 是声明源

前端模块定义和插件 Manifest 负责声明：

- 页面路由和页面组件
- 菜单 ID、名称、路径和权限
- 资源定义
- 权限和依赖

模块和插件还可以声明配置入口，但配置内容不进入平台管理页的通用表单：

- `route`：专属配置页面路由
- `label`：配置入口名称
- `permission`：配置所需权限，建议使用 `plugins.configure` 或领域权限
- `schema`：业务配置模型标识
- `secretFields`：仅用于脱敏和提交策略提示，不能让前端回显明文

`/admin/modules` 和 `/admin/plugins` 只展示配置状态并提供入口；支付密钥、证书等由后端配置 API 接收、加密存储和审计，业务插件自己的页面负责表单和连接测试。

后端继续负责最终鉴权。前端菜单隐藏只属于 UX 优化，不能替代 API 授权。

后续应为后端模块补充可序列化的模块 Manifest/Info，使后端模块目录和前端模块定义共享稳定 ID、名称和版本语义；本阶段不引入动态下载或第三方代码执行。

### 3. 生命周期

模块：应用启动时注册路由、菜单和资源，并保持可用。

插件：

- 启用：执行 `setup`，挂载菜单、路由和资源。
- 停用：执行 `teardown`，卸载菜单、路由和资源。
- 依赖未满足时不能启用。

运行时注册必须可重复执行且不能产生重复菜单、重复路由或重复资源。插件停用后，其菜单和路由不能继续访问；后端 API 仍必须独立校验权限和插件状态。

### 4. Runtime Bridge

新增一个前端运行时桥接层，统一接收 `ModuleRegistry` 和 `PluginRegistry` 的注册结果，并同步到：

- `NavigationRegistry`
- Vue Router 动态路由
- `resourceRegistry`

桥接层只负责运行时装配，不把注册逻辑散落到页面组件中。模块和插件管理页只调用管理 API，不直接操纵导航注册表。

### 5. Screenshot catalog

删除旧的临时或低清截图，只保留一组可复现的功能截图：

- `dashboard.png`
- `modules.png`
- `plugins.png`
- `users.png`
- `roles.png`
- `permissions.png`
- `settings.png`
- `media.png`
- `audit.png`
- `api-docs.png`

截图使用 `1440×900` 视口和 `deviceScaleFactor=2`。不截图登录页，不把账号、密码、Cookie 或临时浏览器状态写入仓库。README 只引用这组截图。

## Compatibility

- `/api/extensions`、`/admin/extensions` 保留兼容，不作为新代码的首选入口。
- 现有 `/api/modules`、`/api/plugins` 继续提供管理目录和状态操作。
- 现有静态核心菜单保持可用，动态菜单在其后合并。

## Verification

必须覆盖：

1. 模块注册后自动出现在导航并可打开其页面。
2. 插件启用后出现菜单/路由，停用后菜单消失且路由不可访问。
3. 重复注册不会产生重复菜单、路由或资源。
4. 权限不足时菜单隐藏，直接访问 API 仍返回拒绝。
5. 旧扩展入口仍可重定向或访问兼容详情。
6. Playwright 验收模块管理、插件管理、业务菜单和插件停用流程。
7. 所有公开截图达到目标像素尺寸，README 不引用已删除截图。

## Non-goals

本阶段不实现：

- 从网络下载或热加载插件
- 第三方插件沙箱和签名校验
- 自动生成后端菜单数据库
- 删除旧扩展 API
- 多租户级别的菜单定制
