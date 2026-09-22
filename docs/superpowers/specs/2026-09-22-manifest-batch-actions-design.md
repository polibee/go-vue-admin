# Manifest 驱动批量 Action 设计

## 目标

把资源批量操作统一纳入 Resource Manifest。资源声明 Action 后，后端提供统一执行入口，前端根据 Manifest 和当前权限自动展示批量操作，不再为 users、roles 或单个业务资源各自维护批量接口和页面分支。

本阶段以用户状态批量修改作为第一条真实 Action，验证通用契约；后续资源可以复用同一入口增加自己的受控 Handler。

## 范围与边界

### 本阶段包含

- 统一批量 Action 请求/响应契约；
- Manifest Action 的名称、标签、类型和权限校验；
- `POST /admin/resources/{resource}/actions/{action}` 执行入口；
- ID 数量限制、重复 ID 归一化、数据范围校验和 Action 权限校验；
- 受控 Handler 注册表，不允许客户端提交 SQL、表名、字段名或脚本；
- `set-status` Handler 迁移现有用户批量状态逻辑；
- 统一成功数量、失败数量和失败明细响应；
- 前端资源列表按 Manifest 自动显示批量 Action；
- OpenAPI、类型化 Client、测试、审计和文档同步。

### 本阶段不包含

- 任意字段批量更新；
- 任意资源删除或批量导入；
- 客户端定义 Action、SQL、脚本或 Handler；
- 跨资源批量事务；
- 请求/响应完整内容审计与脱敏；
- 自动执行数据库迁移。

## 核心契约

Manifest 继续使用现有 Action：

```go
type Action struct {
    Name       string `json:"name"`
    Label      string `json:"label"`
    Kind       string `json:"kind"`
    Permission string `json:"permission"`
}
```

批量入口：

```text
POST /admin/resources/{resource}/actions/{action}
```

请求体：

```json
{
  "ids": [1, 2, 3],
  "params": {
    "status": "disabled"
  }
}
```

约束：

- `ids` 必须非空，服务端去重并限制最大数量；
- Action 必须存在于目标 Manifest；
- 当前用户必须拥有 Action 声明的权限；
- 每个 ID 必须通过该 Action 的数据范围校验；
- `params` 只由对应 Handler 解析，未知参数拒绝；
- 返回结果不得泄露无权限记录的存在性。

响应：

```json
{
  "data": {
    "action": "set-status",
    "requested": 3,
    "succeeded": 2,
    "failed": 1,
    "failures": [
      {"id": 9, "code": "RESOURCE_NOT_FOUND"}
    ]
  }
}
```

失败明细只返回稳定错误码和 ID，不返回密码、Token 或其他敏感数据。

## 后端设计

控制器只负责解析请求、查找 Manifest、校验公共边界和调用 Handler。Handler 由内部注册表按 `Kind` 解析，负责领域校验、数据范围、字段权限和审计。资源表名、字段名和 SQL 仍只能来自服务端 Manifest 或固定 Handler，不能由请求参数决定。

`set-status` 仅允许 Manifest 已声明的状态字段和值，并复用现有用户状态服务/审计语义。旧的用户专用批量状态入口在实现阶段移除或改为内部调用统一执行器，不再作为前端或公开架构的第二套流程。

Action 的数据范围按逐 ID 检查执行；如果部分记录无权访问，允许其他合法记录继续执行，并在结果中返回稳定失败明细。执行前后均不扩大字段权限，批量 Action 不绕过字段权限模型。

## 前端设计

资源列表维护统一选中 ID 集合。当前 Manifest 存在可执行的批量 Action 且当前用户拥有其权限时，显示 Action 菜单；表单参数由 Action `kind` 决定，第一阶段只提供状态选择器。提交后显示汇总结果、刷新列表并清理选择。

前端不根据资源名称判断是否显示批量按钮。`users` 的 `set-status` 只是第一个 Manifest Action 示例，新增资源只需声明 Action 并提供受控 Handler。

## 验收标准

- 无 Action 声明的资源无法通过统一入口执行；
- 无权限用户得到稳定的禁止响应；
- 超过数量上限、空 ID 或未知参数被拒绝；
- `set-status` 可对用户多选执行，并正确处理成功与失败明细；
- own 数据范围不会修改其他用户记录；
- 前端只展示 Manifest 声明且用户有权限的批量 Action；
- 后端、OpenAPI、Client、前端类型检查和生产构建通过；
- 不执行数据库迁移，不新增依赖，不推送 GitHub。
