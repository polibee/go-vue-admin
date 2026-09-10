# Resource Manifest 与代码生成

Resource Manifest 是数据库结构和后台业务配置之间的稳定边界。生成器可以从数据库推断字段，但生成结果必须落成 Manifest，之后的代码生成只读取 Manifest，不直接依赖数据库连接。

## 最小格式

```yaml
id: products
table: products
label: 产品
primary_key: id
soft_delete: false
audit: false
permissions:
  list: products.view
  get: products.view
  create: products.create
  update: products.update
  delete: products.delete
fields:
  - name: name
    type: text
    label: 名称
    required: true
    searchable: true
    sortable: true
  - name: status
    type: select
    label: 状态
    options:
      - value: draft
        label: 草稿
      - value: active
        label: 启用
```

正式结构见 `contracts/schemas/resource-manifest.json`。

## 生成命令

```bash
go -C backend run ./cmd/admin-gen module catalog
`admin-gen module` and the initial MySQL-backed `resource` discovery stage are implemented. The resource command consumes the Manifest as its only intermediate representation:

```bash
go -C backend run ./cmd/admin-gen resource --module catalog --table products
```
```

生成结果位于：

```text
modules/catalog/
├── resources/<resource>.yaml
├── backend/resources/<resource>.go
└── admin/resources/<resource>.ts
```

资源后端文件包含 `RegisterRoutes(auth)`，它使用通用 GORM Repository 和标准 Core Resource Controller，并由生成器同步写入模块入口、应用路由和 `backend/go.mod` replace。生成器只允许 Manifest 声明的字段进入 SQL 写入和排序白名单。运行时由 `DB_CONNECTION=mysql|postgres` 选择 GORM 方言，生成代码不包含平台专用数据库实现。

数据库反向读取、关系推断和完整 CRUD 生成以 Manifest 为中间结构；当前已支持通用 GORM 单表 CRUD，Manifest 是后续枚举、外键、软删除和审计能力的扩展点。

## 数据库选择

`RESOURCE_PROVIDER=database` 启用持久化资源。`DB_CONNECTION` 是唯一的数据库方言选择项，支持 `mysql` 和 `postgres`（`pgsql`、`postgresql` 作为兼容别名）。Goravel 的迁移连接与资源模块独立持有的 GORM 连接池读取同一组连接环境变量，并在应用关闭时统一释放资源连接池。

数据库切换只需要修改配置，不需要修改生成的模块代码：

```env
RESOURCE_PROVIDER=database
DB_CONNECTION=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_DATABASE=go_vue_admin
DB_USERNAME=postgres
DB_PASSWORD=
```

## 约束

- 生成器只写入目标 Module，不写入 `backend/app/core`。
- 数据库字段只能提供默认推断，业务标签、权限和选项允许人工覆盖。
- 生成后的文件应提交到仓库，运行时不扫描数据库生成 Vue 页面。
- 复杂查询和自定义动作放入模块扩展代码，不修改生成器模板。
