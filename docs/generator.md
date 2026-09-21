# CRUD Generator

Generator 的目标是减少重复代码，不是替代架构设计。生成命令在 `backend/` Goravel 项目根目录执行，前端文件写入 `admin/`。

## 当前主流程方向

现有 Resource、Module、Permission 和 Menu 命令已经提供可复用的基础渲染能力，但它们不是最终的业务入口。下一阶段以 `admin:make-resource` 和 `admin:make-crud` 为主入口，由同一份资源定义统一生成后端、权限、菜单、前端列表/表单/详情页面、路由描述、迁移和测试产物。独立 Permission/Menu 命令只保留为兼容性适配器。

当前生成器仍遵守安全边界：生成代码但不自动注册运行时、不修改现有路由或 Sidebar、不写权限/角色数据、不执行迁移。

## 已实现：`admin:make-resource`

当前第一阶段只生成后端 Resource 基础骨架，不自动接入路由、Resource Registry 或数据库迁移执行器。

```text
go run . artisan admin:make-resource posts \
  --label="Posts" \
  --route="/admin/posts" \
  --permission="admin.posts.view" \
  --field=title:text:required \
  --field=published:boolean
```

命令会生成 Resource、Model、Request、Repository、Service、Controller、Routes、Permissions、Manifest Test、README 和 Migration 文件。README 会列出生成文件、Registry/Routes/权限的手动接入步骤以及 `go test ./...` 验证命令。迁移文件只是待审阅的代码产物，必须人工确认后再执行；命令本身不会连接数据库或运行迁移。

生成文件写入 `backend/app/generated/resources/<name>/` 和 `backend/database/migrations/`。如果任一目标文件已经存在，命令会在写入前失败，不会覆盖人工文件，也不会留下半套输出。

生成完成后的人工接入顺序：

1. 审阅生成的接口骨架、字段和迁移；
2. 手动将 Manifest 注册到 Resource Registry；
3. 手动注册需要暴露的 Routes 和权限；
4. 手动执行已审阅的数据库迁移。

## 已实现：`admin:make-module`

模块生成器只创建编译期模块边界，不改变应用启动注册表：

```text
go run . artisan admin:make-module billing
```

命令会在 `backend/app/modules/billing/` 下生成 `module.go`、`model.go`、`request.go`、`repository.go`、`service.go`、`controller.go`、`routes.go`、`resource.go`、`permissions.go`、`events.go` 和 `tests/module_test.go`。它不生成迁移，也不自动注册 Provider、Routes、Permissions、菜单或 Resource Registry。

模块目录同时包含 `README.md`，列出人工接入顺序和检查命令。可以使用只读检查命令确认生成文件是否完整：

```text
go run . artisan admin:check-module billing
```

检查命令只读取文件，不创建、修改或删除文件，也不连接数据库。检查通过后仍需人工注册运行时边界；当前版本没有默认开启的自动注册选项。

## 已实现：`admin:make-crud`

CRUD 生成器组合模块和 Resource 生成器，生成完整的后端边界：

```text
go run . artisan admin:make-crud orders \
  --label="Orders" \
  --route="/admin/orders" \
  --permission="admin.orders.view" \
  --field=number:text:required \
  --field=paid:boolean
```

它会一次性生成模块骨架、Resource 骨架、两个 README、测试骨架和迁移文件。所有目标文件会统一预检，任一文件冲突时不会写入任何文件。生成结果仍需人工注册运行时边界，迁移仍需人工审阅和执行。

## 已实现：`admin:make-permission`

权限生成器只生成权限常量和说明，不写入权限表或角色绑定：

```text
go run . artisan admin:make-permission orders \
  --action=view \
  --action=create \
  --action=update \
  --action=delete
```

它会生成权限常量、README 和 Golden File 测试所需的稳定输出。权限数据和角色分配仍需通过人工审阅后的 RBAC 流程完成。

## 已实现：`admin:make-menu`

菜单生成器只生成菜单配置骨架和接入说明：

```text
go run . artisan admin:make-menu orders \
  --label="Orders" \
  --route="/orders" \
  --permission="admin.orders.view" \
  --icon="shopping-cart"
```

它会生成菜单常量、README 和 Golden File 稳定输出。菜单与权限会被记录在生成配置中，但不会自动修改前端路由、Sidebar、数据库或现有菜单注册。

## 第一阶段命令

```text
admin:make-module
admin:make-resource
admin:make-crud
admin:make-permission
admin:make-menu
admin:api
```

## 生成范围

```text
Module
Model
Migration
Request
Repository
Service
Controller
Routes
Resource
Permission
Vue Pages
Tests
```

相同输入应尽量产生相同输出；生成代码与人工代码分离；重跑不能覆盖人工 Service、Action、Page 和 Component；没有 AI 时生成器仍必须可用。

AI 只能提供字段标签、关系、筛选和布局建议，不能替复杂业务决定流程、状态机、交易规则、领域事件语义或安全策略。

每个生成器必须有 Golden File/Snapshot 测试，验证输出稳定且重跑不会破坏自定义文件。
