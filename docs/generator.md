# CRUD Generator

Generator 的目标是减少重复代码，不是替代架构设计。生成命令在 `backend/` Goravel 项目根目录执行，前端文件写入 `admin/`。

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
