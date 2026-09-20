# Resource Engine

Resource Engine 是标准 CRUD 加速器，不是强制所有业务使用的低代码平台。

## 三种使用等级

```text
Simple Resource    标准列表、表单、详情和基础操作
Extended Resource  标准资源 + 自定义 Action、Slot 或详情区
Custom Page        普通 Vue 页面 + Goravel Controller/Service
```

## Resource Definition

```text
name
label
route
permissions
fields
columns
filters
actions
navigation
```

标准能力：分页、搜索、排序、筛选、新增、编辑、查看、删除、行操作、批量操作、字段权限和审计钩子。

多步骤流程、审批和状态机、实时数据、图表分析、复杂联动表单和高度定制详情页直接使用 Custom Page。

扩展优先使用明确的 Go 接口和 Vue slot/component，而不是无限增加 JSON 字段。禁止把任意可执行代码序列化进 Resource Manifest。

```text
模块注册 -> Resource 注册 -> 权限/导航注册 -> API 暴露 -> Vue 页面消费
```
