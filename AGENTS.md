# 项目开发约束

## 本机开发服务

- `D:\laragon` 是本项目本机开发环境；PostgreSQL、Redis 等依赖服务由用户通过 Laragon 手动启动。
- 执行开发、联调或浏览器验证前，先检查 `127.0.0.1:5432`、`127.0.0.1:6379` 等端口是否可用；不要反复要求用户说明服务是否已启动。
- 不要擅自更换数据库、Redis 安装路径、连接方式或项目 `.env` 配置。发现服务未监听时，报告具体端口和启动失败原因。
- 后端默认使用 `backend/.env` 的 PostgreSQL 与 Redis 配置；Redis 不得静默替换为内存实现。

## 验证与预览

- 前端开发服务器必须后台启动并记录日志、PID；向用户提供地址前先验证实际端口。
- 后端启动失败时，先读取日志确认数据库、Redis 或迁移问题，再修改代码。

## Git 交付节奏

- 阶段性功能完成并通过对应验证后再统一推送 GitHub；开发过程只保留本地提交，不频繁推送。
- 每个阶段应以可独立验收的功能边界、测试结果和提交作为里程碑。

## 模块化目录约束

- 同类功能不得长期堆放在同一个目录；当一个功能域包含多个实现文件或测试文件时，必须建立领域子目录。
- 后端 Service 统一保留在 `backend/app/services/` 根目录下，并按领域分类：`auth/`、`users/`、`rbac/`、`audit/` 等。禁止把 Service 分散到 `app/core` 或 `app/modules/<name>` 下。
- 后端 Controller、Resource、Registry、Repository 等遵循各自边界：平台基础能力放 `app/core/`，业务模块放 `app/modules/<name>/`，不要因为文件变多就重新混入 Service 根目录。
- 后端 Console 命令按领域分类；当出现多个命令域时使用 `backend/app/console/<domain>/`，例如管理命令、交易命令、数据命令分别独立。未来 Gin、交易程序或其他运行入口必须建立自己的领域目录，不得与 Admin 命令混放。
- 前端共享组件统一放 `admin/src/components/`；基础设施放 `admin/src/core/`；业务页面和业务组件放 `admin/src/modules/<name>/`。
- 新增功能必须先确定所属领域和目录，再创建文件；不得为了方便把不同功能直接放进现有大目录。
- 测试文件跟随被测领域目录放置，保持同一 Go/TypeScript package 的边界；目录迁移必须同步更新 import、文档和测试。
- 目录分类以职责和领域为依据，不按单个文件机械创建目录，也不创建没有实际边界的 `misc`、`common`、`other` 收容目录。
