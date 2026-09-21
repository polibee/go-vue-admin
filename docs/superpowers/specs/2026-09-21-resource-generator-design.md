# Resource-Driven Generator Design

## 目标

Generator 的唯一主入口是资源，而不是彼此独立的代码片段。`admin:make-resource` 使用一份资源定义，统一派生后端 Resource、CRUD 边界、权限、菜单、前端列表/表单/详情页面、路由描述、迁移和测试骨架。

生成器负责产出完整、可审阅的代码集合；它不把新功能悄悄接入现有运行时。生成完成后，开发者根据 README 和只读检查结果，手动注册 Provider、路由、权限、菜单、Resource Registry，并手动审阅和执行迁移。

## 主入口

```text
go run . artisan admin:make-resource posts \
  --label="Posts" \
  --route="/posts" \
  --permission="admin.posts.view" \
  --icon="file-text" \
  --field=title:text:required \
  --field=published:boolean
```

不保留 `admin:make-crud`。资源生成已经包含完整 CRUD 所需的边界，避免用户在两个入口之间选择，也避免维护第二套解析和渲染流程。

权限和菜单只保留为 Resource Generator 内部渲染器，不再提供独立 Artisan 命令；它们不得产生与资源流水线不一致的命名和目录结构。

## 输入定义

资源定义包含：

- `name`：小写字母、数字和连字符组成，且以字母开头；
- `label`：管理端显示名称；
- `route`：资源前端路由；
- `permission`：资源查看权限；
- `icon`：菜单图标标识；
- 有序字段列表：`name:type[:required]`；
- 可选列、筛选、动作和关系元数据，后续按兼容字段扩展。

第一阶段字段类型为 `text`、`email`、`password`、`integer`、`boolean` 和 `select`。字段定义只解析一次，所有后端、前端、权限、菜单和测试产物都从规范化后的 Spec 派生。

## 生成产物

一次资源生成至少包含以下边界：

```text
backend/app/generated/resources/<name>/manifest.go
backend/app/generated/resources/<name>/model.go
backend/app/generated/resources/<name>/request.go
backend/app/generated/resources/<name>/repository.go
backend/app/generated/resources/<name>/service.go
backend/app/generated/resources/<name>/controller.go
backend/app/generated/resources/<name>/routes.go
backend/app/generated/resources/<name>/permissions.go
backend/app/generated/resources/<name>/manifest_test.go
backend/app/generated/resources/<name>/README.md
backend/database/migrations/<timestamp>_create_<name>_table.go

admin/src/generated/resources/<name>/resource.ts
admin/src/generated/resources/<name>/api.ts
admin/src/generated/resources/<name>/menu.ts
admin/src/generated/resources/<name>/routes.ts
admin/src/generated/resources/<name>/pages/<Name>ListPage.vue
admin/src/generated/resources/<name>/pages/<Name>FormPage.vue
admin/src/generated/resources/<name>/pages/<Name>DetailPage.vue
admin/src/generated/resources/<name>/<name>.test.ts
```

页面产物必须复用现有 Admin Shell、ResourceList、Resource Form、Resource Detail 和权限语义，不复制一套新的 UI 基础组件。生成的 List/Form/Detail 页面可以是薄包装器，负责传入资源元数据和权限；复杂业务页面才允许人工扩展。

权限和菜单文件是资源生成结果的一部分：

- 权限常量由 `admin.<name>.<action>` 派生；
- 菜单配置包含 label、route、icon 和 permission；
- 路由描述包含列表、新建、编辑和详情入口；
- API 契约包含 list、show、create、update、delete 五类操作，并明确对应的 `/api/v1/admin/<name>` 路径；
- 所有这些文件只写入生成目录，不自动修改现有 Registry、Sidebar 或 router 文件。

## 安全边界

- 不自动注册 Provider、路由、权限、菜单或 Resource Registry。
- 不自动执行迁移，不连接数据库，不修改现有迁移或 Seeder。
- 允许生成新的迁移文件，但迁移必须人工审阅后执行。
- 所有目标文件在写入前统一预检；任一冲突都导致整次生成失败，不能留下半套输出。
- 默认不覆盖任何人工文件；不提供隐式强制覆盖。
- 生成文件包含对应生成器标记和 README 接入步骤。
- 不引入第三方依赖，不做数据库反向推断，不调用 AI 生成业务规则。

## 内部架构

流水线固定为：

```text
Input
  -> ResourceSpec
  -> BackendRenderer
  -> PermissionRenderer
  -> MenuRenderer
  -> FrontendRenderer
  -> TestAndReadmeRenderer
  -> ConflictAwareWriter
```

每个 Renderer 只消费同一份 `ResourceSpec` 并返回内存中的 Artifact；Writer 最后统一处理冲突和写入。权限和菜单 Renderer 不提供独立命令，不得自行解析另一套输入。

## 测试与验收

- Spec 测试验证字段、权限、菜单、路由和页面路径的一致性。
- Golden File 测试覆盖一份完整 `posts` 资源的后端、权限、菜单、前端页面、README 和迁移产物。
- 冲突测试证明任一目标文件冲突时所有目录都不写入。
- 真实 smoke 测试执行 `admin:make-resource`，确认生成全链路文件但不执行迁移、不修改注册文件。
- 生成产物必须通过后端 `go test ./... -count=1` 和前端 `vue-tsc`、构建及测试。
- 表单页面统一复用共享 `ResourceFormView`，不复制用户或角色专用表单实现。

## 后续边界

后续可以增加字段关系、筛选器、批量动作和可选页面模板，但必须继续从同一份 ResourceSpec 派生。插件系统、数据库反向生成、AI 业务规则和自动运行时注册需要另行设计，不属于本阶段。
