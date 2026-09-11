# 快速开始

## 环境要求

- Node.js 24 或更高版本
- pnpm 10
- Go 1.27
- MySQL 8 或 PostgreSQL

## 开发环境启动

在仓库根目录安装管理端依赖：

    pnpm install

为后端提供开发密钥和演示账号：

    APP_ENV=local
    APP_KEY=replace-with-a-32-character-dev-key
    AUTH_BOOTSTRAP_EMAIL=admin@example.com
    AUTH_BOOTSTRAP_PASSWORD=test-only-password
    RESOURCE_PROVIDER=memory

启动后端：

    go -C backend run .

启动管理端：

    pnpm --dir admin run dev -- --host 0.0.0.0

默认地址是 http://127.0.0.1:5173/login，开发代理会探测后端端口。

管理端默认使用 HTTP Resource Provider，因此用户、角色和权限的 CRUD 会由后端保存，浏览器刷新不会丢失刚刚创建的数据。只有在明确需要离线演示时，才将 VITE_RESOURCE_PROVIDER 临时设置为 memory。

## 语言和主题

管理端默认显示简体中文，右侧边栏可以切换 English。语言选择保存在浏览器本地存储中，业务模块应通过 vue-i18n 消息键提供文案，不要在组件中写语言分支。

主题由两层组成：浅色/深色/跟随系统，以及 shadcn 默认、Semi Design、微信风格三套语义色板。色板只覆盖 CSS 变量，不替换 shadcn-vue 组件和底层交互实现。

## 生产提醒

不要把演示账号、APP_KEY 或数据库密码提交到仓库。生产环境应关闭 bootstrap 账号接口，并由部署系统注入密钥。
