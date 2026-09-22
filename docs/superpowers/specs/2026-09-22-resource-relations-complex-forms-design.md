# Resource 关系与复杂表单扩展设计

## 目标

在现有 Resource Registry、统一 CRUD API、字段权限和批量 Action 契约之上，增加可控的资源关系与基础复杂表单能力。首期只解决通用后台中高频、边界清晰的场景，不把审批、状态机、跨聚合事务或任意 Vue 代码塞入 ResourceSpec。

本阶段仍由 `admin:make-resource` 作为唯一生成入口，后端 Manifest、OpenAPI、前端页面、生成器和测试使用同一份资源描述。

## 首期范围

### 资源关系

Resource Manifest 增加声明式关系元数据：

```go
type Relation struct {
    Name         string   `json:"name"`
    Kind         string   `json:"kind"` // belongsTo | hasMany
    Resource     string   `json:"resource"`
    ForeignField string   `json:"foreign_field"`
    LabelField   string   `json:"label_field"`
    Selectable   bool     `json:"selectable"`
    Multiple     bool     `json:"multiple"`
    Permission   string   `json:"permission,omitempty"`
    FilterFields []string `json:"filter_fields,omitempty"`
}
```

- `belongsTo` 首期支持表单中的单选关系；`multiple=true` 暂不允许用于 `belongsTo`；
- `hasMany` 首期只在详情页展示关联摘要，不自动级联写入；
- `Resource`、`ForeignField`、`LabelField` 和 `FilterFields` 必须由服务端 Registry 解析，客户端不能提交任意表名或字段名；
- 关系目标资源必须存在，关系展示和选项接口同时执行目标资源权限与数据范围；
- 关系选择结果只提交 ID，后端写入前重新验证目标记录可访问且关系字段类型合法；
- 关系查询使用统一入口 `GET /api/v1/admin/{resource}/relations/{relation}/options`，返回受权限和范围过滤的 `{value,label}` 选项；
- 不自动推断 ORM 关系，不生成级联删除，不新增关系专用业务 Service。

### 表单扩展

ResourceSpec 增加有限的布局元数据：

```go
type FormGroup struct {
    Name    string   `json:"name"`
    Label   string   `json:"label"`
    Columns int      `json:"columns,omitempty"`
    Fields  []string `json:"fields"`
}

type DetailSection struct {
    Name   string   `json:"name"`
    Label  string   `json:"label"`
    Fields []string `json:"fields"`
}
```

- 未声明分组时，继续按当前字段顺序渲染默认分组；
- `Columns` 只允许受控值 `1/2/3/4`；
- 分组和详情区块只能引用当前 Manifest 已声明字段或关系，未知字段在 Registry 校验阶段拒绝；
- 基础联动只允许声明“字段 A 的值满足固定条件时显示/启用字段 B”，条件值来自服务端字段/选项契约；
- 首期不允许表达任意表达式、脚本、远程 URL、组件路径或自定义 JavaScript；
- 复杂表单由业务模块自己的 Custom Page 实现，不扩张 `ResourceFormView` 成为业务流程引擎。

## 后端契约

关系选项接口只返回当前用户可见的目标记录。目标资源的 view 权限、数据范围和字段 readable 策略全部生效；`LabelField` 不可读时关系配置无效。

创建和更新流程在持久化前执行：

1. Resource Manifest 字段 writable 校验；
2. 关系字段是否声明为可写/可选择；
3. 目标资源和关系名称是否来自服务端 Registry；
4. 目标 ID 是否存在且在当前用户权限和数据范围内；
5. 外键/多选值类型和数量限制校验；
6. 通过资源写入 Service 完成保存。

非法关系 ID 返回稳定的 `RELATION_NOT_FOUND` 或 `RELATION_FORBIDDEN`，不泄露目标记录是否存在。关系选项接口不得返回未授权目标记录。

## 前端契约

共享 `ResourceFormView` 只负责通用字段、分组、布局、关系选项加载和错误展示；资源模块仍通过 `resourceDefinition` 提供元数据和 `resourceApi` 提供 API。

- `belongsTo` 使用受控 Select；
- `hasMany` 在详情页显示只读关联区块，首期不提供通用编辑器；
- 关系选项按需加载，不把所有资源数据塞入初始 Manifest；
- 关系接口失败时表单不可提交，并显示统一错误；
- 后端返回的字段裁剪和权限错误优先于前端显示逻辑；
- 生成页面不为每个资源复制一套表单组件。

## 生成器行为

`admin:make-resource` 支持生成关系、表单分组、详情区块的元数据和 README 说明，同时生成关系选项 API 的客户端契约。生成器：

- 不自动创建关系业务 Handler、Repository 或自定义组件；
- 不自动执行迁移；
- 不覆盖人工维护的 Service、Registry、路由和业务页面；
- 对非法关系目标、重复分组字段和未知字段执行冲突预检；
- README 必须列出关系目标、权限审阅、迁移审阅、手动注册点和验证命令。

## 验收标准

- Resource Manifest 能稳定描述 `belongsTo` 和 `hasMany`；
- 关系选项只返回目标资源权限和数据范围允许的记录；
- 非法关系 ID 无法通过创建或更新写入；
- readable/writable/敏感字段策略在关系选项和表单提交中继续生效；
- 表单分组、列布局、详情区块和基础联动能被前端稳定渲染；
- 生成器 Golden、OpenAPI、TypeScript Client 和真实示例资源同步；
- 未声明关系的现有 users、roles、permissions 页面无回归；
- 不新增依赖、不执行迁移、不推送 GitHub。

## 暂不实现

- many-to-many 通用编辑器和级联保存；
- 关系批量导入、级联删除和跨资源事务；
- 任意表达式、脚本、动态组件和远程组件；
- 审批流、状态机、实时联动和复杂业务计算；
- 请求/响应完整审计与脱敏。
