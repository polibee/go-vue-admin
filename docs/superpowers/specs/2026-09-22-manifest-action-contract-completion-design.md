# Manifest Action 完整契约补齐设计

## 目标

补齐总 Resource Capabilities 设计中尚未落地的批量 Action 元数据和结果契约，使 Resource Manifest 能明确描述 Action 是否允许批量执行、使用哪种 payload 契约，并让前端和生成器能够安全消费同一份定义。

本阶段只扩展现有批量 Action，不新增任意字段更新、导入、脚本或跨资源事务能力。

## 核心契约

```go
type Action struct {
    Name       string `json:"name"`
    Label      string `json:"label"`
    Kind       string `json:"kind"`
    Permission string `json:"permission"`
    Batch      bool   `json:"batch"`
    Payload    string `json:"payload,omitempty"`
}
```

- `Batch=false` 的 Action 只能用于行级操作，不能进入批量入口；
- `Batch=true` 才能调用统一批量 API；
- `Payload` 是服务端已注册的契约标识，不是客户端可执行代码；
- `Kind` 决定 Handler，`Payload` 决定输入校验契约；二者必须由服务端白名单解析；
- `set-status` 使用 `Kind=user-status`、`Batch=true`、`Payload=user-status`。

请求体统一使用 `payload`，不再使用 `params`：

```json
{
  "ids": [1, 2, 3],
  "payload": {"status": "disabled"}
}
```

响应补充 skipped 结果：

```json
{
  "data": {
    "action": "set-status",
    "requested": 3,
    "succeeded": 1,
    "failed": 1,
    "skipped": 1,
    "failures": [{"id": 9, "code": "RESOURCE_NOT_FOUND"}],
    "skips": [{"id": 10, "code": "OUT_OF_SCOPE"}]
  }
}
```

无权限或范围外记录使用 skipped，不暴露记录存在性；参数错误、领域校验失败使用 failures。Handler 自己负责领域事务，Manifest 不声明伪造的跨资源事务策略。

## 生成器和前端

`admin:make-resource` 的 ResourceSpec 增加可选 Action 定义，生成稳定的 `batch`、`kind`、`payload` 元数据；没有 Action 定义时标准 CRUD Action 默认 `batch=false`。生成器不生成任意业务 Handler，只输出 README 中的 Handler 注册清单和契约说明。

前端只展示 `batch=true` 且权限通过的 Action；`kind` 对应参数表单，未知 Kind 显示不可执行状态而不是猜测字段。OpenAPI、TypeScript Client、Golden File 和真实示例资源同步更新。

## 验收标准

- Manifest 能区分行级 Action 和批量 Action；
- 缺失或未知 Payload 契约不能执行；
- `set-status` 使用 `payload` 并返回 succeeded/failed/skipped；
- 越权记录进入 skipped，不泄露存在性；
- 生成器、OpenAPI、前端和 README 输出一致；
- 旧 `params` 请求不作为兼容入口保留；
- 不执行迁移、不新增依赖、不推送 GitHub。
