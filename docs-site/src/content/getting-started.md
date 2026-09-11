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

## 生产提醒

不要把演示账号、APP_KEY 或数据库密码提交到仓库。生产环境应关闭 bootstrap 账号接口，并由部署系统注入密钥。
