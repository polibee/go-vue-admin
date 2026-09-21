# Module Layout

项目按业务模块组织新增代码。共享基础设施和模块业务代码分离，生成器只向模块目录写入新产物。

## Backend

```text
backend/app/
├── core/
│   ├── resource/                  # Resource Engine，只放通用类型和 Registry
│   ├── admin/controllers/         # 管理端基础 HTTP 入口
│   └── auth/controllers/          # 认证基础 HTTP 入口
└── modules/
    ├── admin/registry/            # 已启用的管理端业务资源注册表
    └── <resource>/                # 生成器输出的业务模块
        ├── resource/
        │   ├── manifest.go
        │   ├── model.go
        │   ├── request.go
        │   ├── repository.go
        │   ├── service.go
        │   ├── controller.go
        │   ├── routes.go
        │   ├── permissions.go
        │   ├── manifest_test.go
        │   └── README.md
        ├── permissions/
        │   ├── permissions.go
        │   └── README.md
        └── menu/
            ├── menu.go
            └── README.md
```

迁移仍统一放在 `backend/database/migrations/`，因为迁移属于数据库边界，不属于单个 Go 包。运行时 Registry、Provider、Routes、权限和菜单不会自动修改。

## Frontend

```text
admin/src/modules/<resource>/
├── api.ts
├── resource.ts
├── menu.ts
├── routes.ts
├── pages/
│   ├── <Resource>ListPage.vue
│   ├── <Resource>FormPage.vue
│   └── <Resource>DetailPage.vue
└── <resource>.test.ts
```

共享 UI 和字段渲染器放在 `admin/src/components/`（其中 `ui/` 是通用组件，`resource/` 是资源组件）；通用资源页面放在 `admin/src/core/resource/pages/`。模块页面通过组合共享组件实现，复杂业务可以在模块目录内增加 `components/`，不再使用全局 `views/`。

## Migration boundary

现有前端页面已迁入 `core/pages/` 或 `modules/<name>/pages/`；后端基础入口已迁入 `core/`，用户业务入口位于 `modules/users/`。新生成资源必须遵循模块目录，不再向旧平铺目录写入文件。

## Naming rules

- `core` 只放平台基础设施，不放业务资源，也不放共享 UI 组件。
- `components` 位于前端 `src` 根目录，所有模块共享的 UI 和资源组件从这里导入。
- `modules/<name>` 是业务边界；生成器的后端与前端产物必须使用同名模块目录。
- `generated` 仅用于真正由工具维护的客户端产物，不作为业务模块目录。
