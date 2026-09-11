# 开发业务模块

## 模块结构

新模块应包含以下文件：

    modules/catalog/
    ├── module.yaml
    ├── backend/
    ├── admin/
    ├── README.md
    └── acceptance.md

模块必须显式声明 ID、版本、依赖和权限。模块不能通过私有路径访问另一个模块的实现。

## 生成模块

使用 Go 生成器，不依赖 Bash 或 PowerShell 的平台差异：

    go -C backend run ./cmd/admin-gen module catalog

生成后补充模块的业务定义、路由、权限和验收测试。不要把生成器输出直接当成完整业务实现。

## 菜单策略

一个业务入口直接提升为一级菜单。多个入口通过 owner 自动形成可折叠分组，也可以在导航声明中提供 group、groupLabel 和 order。

模块管理入口 /admin/modules 只用于模块生命周期和配置，不放业务订单、商品等页面。

## 交付检查

    pnpm run openapi:generate
    pnpm run typecheck
    pnpm run test
    pnpm run build
    bash scripts/module-check.sh modules/catalog
