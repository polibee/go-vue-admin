# Go Vue Admin

基于 Goravel、Vue 3、TypeScript 和 shadcn-vue 的模块化通用后台管理平台。系统以 Resource Manifest 驱动资源列表、表单、详情、权限、菜单、搜索和批量操作。

[![RackNerd VPS](https://img.shields.io/badge/RackNerd-VPS-2563eb?style=for-the-badge)](https://my.racknerd.com/aff.php?aff=7572)
[![Vast.ai GPU Cloud](https://img.shields.io/badge/Vast.ai-GPU%20Cloud-7c3aed?style=for-the-badge)](https://cloud.vast.ai/?ref_id=91181)

## 主要能力

- JWT 登录、刷新、退出，以及 PostgreSQL 权威 Refresh Token 存储；Redis 作为加速层。
- 用户、角色、权限、数据范围和字段权限。
- Resource Registry 驱动的列表、创建、编辑、详情和删除。
- 搜索、筛选、排序、分页、CSV 导出。
- 当前页、跨页、按筛选条件全选，以及批量删除、批量更新和 Manifest Action。
- 资源关系、分组表单和字段依赖扩展点。
- 全局搜索、菜单分组、权限过滤。
- 请求/响应审计、敏感字段脱敏、手动清理和定期清理。
- admin:make-resource 生成后端 Resource、迁移、权限、菜单、前端 ResourceSpec、测试和 README；普通资源默认复用通用页面，复杂资源可显式覆盖。

## 本地运行

项目依赖由 Laragon 管理，启动 PostgreSQL 和 Redis 后执行：

~~~powershell
cd backend
go run .
~~~

另开终端启动管理面板：

~~~powershell
cd admin
pnpm install
pnpm dev
~~~

端口约定：

- Go API：http://127.0.0.1:3000
- Go/Vue 管理面板：http://127.0.0.1:5180
- Python 管理面板（独立项目）：http://127.0.0.1:5173

数据库迁移必须先审阅生成文件，再按项目流程手动执行；生成器不会静默执行迁移。

## 资源生成

~~~powershell
cd backend
go run . admin:make-resource announcements --fields="title:text:required,status:select:required:draft=Draft|published=Published"
~~~

资源生成是面向业务开发的主流程：生成后由后台面板自动发现并展示。生成器不会覆盖人工文件；复杂业务应放入对应的 backend/app/modules/<module> 和 admin/src/modules/<module>，不要强行套用通用 Resource。

## 验证

~~~powershell
cd backend
go test ./...

cd ..\admin
pnpm exec vue-tsc --noEmit
pnpm run build
~~~

## 目录约定

- 后端基础能力：backend/app/core
- 后端业务模块：backend/app/modules/<module>
- 后端 Service：backend/app/services/<domain>
- 后端命令：backend/app/console/<domain>
- 前端共享组件：admin/src/components
- 前端基础设施：admin/src/core
- 前端业务页面和组件：admin/src/modules/<module>

生产部署前仍需按实际环境复核密钥、数据库、Redis、反向代理、日志保留策略、备份、监控和迁移流程。本项目 README 描述开发基线，不替代环境验收。
