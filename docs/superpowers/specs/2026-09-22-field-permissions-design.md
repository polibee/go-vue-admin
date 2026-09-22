# Resource 字段权限设计

## 目标

在现有 Resource、RBAC 和统一 CRUD 契约上增加字段级访问控制，使后台资源能够安全地控制字段展示、读取、写入、搜索和导出。后端负责最终裁剪和校验，前端只负责根据契约改善操作体验。

本阶段不建立第二套表单系统，不把字段权限散落到 Controller 或业务 Service，也不实现请求/响应完整审计脱敏。

## 范围与边界

### 本阶段包含

- Resource Manifest 字段策略：`visible`、`readable`、`writable`、`sensitive`；
- 统一字段策略解析服务；
- 列表、详情、导出响应的字段裁剪；
- 搜索、排序和导出对不可读/敏感字段的保护；
- 创建、更新对不可写字段的拒绝；
- 角色权限字段覆盖配置；
- OpenAPI、生成器、前端资源页面、测试和 README 同步；
- 既有资源未声明策略时保持现有行为。

### 本阶段不包含

- 请求/响应完整内容审计与脱敏日志；
- 字段级数据加密或密钥管理；
- 任意运行时脚本或动态 SQL 字段权限；
- 复杂关系字段、联动表单和批量 Action；
- 自动向角色授予字段权限；
- 自动执行数据库迁移。

## 核心契约

`resource.Field` 增加以下元数据：

```go
type Field struct {
    Name      string   `json:"name"`
    Label     string   `json:"label"`
    Type      string   `json:"type"`
    Required  bool     `json:"required,omitempty"`
    Options   []Option `json:"options,omitempty"`
    Visible   bool     `json:"visible"`
    Readable  bool     `json:"readable"`
    Writable  bool     `json:"writable"`
    Sensitive bool     `json:"sensitive"`
}
```

为保持现有 Manifest 兼容，未显式声明的策略按以下默认值处理：

```text
visible=true
readable=true
writable=true
sensitive=false
```

`visible=false` 表示字段不进入公开资源契约和前端页面；`readable=false` 表示后端响应中不得出现该字段；`writable=false` 表示客户端提交该字段时返回稳定的字段权限错误；`sensitive=true` 表示字段即使可读，也默认不参与搜索、CSV 导出和普通审计 metadata。

字段名、表名、排序字段和查询字段仍然只能来自已校验的 Manifest 元数据，不能由客户端拼接。

## 角色字段覆盖

Manifest 定义资源的安全上限，角色权限只能进一步收紧，不能重新开放 Manifest 已禁止的字段。

建议增加 `permission_role_field` 关系表：

```text
role_id
permission_id
field_name
readable
writable
created_at
updated_at
```

未配置覆盖时继承 Manifest 默认策略。用户拥有多个角色时，字段权限采用“任一角色明确允许即可允许”的合并规则，但最终结果仍受 Manifest 上限约束。这样与现有数据范围权限的 `all` 优先规则一致，并避免误配置一个角色导致用户失去另一个角色明确授予的能力。

字段覆盖只针对已存在于对应 Resource Manifest 的字段；未知字段、未知权限和非法组合在保存时返回 `422 VALIDATION_ERROR`。

## 后端执行边界

统一的字段策略服务在资源 Controller 与 Repository/ORM 查询之间提供以下能力：

```text
Resolve(user, manifest, action) -> effective field policy
SelectReadable(record, policy) -> response record
ValidateReadableQuery(field, policy) -> allow/deny
ValidateWritablePayload(payload, policy) -> filtered payload or field error
```

具体规则：

- 列表和详情只返回 `visible && readable` 字段；
- 导出使用同一套响应裁剪，并排除 `sensitive` 字段；
- 搜索字段必须同时满足 `visible && readable && !sensitive`；
- 排序字段必须满足 `visible && readable && !sensitive`，否则使用安全默认排序或返回稳定校验错误；
- 创建和更新只接受 `writable` 字段；提交不可写字段不得静默成功，返回包含字段名的稳定权限错误；
- 客户端提交只读字段不能通过忽略、空值或额外 JSON 属性绕过；
- 数据范围过滤先于字段裁剪，字段权限不能扩大数据范围；
- 用户、角色等专用资源也必须经过统一字段策略边界，保留其领域 Service 的业务校验。

后端裁剪发生在序列化响应之前，不能只依赖 Vue 隐藏列或表单项。

## 前端行为

前端继续使用现有 Resource 列表、表单和详情页面：

- `visible=false` 字段不渲染；
- `writable=false` 字段不渲染输入控件；
- `readable=false` 字段不渲染详情和列表列；
- 后端返回的 Resource Manifest 成为页面元数据来源；
- 前端不自行推断敏感字段，也不复制一套权限判断；
- 后端返回字段权限错误时显示稳定的字段级提示。

## 生成器行为

`admin:make-resource` 继续是唯一公共入口。字段定义增加显式策略选项，示例：

```text
--field password:string:sensitive
--field owner_id:integer:readonly
```

生成器必须：

- 校验字段策略组合；
- 在后端 Manifest 和前端资源元数据中生成相同策略；
- 在 README 中说明字段读取和写入行为；
- 生成迁移但不执行迁移；
- 更新生成专属 discovery 文件；
- 不自动注册角色字段覆盖；
- 保持默认资源输出确定性，未声明策略时输出兼容默认值。

## API 与错误契约

角色字段覆盖沿用现有角色权限管理入口，不新增第二套 RBAC 页面。请求至少包含：

```json
{
  "permission_ids": [10],
  "scopes": {"10": "all"},
  "fields": {
    "10": {
      "email": {"readable": true, "writable": false}
    }
  }
}
```

非法字段策略、未知字段、提交不可写字段使用现有错误格式，并保持 `422 VALIDATION_ERROR`；资源不存在或无权访问继续使用现有资源错误语义，不通过错误内容泄露字段存在性。

OpenAPI、生成客户端、RBAC 页面和 Resource 页面必须同步更新，不能手工维护互相冲突的字段权限类型。

## 验收标准

1. Manifest 未声明策略的现有资源行为不回归。
2. `visible=false` 字段不出现在资源契约、列表、详情和表单中。
3. `readable=false` 字段不会因列表、详情、搜索、排序、导出或直接 API 请求泄露。
4. `writable=false` 字段在创建和更新时均被拒绝，客户端不能用额外 JSON 属性绕过。
5. `sensitive=true` 字段默认不参与搜索和 CSV 导出。
6. 多角色字段权限按“任一角色允许即可，但不能突破 Manifest 上限”合并。
7. 字段策略与 `all/own` 数据范围同时存在时，先完成数据范围过滤，再完成字段裁剪。
8. 生成器为默认字段和受限字段分别生成稳定的后端、前端、OpenAPI、README 和测试产物。
9. 后端单元/边界测试、OpenAPI 测试、前端类型检查、Resource 页面测试和 Golden File 测试通过。
10. 迁移经过人工审阅后才能执行；生成器和测试不得自动执行迁移。

## 分阶段实现顺序

1. Resource 字段策略、默认值和后端裁剪服务；
2. 创建/更新字段保护与查询字段保护；
3. `permission_role_field` 持久化、RBAC API 和角色页面配置；
4. 生成器、OpenAPI、前端页面和真实示例资源验收。
