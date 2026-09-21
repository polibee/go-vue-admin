# CRUD Generator

Generator 的目标是减少重复代码，不是替代架构设计。生成命令在 `backend/` Goravel 项目根目录执行，前端文件写入 `admin/`。

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

命令会生成 Resource、Model、Request、Repository、Service、Controller、Routes、Permissions、Manifest Test 和 Migration 文件。迁移文件只是待审阅的代码产物，必须人工确认后再执行；命令本身不会连接数据库或运行迁移。

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
