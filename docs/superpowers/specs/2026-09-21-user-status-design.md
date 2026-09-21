# 用户状态模型设计

## 目标

将用户状态从布尔字段 `is_active` 升级为单一的字符串状态字段 `status`，让管理员可以明确管理用户的生命周期状态，并让认证、RBAC、资源列表和表单使用同一套状态语义。

本阶段只实现状态管理，不引入邀请注册、自动锁定策略、解锁审计或状态历史表。

## 状态定义

| 状态 | 含义 | 是否允许登录 | 管理员可选 |
| --- | --- | --- | --- |
| `active` | 正常使用 | 是 | 是 |
| `disabled` | 管理员停用 | 否 | 是 |
| `locked` | 管理员锁定 | 否 | 是 |

`status` 是唯一真实来源，不保留 `is_active` 兼容字段。`locked` 当前表示人工锁定，后续若增加登录失败自动锁定，仍复用该状态。

## 数据库迁移

新增迁移，针对已有 `users` 表执行：

1. 增加非空字符串列 `status`，默认值为 `active`。
2. 根据旧数据迁移：`is_active = true` 写入 `active`，`is_active = false` 写入 `disabled`。
3. 删除 `is_active` 列。
4. 为 `status` 增加允许值约束（若当前数据库抽象层不支持原生 CHECK，则由服务层和请求校验共同保证）。

新建数据库直接创建 `status`，不创建 `is_active`。

## 后端边界

### 模型与服务

- `models.User` 删除 `IsActive`，新增 `Status string`。
- `Public()` 只输出 `status`，不再输出 `is_active`。
- 用户服务校验状态必须是 `active`、`disabled` 或 `locked`。
- 创建用户默认 `active`；更新用户必须显式接收并校验 `status`。
- “最后一个有效管理员”规则改为统计 `status = active` 的管理员。
- 用户角色绑定和 RBAC 列表只把 `status = active` 视为可用账号。

### 认证

- 只有 `status = active` 的账号允许登录。
- `disabled` 和 `locked` 使用统一的账号不可用错误，不向未授权用户泄露内部状态差异。
- 当前用户接口返回 `status`。

### 资源接口

- 用户资源字段使用 `status`，类型为 `select`。
- 选项固定为 `active`、`disabled`、`locked`，同时提供中英文展示文案。
- 用户列表支持 `status` 查询筛选。
- 用户列表列显示状态 Badge，不再读取 `is_active`。

## 前端交互

- 用户创建表单默认选择“正常”。
- 用户编辑表单使用状态下拉框替代 Active 开关。
- 用户列表增加状态筛选，并使用语义 Badge 展示三种状态。
- `disabled` 与 `locked` 使用不同的文案和视觉语义，但不使用危险的高饱和装饰。
- 当当前用户试图将最后一个有效管理员改为 `disabled` 或 `locked`，显示后端返回的明确错误。
- 用户密码生成、显示/隐藏和复制交互保持不变。

## API 契约变化

用户创建和更新请求从：

```json
{
  "name": "Ada",
  "email": "ada@example.com",
  "password": "...",
  "locale": "zh-CN",
  "is_active": true
}
```

改为：

```json
{
  "name": "Ada",
  "email": "ada@example.com",
  "password": "...",
  "locale": "zh-CN",
  "status": "active"
}
```

响应、资源详情和列表中的用户记录都只返回 `status`。

## 测试与验收

- 迁移测试覆盖旧 `true/false` 到 `active/disabled` 的映射，并确认 `is_active` 被删除。
- 服务测试覆盖合法状态、非法状态、默认状态和最后一个 active 管理员保护。
- 认证测试覆盖 active 可登录、disabled/locked 被拒绝。
- 资源列表测试覆盖 status 筛选和返回字段。
- 前端测试覆盖状态选项、默认值、状态序列化和列表显示。
- 最终验证：`go test ./...`、前端测试、`vue-tsc -b`、`npm run build`，并在右侧预览完成创建/编辑/筛选流程。

## 非目标与风险

- 本阶段不实现登录失败次数统计、自动锁定、状态历史、审计日志和批量状态变更。
- 删除 `is_active` 是破坏性 API 变化，旧客户端必须同步升级；本项目当前只有同仓库前端客户端，因此采用一次性契约切换。
- 迁移必须先复制状态再删除旧列，并在 PostgreSQL 上验证顺序，避免现有账号状态丢失。
