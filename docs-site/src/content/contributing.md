# 贡献与发布

## 修改边界

一个模块或一项基础设施能力使用一个独立提交。不要把生成输出、无关 UI 清理和业务功能混在同一提交。

提交前运行：

    pnpm run lint
    pnpm run typecheck
    pnpm run test
    pnpm run build
    pnpm run openapi:generate
    git diff --exit-code -- contracts/openapi contracts/schemas admin/src/generated
    bash scripts/module-check.sh

## 文档站发布

文档站使用 .github/workflows/developer-docs-pages.yml 手动发布。普通推送不会自动部署，也不会把文档站构建混入后台构建。

在 GitHub 仓库设置中将 Pages 的 Source 设为 GitHub Actions，然后手动运行 Developer Docs Pages workflow。站点只上传 docs-site/dist。

## 公开内容边界

公开文档可以包含模块 API、Manifest、配置约束和集成示例。不要提交账号密码、Cookie、真实密钥、内部网络地址、生产日志或未脱敏截图。
