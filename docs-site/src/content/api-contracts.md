# API 契约与客户端

## 生成流程

后端 OpenAPI 构建源生成三类产物：

    contracts/openapi/openapi.json
    contracts/schemas/*.json
    admin/src/generated/api/*

执行：

    pnpm run openapi:generate

生成文件不要手工修改。需要新增字段时修改 Go OpenAPI 构建源、Schema 或路由元数据，再重新生成并检查差异。

## 管理端调用规则

业务页面使用生成的客户端或 ResourceDataProvider。不要在 Vue 页面中直接调用 fetch、axios 或拼接未声明的 API 路径。

## 外部集成

浏览器只适合调用不含敏感逻辑的用户侧接口。支付签名、密钥、证书、文件处理和 webhook 验证应由业务系统后端调用本项目 API。

当前项目提供通用 OpenAPI 契约和 TypeScript 管理端客户端。未来具体业务上线后，再按稳定业务契约生成 Go、PHP、Java、Python 或其他语言 SDK。

## 文档入口

后台内置的 Scalar 页面用于本地和受保护环境查看当前契约。公开开发文档站只说明稳定集成边界，不暴露后台账号、Cookie 或内部配置。
