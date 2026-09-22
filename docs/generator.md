# Resource Generator

Generator 的目标是减少重复代码，不是替代架构设计。生成命令在 `backend/` Goravel 项目根目录执行，前端文件写入 `admin/`。

新增资源按模块目录生成，目录约定见 [`docs/module-layout.md`](./module-layout.md)。前端不再使用平铺 `views/` 目录；通用资源页面位于 `admin/src/core/resource/pages/`。

## 当前主流程方向

`admin:make-resource` 是唯一主入口，由同一份 ResourceSpec 统一生成后端、权限、菜单、前端 ResourceSpec、路由描述、迁移、测试和 README。普通资源默认使用 `Manifest + 通用列表/表单/详情页面`；只有显式声明 `pageMode: custom` 的复杂资源才生成并接入专用页面。Permission/Menu 只作为内部渲染器存在，不提供独立命令；`admin:make-crud` 已移除，避免形成第二套生成流程。

当前生成器采用自动发现：生成代码后更新生成专属 discovery 文件，资源自动进入后台 Registry、router 和 Sidebar；不修改人工维护文件，不写权限/角色数据，不执行迁移。

## `admin:make-resource`

```text
go run . artisan admin:make-resource posts \
  --label="Posts" \
  --route="/admin/posts" \
  --permission="admin.posts.view" \
  --field=title:text:required \
  --field=published:boolean \
  --field=status:select:required:active=Active|disabled=Disabled
```

字段格式为 `name:type[:required[:value=Label|value=Label]]`；选项仅适用于 `select`，例如 `status:select:required:active=Active|disabled=Disabled`。命令会生成 Resource、Model、Request、Repository、Service、Controller、Routes、Permissions、菜单、ResourceSpec、统一 API 契约、测试、README 和 Migration 文件，并更新生成专属 discovery 文件。普通资源不会生成资源专用 ListPage/FormPage/DetailPage，而是自动进入核心通用页面；`pageMode: custom` 才会生成专用页面覆盖。通用表单根据 Manifest 字段、关系、分组和依赖渲染，列表页根据字段类型自动提供搜索、筛选、排序、详情、编辑、删除和批量操作；API 契约包含 list/show/create/update/delete 五类操作；迁移文件只是待审阅的代码产物，必须人工确认后再执行；命令本身不会连接数据库或运行迁移。

字段可追加权限修饰符：`sensitive` 表示默认不搜索和导出，`readonly` 表示不可写，`hidden` 表示不出现在公开资源页面；例如 `password:text:sensitive`、`owner_id:integer:readonly`。多个修饰符可组合，生成器会把完整的 `visible/readable/writable/sensitive` 策略同步写入后端 Manifest 和前端元数据。

需要限制为本人数据的资源可使用 `--scope=own --owner-field=owner_id`，其中 `owner_id:integer` 必须同时出现在字段定义中。生成后在 RBAC 中为角色把对应权限设置为“仅本人数据”；后端创建时以当前登录用户为 owner，并在列表、详情、搜索、导出、更新和删除时强制应用范围。

关系和基础布局可通过受控元数据生成：

```text
--relation=customer:belongsTo:customers:customer_id:id:name:selectable
--form-group=main:Main:2:customer_id|status
--detail-section=summary:Summary:status
```

关系格式为 `name:kind:target_resource:field:foreign_field:label_field[:selectable]`。关系选项由后端统一接口提供，并再次执行目标资源权限、数据范围和字段可读校验；前端不能提交任意表名或字段名。`hasMany` 首期只生成详情只读区块，不自动级联写入。

生成文件写入 `backend/app/modules/<name>/`、`admin/src/modules/<name>/` 和 `backend/database/migrations/`。如果任一目标文件已经存在，命令会在写入前失败，不会覆盖人工文件，也不会留下半套输出。

生成完成后的人工审阅顺序：

1. 审阅生成的接口骨架、字段和迁移；
2. 检查生成专属 discovery 文件是否包含新资源；普通资源会通过该文件自动进入 Registry、通用路由和 Admin Shell；
3. 手动审阅并执行数据库迁移；
4. 按业务需要补充非通用动作和授权规则。

## 模块检查

模块检查只读确认资源生成结果是否完整，不创建第二套模块生成流程：

```text
go run . artisan admin:check-module billing
```

检查命令只读取文件，不创建、修改或删除文件，也不连接数据库；除资源骨架外，还会检查后端 Registry 和前端 router/Shell 使用的生成专属 discovery 文件。检查通过后资源已经由 discovery 文件自动进入运行时。

## 当前命令

```text
admin:make-resource
admin:api
admin:check-module
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
Resource Manifest
Permission
Generic Page Registration
Tests
```

相同输入应尽量产生相同输出；生成代码与人工代码分离；重跑不能覆盖人工 Service、Action、Page 和 Component；没有 AI 时生成器仍必须可用。

AI 只能提供字段标签、关系、筛选和布局建议，不能替复杂业务决定流程、状态机、交易规则、领域事件语义或安全策略。

每个生成器必须有 Golden File/Snapshot 测试，验证输出稳定且重跑不会破坏自定义文件。
