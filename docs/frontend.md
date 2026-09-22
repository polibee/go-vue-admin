# Vue Admin 前端规范

## 技术栈

```text
Vue 3
TypeScript
Vite
Vue Router
Pinia
shadcn-vue
Tailwind CSS
TanStack Table
VueUse
```

## shadcn-vue 原则

项目直接使用 shadcn-vue 官方组件和官方视觉语言，不另起一套 Filament 风格，也不重新设计基础 UI。

```text
components/ui/       shadcn-vue 官方组件
components/admin/    页面布局和业务组合组件
modules/             业务页面
```

禁止自定义或修改 Button、Input、Dialog、Table、Select 等基础组件；禁止引入第二套 UI Framework；禁止为每个页面创造不同的基础控件风格。components/admin 只能组合官方组件，不能重新实现基础控件。

## 页面和模块

后台壳层统一提供 Sidebar、Header、Breadcrumb 和 Page Content。标准页面包括 Resource List、Resource Create/Edit、Resource Detail 和 Custom Business Page。自定义页面可以自由布局，但必须复用 shadcn-vue 官方组件和已有 API、权限、状态处理。

资源导航采用“单资源单入口”约定：每个 Resource Manifest 只生成一个规范列表入口，资源列表页只展示当前资源，不再提供跨资源切换器。Manifest 的 `navigation.group`、`navigation.order` 和 `navigation.hidden` 决定侧边栏分组、顺序和可见性；内置用户、角色、权限属于 `system`，生成资源默认属于 `business`。Dashboard 卡片和全局搜索是发现入口，不替代侧边栏，也不得重复生成资源菜单。

```text
admin/src/modules/users/
├── index.ts
├── routes.ts
├── permissions.ts
├── resource.ts
├── pages/
├── components/
└── composables/
```

业务代码优先使用 generated API Client，禁止在页面中散落重复的 `fetch`/`axios` 请求。

## 多语言

所有语言包统一放在 `admin/src/locales/`，按 `locale/namespace` 拆分，并由 vue-i18n 加载。详细规则见 [i18n.md](./i18n.md)。插件未来通过宿主提供的 vue-i18n 接入点声明自己的命名空间，不能覆盖 Core 或其他插件。
